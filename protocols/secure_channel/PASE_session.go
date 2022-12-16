package secure_channel

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"github.com/galenliu/chip/crypto"
	"github.com/galenliu/chip/lib"
	"github.com/galenliu/chip/lib/bitflags"
	"github.com/galenliu/chip/lib/setup"
	"github.com/galenliu/chip/lib/tlv"
	"github.com/galenliu/chip/messageing"
	"github.com/galenliu/chip/messageing/transport"
	"github.com/galenliu/chip/messageing/transport/raw"
	"github.com/galenliu/chip/messageing/transport/session"
	"github.com/galenliu/chip/platform/system"
	"github.com/moznion/go-optional"
	"golang.org/x/exp/rand"
	"golang.org/x/exp/slices"
	log "golang.org/x/exp/slog"
	"hash"
	"io"
)

const kPBKDFParamRandomNumberSize = 32
const defaultCommissioningPasscodeId uint16 = 0
const kSpake2pContext = "CHIP PAKE V1 Commissioning"

type PASESession struct {
	pairingSession

	delegate             SessionEstablishmentDelegate
	sessionManager       *transport.SessionManager
	role                 session.Role
	secureSessionHolder  *transport.SessionHolderWithDelegate
	pbkdfLocalRandomData []byte
	havePBKDFParameters  bool
	commissioningHash    hash.Hash
	nextExpectedMsg      messageing.MessageType
	setupPINCode         uint32

	iterationCount uint32
	saltLength     int
	salt           []byte

	ke []byte

	spake2p      crypto.P256Sha256HkdfHmac
	paseVerifier *crypto.Spake2pVerifier

	pairingComplete bool
	mPeerCATs       lib.CATValues
}

func (s *PASESession) PeerCATs() lib.CATValues {
	return s.mPeerCATs
}

func (s *PASESession) DeriveSecureSession(ctx *transport.CryptoContext) error {
	//TODO implement me
	panic("implement me")
}

func (s *PASESession) LocalScopedNodeId() lib.ScopedNodeId {
	return lib.NewScopedNodeId(lib.UndefinedNodeId(), lib.UndefinedFabricIndex())
}

func (s *PASESession) Peer() lib.ScopedNodeId {
	return lib.NewScopedNodeId(lib.DefaultCommissioningPasscodeId.NodeId(), lib.UndefinedFabricIndex())
}

func NewPASESession() *PASESession {
	paseSession := &PASESession{
		delegate:            nil,
		sessionManager:      nil,
		role:                0,
		secureSessionHolder: nil,
		commissioningHash:   sha256.New(),
	}
	paseSession.pairingSession = NewPairingSessionImpl()
	paseSession.pairingSession.RemoteMRPConfig = messageing.DefaultMRPConfig()
	paseSession.pairingSession.base = paseSession
	return paseSession
}

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
	if err = s.pairingSession.allocateSecureSession(sessionManager, s.Peer()); err != nil {
		return err
	}

	if s.LocalSessionId().IsNone() {
		return s.invalidArgumentError("SessionID nil")
	}
	log.Info("SecureChannel assigned local session key ID %d", s.LocalSessionId().Unwrap())
	if setupCode >= 1<<setup.PINCodeFieldLengthInBits {
		return s.invalidArgumentError("Setup Code too long")
	}
	s.setupPINCode = setupCode
	return err
}

func (s *PASESession) WaitForPairing(sessionManager *transport.SessionManager,
	verifier *crypto.Spake2pVerifier,
	pbkdf2IterCount uint32,
	salt []byte,
	mrpLocalConfig *messageing.ReliableMessageProtocolConfig,
	delegate SessionEstablishmentDelegate) (err error) {
	defer func(err *error) {
		if *err != nil {
			s.clear()
		}
	}(&err)

	if err = s.Init(sessionManager, setup.PINCodeUndefinedValue, delegate); err != nil {
		return err
	}
	s.role = session.RoleResponder
	s.salt = salt
	s.saltLength = len(salt)
	s.paseVerifier = verifier
	s.iterationCount = pbkdf2IterCount
	s.nextExpectedMsg = PBKDFParamRequest
	s.pairingComplete = false
	s.RemoteMRPConfig = mrpLocalConfig

	log.Debug("PASESession Waiting for PBKDF param request")

	return nil
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
	s.role = session.RoleInitiator
	s.exchangeContext = exchangeCtxt
	s.exchangeContext.SessionHandle().Session.(*session.Unauthenticated).MarkActiveRx()
	s.exchangeContext.UseSuggestedResponseTimeout(kExpectedHighProcessingTime)
	s.LocalMRPConfig = mrpLocalConfig
	if err = s.SendPBKDFParamRequest(); err != nil {
		s.clear()
		return err
	}
	s.delegate.OnSessionEstablishmentStarted()
	return nil
}

func (s *PASESession) SendPBKDFParamRequest() (err error) {
	log.Debug("PASESession SendPBKDFParamRequest")
	s.pairingSession.LocalSessionId()

	if s.LocalSessionId().IsNone() {
		return s.invalidArgumentError("PASESession Local Session Id is nil")
	}
	s.pbkdfLocalRandomData = make([]byte, kPBKDFParamRandomNumberSize)
	_, err = rand.Read(s.pbkdfLocalRandomData)
	if err != nil {
		return err
	}

	buf := bytes.NewBuffer(make([]byte, 0))
	e := tlv.NewEncoder(buf)
	outContainerType := tlv.TypeNotSpecified
	outContainerType, err = e.StartContainer(tlv.AnonymousTag(), tlv.TypeStructure)

	if err = e.Put(tlv.ContextTag(1), s.pbkdfLocalRandomData); err != nil {
		return err
	}

	if err = e.Put(tlv.ContextTag(2), uint64(s.LocalSessionId().Unwrap())); err != nil {
		return err
	}
	if err = e.Put(tlv.ContextTag(3), uint64(defaultCommissioningPasscodeId)); err != nil {
		return err
	}

	if err = e.Put(tlv.ContextTag(4), s.havePBKDFParameters); err != nil {
		return err
	}

	if s.LocalMRPConfig != nil {
		if err = s.pairingSession.encodeMRPParameters(e, tlv.ContextTag(5), s.LocalMRPConfig); err != nil {
			return err
		}
	}
	if err := e.EndContainer(outContainerType); err != nil {
		return err
	}

	s.commissioningHash.Write(buf.Bytes())
	if err = s.exchangeContext.SendMessage(PBKDFParamRequest, buf.Bytes(), bitflags.Some(messageing.ExpectResponse)); err != nil {
		return err
	}
	s.nextExpectedMsg = PBKDFParamResponse
	log.Debug("Secure Channel Sent PBKDF param request")
	return err
}

func (s *PASESession) OnMessageReceived(context *messageing.ExchangeContext, header *raw.PayloadHeader, data *system.PacketBufferHandle) (err error) {
	defer func(err *error) {
		if *err != nil {
			s.pairingSession.discardExchange()
			s.clear()
			log.Error("SecureChannel Failed during PASE session setup", *err)
			s.pairingSession.notifySessionEstablishmentError(*err)
		}
	}(&err)

	if err = s.validateReceivedMessage(context, header, data); err != nil {
		return err
	}
	msgType := MsgType(header.MessageType)
	switch msgType {
	case PBKDFParamRequest:
		err = s.HandlePBKDFParamRequest(data)
	case PBKDFParamResponse:
		err = s.HandlePBKDFParamResponse(data)
	case PASEPake1:
		err = s.HandleMsg1AndSendMsg2(data)
	case PASEPake2:
		err = s.HandleMsg2AndSendMsg3(data)
	case PASEPake3:
		err = s.HandleMsg3(data)
	case CASEStatusReport:
		s.pairingSession.handleStatusReport(data, s.nextExpectedMsg == CASEStatusReport)
	default:
		return fmt.Errorf("invaild message type:%d", msgType)
	}
	return err
}

func (s *PASESession) OnSessionReleased() {
	s.pairingSession.OnSessionReleased()
	s.clear()
}

func (s *PASESession) OnResponseTimeout(ec *messageing.ExchangeContext) {
	if ec == nil {
		log.Info("PASESession::OnResponseTimeout was called by null exchange")
		return
	}
	if s.exchangeContext != ec {
		log.Info("PASESession::OnResponseTimeout exchange doesn't match")
		return
	}
	log.Info("PASESession timed out while waiting for a response from the peer.", "ExpectedMessageType", s.nextExpectedMsg)
	s.discardExchange()
	s.clear()
	s.pairingSession.NotifySessionEstablishmentError(ErrorTimeOut)
}

func (s *PASESession) OnExchangeClosing(ec *messageing.ExchangeContext) {
	//TODO implement me
	panic("implement me")
}

func (s *PASESession) GetMessageDispatch() messageing.ExchangeMessageDispatch {
	//return SessionEstablishmentExchangeDispatchInstance()
	return nil
}

func (s *PASESession) OnUnsolicitedMessageReceived(header *raw.PayloadHeader) (messageing.ExchangeDelegate, error) {
	return s, nil
}

func (s *PASESession) OnExchangeCreationFailed(delegate messageing.ExchangeDelegate) {
	//TODO implement me
	panic("implement me")
}

func (s *PASESession) HandlePBKDFParamRequest(w io.Reader) (err error) {
	lib.MatterTraceEventScope("PASESession", "HandlePBKDFParamRequest")
	defer func(err *error) {
		if *err != nil {
			s.sendStatusReport(s.exchangeContext, protocolCodeInvalidParam)
		}
	}(&err)
	var initiatorRandom []byte
	var initiatorSessionId uint16
	var passcodeId = defaultCommissioningPasscodeId
	var hasPBKDFParameters bool

	tlvDecode := tlv.NewDecoder(w)
	containerType := tlv.TypeStructure
	if err = tlvDecode.NextType(containerType, tlv.AnonymousTag()); err != nil {
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
		if err = s.pairingSession.decodeMRPParametersIfPresent(tlvDecode, tlv.ContextTag(5)); err != nil {
			return err
		}
		s.exchangeContext.SessionHandle().Session.(*session.Unauthenticated).SetRemoteMRPConfig(s.RemoteMRPConfig)
	}
	if err = s.SendPBKDFParamResponse(initiatorRandom, hasPBKDFParameters); err != nil {
		return err
	}
	s.delegate.OnSessionEstablishmentStarted()
	return nil
}

func (s *PASESession) SendPBKDFParamResponse(initiatorRandom []byte, initiatorHasPBKDFParams bool) (err error) {

	log.Debug("PASESession SendPBKDFParamResponse")

	if s.LocalSessionId().IsNone() {
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
	if err = tlvEncode.PutU16(tlv.ContextTag(3), s.LocalSessionId().Unwrap()); err != nil {
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

	if s.LocalMRPConfig != nil {
		log.Info("SecureChannel Including MRP parameters in PBKDF param response")
		if err = s.pairingSession.encodeMRPParameters(tlvEncode, tlv.ContextTag(5), s.LocalMRPConfig); err != nil {
			return err
		}
	}

	s.commissioningHash.Write(buf.Bytes())

	if err = s.setupSpake2p(); err != nil {
		return err
	}
	if err = s.exchangeContext.SendMessage(PBKDFParamResponse, buf.Bytes(), bitflags.Some(messageing.ExpectResponse)); err != nil {
		return err
	}
	s.nextExpectedMsg = PASEPake1
	return nil
}

func (s *PASESession) HandlePBKDFParamResponse(msg *system.PacketBufferHandle) (err error) {
	lib.MatterTraceEventScope("PASESession", "HandlePBKDFParamResponse")
	defer func(err *error) {
		if *err != nil {
			s.sendStatusReport(s.exchangeContext, protocolCodeInvalidParam)
		}
	}(&err)
	s.commissioningHash.Write(msg.Bytes())

	var random []byte
	var salt []byte
	var responderSessionId uint16
	var serializedWS []byte

	d := tlv.NewDecoder(msg)
	var containerType = tlv.TypeStructure
	if err = d.NextType(containerType, tlv.AnonymousTag()); err != nil {
		return err
	}
	if containerType, err = d.EnterContainer(); err != nil {
		return err
	}

	// Initiator's random value
	if err = d.Next(); err != nil {
		return err
	}
	if d.GetTag().Number() != 1 {
		return d.TagError(d.GetTag())
	}
	if random, err = d.GetBytes(); err != nil {
		return err
	}
	if slices.Equal(random, s.pbkdfLocalRandomData) {
		return s.invalidArgumentError(fmt.Sprintf("PASE  RandomData err:%d", random))
	}

	// RoleResponder's random value
	if err = d.Next(); err != nil {
		return err
	}
	if d.GetTag().Number() != 2 {
		return d.TagError(d.GetTag())
	}
	if random, err = d.GetBytes(); err != nil {
		return err
	}

	if err = d.Next(); err != nil {
		return err
	}
	if d.GetTag().Number() != 3 {
		return d.TagError(d.GetTag())
	}
	if responderSessionId, err = d.GetU16(); err != nil {
		return err
	}
	log.Debug("SecureChannel Peer assigned session ID %d", responderSessionId)
	s.peerSessionId = optional.Option[uint16]{responderSessionId}

	if s.havePBKDFParameters {
		if err = d.Next(); err != tlv.EndOfTLVError {
			if err = s.decodeMRPParametersIfPresent(d, tlv.ContextTag(5)); err != nil {
				return err
			}
			s.exchangeContext.SessionHandle().Session.(*session.Unauthenticated).SetRemoteMRPConfig(s.RemoteMRPConfig)
		}
		salt = s.salt
	} else {

		if err = d.Next(); err != nil {
			return err
		}
		if containerType, err = d.EnterContainer(); err != nil {
			return err
		}

		if err = d.Next(); err != nil {
			return err
		}
		if d.GetTag().Number() != 1 {
			return d.TagError(d.GetTag())
		}

		if s.iterationCount, err = d.GetU32(); err != nil {
			return err
		}

		if err = d.Next(); err != nil {
			return err
		}
		if d.GetTag().Number() != 2 {
			return d.TagError(d.GetTag())
		}
		if salt, err = d.GetBytes(); err != nil {
			return err
		}

		if err = d.ExitContainer(containerType); err != nil {
			return err
		}
		if err = d.Next(); err != tlv.EndOfTLVError {
			if err = s.decodeMRPParametersIfPresent(d, tlv.ContextTag(5)); err != nil {
				return err
			}
			s.exchangeContext.SessionHandle().Session.(*session.Unauthenticated).SetRemoteMRPConfig(s.RemoteMRPConfig)
		}
	}
	if err = s.setupSpake2p(); err != nil {
		return err
	}
	if serializedWS, err = crypto.ComputeWS(s.iterationCount, s.setupPINCode, salt); err != nil {
		return err
	}
	if err = s.spake2p.BeginProver(serializedWS); err != nil {
		return err
	}
	if err = s.SendMsg1(); err != nil {
		return err
	}
	return
}

func (s *PASESession) SendMsg1() (err error) {

	lib.MatterTraceEventScope("PASESession", "SendMsg1")
	buf := new(bytes.Buffer)
	e := tlv.NewEncoder(buf)
	var outerContainerType = tlv.TypeNotSpecified
	if outerContainerType, err = e.StartContainer(tlv.AnonymousTag(), tlv.TypeStructure); err != nil {
		return err
	}

	var X = make([]byte, crypto.MaxPointLength)
	var Y = make([]byte, 0)
	if err = s.spake2p.ComputeRoundOne(X, Y); err != nil {
		return err
	}
	if err = e.Put(tlv.ContextTag(1), X); err != nil {
		return err
	}
	if err = e.EndContainer(outerContainerType); err != nil {
		return err
	}
	if err = s.exchangeContext.SendMessage(PASEPake1, buf.Bytes(), bitflags.Some(messageing.ExpectResponse)); err != nil {
		return err
	}

	log.Debug("SecureChannel Sent spake2p msg1")
	s.nextExpectedMsg = PASEPake2
	return nil
}

func (s *PASESession) HandleMsg1AndSendMsg2(w io.Reader) (err error) {

	lib.MatterTraceEventScope("PASESession", "HandleMsg1AndSendMsg2")
	defer func(err *error) {
		if *err != nil {
			s.sendStatusReport(s.exchangeContext, protocolCodeInvalidParam)
		}
	}(&err)

	var Y = make([]byte, crypto.MaxPointLength)
	var X = make([]byte, 0)
	var verifier = make([]byte, crypto.MaxHashLength)

	log.Debug("SecureChannel Received spake2p msg1")

	tlvDecoder := tlv.NewDecoder(w)
	var containerType = tlv.TypeStructure
	if err = tlvDecoder.NextType(containerType, tlv.AnonymousTag()); err != nil {
		return err
	}
	if containerType, err = tlvDecoder.EnterContainer(); err != nil {
		return err
	}

	if err = tlvDecoder.Next(); err != nil {
		return err
	}
	if tlvDecoder.GetTag().Number() != 1 {
		return tlvDecoder.TagError(tlvDecoder.GetTag())
	}
	if err = tlvDecoder.Get(X); err != nil {
		return err
	}
	if err = s.spake2p.BeginVerifier(nil, 0, nil, 0, s.paseVerifier.W0, s.paseVerifier.ML); err != nil {
		return err
	}

	if err = s.spake2p.ComputeRoundOne(X, Y); err != nil {
		return err
	}
	if err = s.spake2p.ComputeRoundTwo(X, verifier); err != nil {
		return err
	}

	{
		buf := new(bytes.Buffer)
		tlvEncoder := tlv.NewEncoder(buf)
		var outerContainerType = tlv.TypeNotSpecified
		if outerContainerType, err = tlvEncoder.StartContainer(tlv.AnonymousTag(), tlv.TypeStructure); err != nil {
			return err
		}
		if err = tlvEncoder.PutBytes(tlv.ContextTag(1), Y); err != nil {
			return err
		}
		if err = tlvEncoder.PutBytes(tlv.ContextTag(2), verifier); err != nil {
			return err
		}
		if err = tlvDecoder.ExitContainer(outerContainerType); err != nil {
			return err
		}
		if err = s.exchangeContext.SendMessage(PASEPake2, buf.Bytes(), bitflags.Some(messageing.ExpectResponse)); err != nil {
			return err
		}
		s.nextExpectedMsg = PASEPake3
	}
	log.Debug("SecureChannel Sent spake2p msg2")

	return nil
}

func (s *PASESession) HandleMsg2AndSendMsg3(w io.Reader) (err error) {
	lib.MatterTraceEventScope("PASESession", "HandleMsg2AndSendMsg3")
	defer func(err *error) {
		if *err != nil {
			s.sendStatusReport(s.exchangeContext, protocolCodeInvalidParam)
		}
	}(&err)

	var verifier = make([]byte, crypto.MaxHashLength)
	var Y []byte
	var peerVerifier []byte

	log.Debug("SecureChannel Received spake2p msg2")

	var containerType = tlv.TypeStructure
	tlvDecoder := tlv.NewDecoder(w)
	if err = tlvDecoder.NextType(containerType, tlv.AnonymousTag()); err != nil {
		return err
	}
	if containerType, err = tlvDecoder.EnterContainer(); err != nil {
		return err
	}

	if err = tlvDecoder.Next(); err != nil {
		return err
	}
	if tlvDecoder.GetTag().Number() != 1 {
		return tlvDecoder.TagError(tlvDecoder.GetTag())
	}
	if Y, err = tlvDecoder.GetBytes(); err != nil {
		return err
	}

	if err = tlvDecoder.Next(); err != nil {
		return err
	}
	if tlvDecoder.GetTag().Number() != 2 {
		return tlvDecoder.TagError(tlvDecoder.GetTag())
	}
	if peerVerifier, err = tlvDecoder.GetBytes(); err != nil {
		return err
	}

	if err = s.spake2p.ComputeRoundTwo(Y, verifier); err != nil {
		return err
	}
	if err = s.spake2p.KeyConfirm(peerVerifier); err != nil {
		return err
	}

	s.ke = make([]byte, crypto.MaxHashLength)
	if err = s.spake2p.GetKeys(s.ke); err != nil {
		return err
	}

	//SendMessage
	{
		buf := new(bytes.Buffer)
		tlvEncoder := tlv.NewEncoder(buf)
		var outerContainerType = tlv.TypeNotSpecified
		if outerContainerType, err = tlvEncoder.StartContainer(tlv.AnonymousTag(), tlv.TypeStructure); err != nil {
			return err
		}
		if err = tlvEncoder.PutBytes(tlv.ContextTag(1), verifier); err != nil {
			return err
		}
		if err = tlvEncoder.EndContainer(outerContainerType); err != nil {
			return err
		}
		if err = s.exchangeContext.SendMessage(PASEPake3, buf.Bytes(), bitflags.Some(messageing.ExpectResponse)); err != nil {
			return err
		}
		s.nextExpectedMsg = CASEStatusReport
	}
	log.Debug("SecureChannel Sent spake2p msg3")
	return nil
}

func (s *PASESession) HandleMsg3(w io.Reader) (err error) {
	lib.MatterTraceEventScope("PASESession", "HandleMsg3")
	defer func(err *error) {
		if *err != nil {
			s.sendStatusReport(s.exchangeContext, protocolCodeInvalidParam)
		}
	}(&err)

	log.Debug("SecureChannel Received spake2p msg3")
	s.nextExpectedMsg = nil

	var containerType = tlv.TypeStructure
	var peerVerifier = make([]byte, crypto.MaxHashLength)

	tlvDecoder := tlv.NewDecoder(w)
	if err = tlvDecoder.NextType(containerType, tlv.AnonymousTag()); err != nil {
		return err
	}
	if containerType, err = tlvDecoder.EnterContainer(); err != nil {
		return err
	}

	if err = tlvDecoder.Next(); err != nil {
		return err
	}
	if tlvDecoder.GetTag().Number() != 1 {
		return tlvDecoder.TagError(tlvDecoder.GetTag())
	}
	if err = tlvDecoder.Get(peerVerifier); err != nil {
		return err
	}

	if err = s.spake2p.KeyConfirm(peerVerifier); err != nil {
		return err
	}

	s.ke = make([]byte, crypto.MaxHashLength)
	if err = s.spake2p.GetKeys(s.ke); err != nil {
		return err
	}

	s.sendStatusReport(s.exchangeContext, protocolCodeSuccess)

	s.finish()

	return nil
}

func (s *PASESession) OnSuccessStatusReport() {
	s.finish()
}

func (s *PASESession) invalidArgumentError(val any) error {
	return fmt.Errorf("PASESession invaild argument:%v", val)
}

func (s *PASESession) validateReceivedMessage(context *messageing.ExchangeContext, header *raw.PayloadHeader, data *system.PacketBufferHandle) error {
	return nil
}

func (s *PASESession) generatePASEVerifier(verifier *crypto.Spake2pVerifier,
	pbkdf2IterCount uint32,
	salt []byte,
	useRandomPIN bool,
	setupPINCode uint32) error {
	lib.MatterTraceEventScope("PASESession", "generatePASEVerifier")
	if useRandomPIN {
		setupPINCode = rand.Uint32()
		setupPINCode = (setupPINCode & setup.PINCodeMaximumValue) + 1
	}
	return verifier.Generate(pbkdf2IterCount, salt, setupPINCode)
}

func (s *PASESession) setupSpake2p() error {
	log.Debug("PASESession setupSpake2p")
	return s.spake2p.Init(s.commissioningHash.Sum(nil))
}

func (s *PASESession) clear() {
	s.ke = nil
	s.nextExpectedMsg = nil
	s.spake2p.Clear()
	s.commissioningHash.Reset()
	s.iterationCount = 0
	s.salt = nil
	s.pairingComplete = false
	s.pairingSession.clear()
}

func (s *PASESession) finish() {
	s.pairingComplete = true
	s.pairingSession.finish()
}
