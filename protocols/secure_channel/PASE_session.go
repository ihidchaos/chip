package secure_channel

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"github.com/galenliu/chip/crypto"
	"github.com/galenliu/chip/lib/setup"
	"github.com/galenliu/chip/lib/tlv"
	"github.com/galenliu/chip/messageing"
	"github.com/galenliu/chip/messageing/transport"
	"github.com/galenliu/chip/messageing/transport/raw"
	"github.com/galenliu/chip/messageing/transport/session"
	"github.com/galenliu/chip/platform/system"
	"github.com/moznion/go-optional"
	"golang.org/x/exp/rand"
	log "golang.org/x/exp/slog"
	"hash"
	"io"
)

type PASESession struct {
	PairingSession
	localMRPConfig       *messageing.ReliableMessageProtocolConfig
	remoteMRPConfig      *messageing.ReliableMessageProtocolConfig
	delegate             SessionEstablishmentDelegate
	exchangeCtxt         *messageing.ExchangeContext
	sessionManager       *transport.SessionManager
	role                 session.Role
	secureSessionHolder  *transport.SessionHolderWithDelegate
	pbkdfLocalRandomData []byte
	havePBKDFParameters  bool
	commissioningHash    hash.Hash
	nextExpectedMsg      messageing.MessageType
	setupPINCode         uint32

	iterationCount uint32
	saltLength     uint16
	salt           []byte

	spake2p crypto.P256Sha256HkdfHmac
}

func NewPASESession() *PASESession {
	s := &PASESession{
		localMRPConfig:      nil,
		remoteMRPConfig:     messageing.DefaultMRPConfig(),
		delegate:            nil,
		exchangeCtxt:        nil,
		sessionManager:      nil,
		role:                0,
		secureSessionHolder: nil,
		commissioningHash:   sha256.New(),
	}
	return s
}

const kSpake2pContext = "CHIP PAKE V1 Commissioning"

func (s *PASESession) Init(
	sessionManager *transport.SessionManager,
	setupCode uint32,
	delegate SessionEstablishmentDelegate,
) (err error) {
	if delegate == nil {
		return s.invalidArgumentError("SessionEstablishmentDelegate nil")
	}
	s.clear()
	s.commissioningHash = sha256.New()
	s.commissioningHash.Write([]byte(kSpake2pContext))
	s.delegate = delegate
	if err = s.PairingSession.AllocateSecureSession(sessionManager); err != nil {
		return err
	}
	sessionId := s.LocalSessionId()
	if sessionId == nil {
		return s.invalidArgumentError("SessionID nil")
	}
	log.Info("SecureChannel assigned local session key ID %d", *sessionId)
	if setupCode >= 1<<setup.PINCodeFieldLengthInBits {
		return s.invalidArgumentError("Setup Code too long")
	}
	s.setupPINCode = setupCode
	return err
}

const kPBKDFParamRandomNumberSize = 32
const defaultCommissioningPasscodeId uint16 = 0

func (s *PASESession) SendPBKDFParamRequest() (err error) {
	log.Debug("PASESession SendPBKDFParamRequest")
	s.PairingSession.LocalSessionId()
	sessionId := s.LocalSessionId()
	if sessionId == nil {
		return s.invalidArgumentError("PASESession Local Session Id is nil")
	}
	s.pbkdfLocalRandomData = make([]byte, kPBKDFParamRandomNumberSize)
	_, err = rand.Read(s.pbkdfLocalRandomData)
	if err != nil {
		return err
	}

	buf := bytes.NewBuffer(make([]byte, 0))
	tlvEncode := tlv.NewEncoder(buf)
	outContainerType := tlv.TypeNotSpecified
	outContainerType, err = tlvEncode.StartContainer(tlv.AnonymousTag(), tlv.TypeStructure)

	if err = tlvEncode.PutBytes(tlv.ContextTag(1), s.pbkdfLocalRandomData); err != nil {
		return err
	}

	if err = tlvEncode.PutUint(tlv.ContextTag(2), uint64(*sessionId)); err != nil {
		return err
	}

	if err = tlvEncode.PutUint(tlv.ContextTag(3), uint64(defaultCommissioningPasscodeId)); err != nil {
		return err
	}

	if err = tlvEncode.PutBoolean(tlv.ContextTag(4), s.havePBKDFParameters); err != nil {
		return err
	}

	if s.localMRPConfig != nil {
		if err = s.PairingSession.encodeMRPParameters(tlvEncode, tlv.ContextTag(5), s.localMRPConfig); err != nil {
			return err
		}
	}
	if err := tlvEncode.EndContainer(outContainerType); err != nil {
		return err
	}

	s.commissioningHash.Write(buf.Bytes())
	if err = s.exchangeCtxt.SendMessage(PBKDFParamRequest, buf.Bytes(), messageing.ExpectResponse); err != nil {
		return err
	}
	s.nextExpectedMsg = PBKDFParamResponse
	log.Debug("Secure Channel Sent PBKDF param request")
	return err
}

func (s *PASESession) Pair(
	sessionManager *transport.SessionManager,
	peerSetUpPINCode uint32,
	mrpLocalConfig *messageing.ReliableMessageProtocolConfig,
	exchangeCtxt *messageing.ExchangeContext,
	delegate SessionEstablishmentDelegate) (err error) {

	log.Debug("Pair PASE Session")
	if exchangeCtxt == nil {
		return s.invalidArgumentError("exchange context nil")
	}
	if err = s.Init(sessionManager, peerSetUpPINCode, delegate); err != nil {
		s.clear()
		return err
	}
	s.role = session.Initiator
	s.exchangeCtxt = exchangeCtxt
	s.exchangeCtxt.GetSessionHandle().Session.(*session.Unauthenticated).MarkActiveRx()
	s.exchangeCtxt.UseSuggestedResponseTimeout(kExpectedHighProcessingTime)
	s.localMRPConfig = mrpLocalConfig
	if err = s.SendPBKDFParamRequest(); err != nil {
		s.clear()
		return err
	}
	s.delegate.OnSessionEstablishmentStarted()
	return nil
}

func (s *PASESession) OnMessageReceived(context *messageing.ExchangeContext, header *raw.PayloadHeader, data *system.PacketBufferHandle) error {
	err := s.validateReceivedMessage(context, header, data)
	if err != nil {
		return err
	}
	msgType := MsgType(header.MessageType())
	switch msgType {
	case PBKDFParamRequest:
		err = s.HandlePBKDFParamRequest(data)
	case PBKDFParamResponse:
		err = s.HandlePBKDFParamResponse(data)
	case PASE_Pake1:
		err = s.HandleMsg1AndSendMsg2(data)
	case PASE_Pake2:
		err = s.HandleMsg2AndSendMsg3(data)
	case PASE_Pake3:
		err = s.HandleMsg3(data)
	default:
		return fmt.Errorf("invaild message type:%d", msgType)
	}
	return err
}

func (s *PASESession) OnResponseTimeout(ec *messageing.ExchangeContext) {
	//TODO implement me
	panic("implement me")
}

func (s *PASESession) OnExchangeClosing(ec *messageing.ExchangeContext) {
	//TODO implement me
	panic("implement me")
}

func (s *PASESession) GetMessageDispatch() messageing.ExchangeMessageDispatchBase {
	//TODO implement me
	panic("implement me")
}

func (s *PASESession) OnUnsolicitedMessageReceived(header *raw.PayloadHeader) (messageing.ExchangeDelegate, error) {
	return s, nil
}

func (s *PASESession) OnExchangeCreationFailed(delegate messageing.ExchangeDelegate) {
	//TODO implement me
	panic("implement me")
}

func (s *PASESession) HandlePBKDFParamRequest(w io.Reader) (err error) {
	log.Debug("PASESession HandlePBKDFParamRequest")

	defer s.PairingSession.SendStatusReport(s.exchangeCtxt, kProtocolCodeInvalidParam)

	var initiatorRandom []byte
	var initiatorSessionId uint16
	var passcodeId = defaultCommissioningPasscodeId
	var hasPBKDFParameters bool

	tlvDecode := tlv.NewDecoder(w)
	containerType := tlv.TypeStructure
	if err = tlvDecode.Type(containerType, tlv.AnonymousTag()); err != nil {
		return err
	}
	if containerType, err = tlvDecode.EnterContainer(); err != nil {
		return err
	}

	if err = tlvDecode.Next(); err != nil {
		return err
	}
	if tlvDecode.GetTag().Number() != 1 {
		return tlvDecode.TagError(tlvDecode.GetTag())
	}
	if initiatorRandom, err = tlvDecode.GetBytes(); err != nil {
		return err
	}

	if err = tlvDecode.Next(); err != nil {
		return err
	}
	if tlvDecode.GetTag().Number() != 2 {
		return tlvDecode.TagError(tlvDecode.GetTag())
	}
	if initiatorSessionId, err = tlvDecode.GetU16(); err != nil {
		return err
	}

	log.Debug("Peer assigned session ID %d", initiatorSessionId)
	s.peerSessionId = optional.Option[uint16]{initiatorSessionId}

	if err = tlvDecode.Next(); err != nil {
		return err
	}
	if tlvDecode.GetTag().Number() != 3 {
		return tlvDecode.TagError(tlvDecode.GetTag())
	}
	if passcodeId, err = tlvDecode.GetU16(); err != nil {
		return err
	}
	if passcodeId == defaultCommissioningPasscodeId {
		return s.invalidArgumentError("invalid passcodeId")
	}

	if err = tlvDecode.Next(); err != nil {
		return err
	}
	if tlvDecode.GetTag().Number() != 4 {
		return tlvDecode.TagError(tlvDecode.GetTag())
	}
	if hasPBKDFParameters, err = tlvDecode.GetBoolean(); err != nil {
		return err
	}

	if err = tlvDecode.Next(); err != nil {
		if err != tlv.EndOfTLVError {
			return err
		}
		if err = s.PairingSession.decodeMRPParametersIfPresent(tlvDecode, tlv.ContextTag(5)); err != nil {
			return err
		}
		s.exchangeCtxt.GetSessionHandle().Session.(*session.Unauthenticated).SetRemoteMRPConfig(s.remoteMRPConfig)
	}
	if err = s.SendPBKDFParamResponse(initiatorRandom, hasPBKDFParameters); err != nil {
		return err
	}
	s.delegate.OnSessionEstablishmentStarted()
	return nil
}

func (s *PASESession) HandlePBKDFParamResponse(w io.Reader) error {
	return nil

}

func (s *PASESession) HandleMsg1AndSendMsg2(w io.Reader) error {
	return nil
}

func (s *PASESession) HandleMsg2AndSendMsg3(w io.Reader) error {
	return nil
}

func (s *PASESession) HandleMsg3(w io.Reader) error {
	return nil
}

func (s *PASESession) HandleStatusReport(w io.Reader) error {
	return nil
}

func (s *PASESession) SendPBKDFParamResponse(initiatorRandom []byte, initiatorHasPBKDFParams bool) (err error) {

	log.Debug("PASESession SendPBKDFParamResponse")
	sessionId := s.LocalSessionId()
	if sessionId == nil {
		return s.invalidArgumentError("session id nil")
	}
	s.pbkdfLocalRandomData = make([]byte, kPBKDFParamRandomNumberSize)
	_, err = rand.Read(s.pbkdfLocalRandomData)
	if err != nil {
		return err
	}

	data := make([]byte, 0)
	buf := bytes.NewBuffer(data)
	tlvEncode := tlv.NewEncoder(buf)

	var outerContainerType = tlv.TypeStructure
	if outerContainerType, err = tlvEncode.StartContainer(tlv.AnonymousTag(), outerContainerType); err != nil {
		return err
	}

	if err = tlvEncode.PutBytes(tlv.ContextTag(1), initiatorRandom); err != nil {
		return err
	}
	if err = tlvEncode.PutBytes(tlv.ContextTag(2), s.pbkdfLocalRandomData); err != nil {
		return err
	}
	if err = tlvEncode.PutU16(tlv.ContextTag(3), *sessionId); err != nil {
		return err
	}
	if !initiatorHasPBKDFParams {
		var pbkdfParamContainer tlv.Type
		if pbkdfParamContainer, err = tlvEncode.StartContainer(tlv.ContextTag(4), tlv.TypeStructure); err != nil {
			return err
		}
		if err = tlvEncode.PutU32(tlv.ContextTag(1), s.iterationCount); err != nil {
			return err
		}
		if err = tlvEncode.PutBytes(tlv.ContextTag(2), s.salt); err != nil {
			return err
		}
		if err = tlvEncode.EndContainer(pbkdfParamContainer); err != nil {
			return err
		}
	}

	if s.localMRPConfig != nil {
		log.Info("SecureChannel Including MRP parameters in PBKDF param response")
		if err = s.PairingSession.encodeMRPParameters(tlvEncode, tlv.ContextTag(5), s.localMRPConfig); err != nil {
			return err
		}
	}

	s.commissioningHash.Write(buf.Bytes())

	if err = s.SetupSpake2p(); err != nil {
		return err
	}
	if err = s.exchangeCtxt.SendMessage(PBKDFParamResponse, buf.Bytes(), messageing.ExpectResponse); err != nil {
		return err
	}
	s.nextExpectedMsg = PASE_Pake1
	return nil
}

func (s *PASESession) invalidArgumentError(val any) error {
	return fmt.Errorf("PASE session invaild argument err:%v", val)
}

func (s *PASESession) validateReceivedMessage(context *messageing.ExchangeContext, header *raw.PayloadHeader, data *system.PacketBufferHandle) error {
	return nil
}

func (s *PASESession) clear() {

}

func (s *PASESession) SetupSpake2p() error {
	log.Debug("PASESession SetupSpake2p")
	return s.spake2p.Init(s.commissioningHash.Sum(nil))
}
