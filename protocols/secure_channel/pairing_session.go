package secure_channel

import (
	"bytes"
	"fmt"
	"github.com/galenliu/chip/lib"
	"github.com/galenliu/chip/lib/tlv"
	"github.com/galenliu/chip/messageing"
	"github.com/galenliu/chip/messageing/transport"
	"github.com/galenliu/chip/messageing/transport/raw"
	"github.com/galenliu/chip/messageing/transport/session"
	"github.com/galenliu/chip/platform/system"
	"github.com/galenliu/chip/protocols"
	"github.com/moznion/go-optional"
	log "golang.org/x/exp/slog"
	"time"
)

type pairingSessionBase interface {
	OnSuccessStatusReport()
	DeriveSecureSession(ctx *session.CryptoContext) error
	LocalScopedNodeId() lib.ScopedNodeId
	Peer() lib.ScopedNodeId
	PeerCATs() lib.CATValues
}

type pairingSession struct {
	role              session.Role
	secureSessionType session.SecureType
	peerSessionId     optional.Option[uint16]

	delegate SessionEstablishmentDelegate

	sessionManager      *transport.SessionManager
	exchangeContext     *messageing.ExchangeContext
	secureSessionHolder *transport.SessionHolderWithDelegate

	LocalMRPConfig  *messageing.ReliableMessageProtocolConfig
	RemoteMRPConfig *messageing.ReliableMessageProtocolConfig
	base            pairingSessionBase
}

func NewPairingSessionImpl() pairingSession {
	return pairingSession{
		role:                0,
		secureSessionType:   session.SecureSessionTypeCASE,
		sessionManager:      nil,
		exchangeContext:     &messageing.ExchangeContext{},
		secureSessionHolder: nil,
		LocalMRPConfig:      nil,
		RemoteMRPConfig:     messageing.GetLocalMRPConfig(),
	}
}

func (s *pairingSession) LocalSessionId() optional.Option[uint16] {
	if s.secureSessionHolder == nil {
		return optional.None[uint16]()
	}
	if s.secureSessionHolder.IsSecure() {
		sessionId := s.secureSessionHolder.Session.(*session.Secure).LocalSessionId()
		return optional.Some(sessionId)
	}
	return nil
}

func (s *pairingSession) GetNewSessionHandlingPolicy() transport.NewSessionHandlingPolicy {
	return transport.StayAtOldSession
}

func (s *pairingSession) CopySecureSession() *transport.SessionHandle {
	//TODO implement me
	panic("implement me")
}

func (s *pairingSession) IsValidPeerSessionId() bool {
	//TODO implement me
	panic("implement me")
}

func (s *pairingSession) DeriveSecureSession(ctx session.CryptoContext) error {
	//TODO implement me
	panic("implement me")
}

func (s *pairingSession) encodeMRPParameters(tlvEncode *tlv.Encoder, tag tlv.Tag, mrpLocalConfig *messageing.ReliableMessageProtocolConfig) error {
	mrpParamsContainer, err := tlvEncode.StartContainer(tag, tlv.TypeStructure)
	if err != nil {
		return err
	}
	if err = tlv.PutUint(tlvEncode, tlv.ContextTag(1), uint64(mrpLocalConfig.IdleRetransTimeout.Milliseconds())); err != nil {
		return err
	}

	if err = tlv.PutUint(tlvEncode, tlv.ContextTag(2), uint64(mrpLocalConfig.ActiveRetransTimeout.Milliseconds())); err != nil {
		return err
	}
	return tlvEncode.EndContainer(mrpParamsContainer)
}

func (s *pairingSession) decodeMRPParametersIfPresent(tlvDecode *tlv.Decoder, expectedTag tlv.Tag) (err error) {
	if tlvDecode.GetTag() != expectedTag {
		return tlvDecode.TagError(tlvDecode.GetTag())
	}

	var tlvElementValue uint32 = 0
	var container = tlv.TypeStructure

	container, err = tlvDecode.EnterContainer()
	if err != nil {
		return err
	}

	if err = tlvDecode.Next(); err != nil {
		return err
	}
	log.Debug("SecureChannel Found MRP parameters in the message")

	if tlvDecode.GetTag().Number() == 1 {
		if tlvElementValue, err = tlvDecode.GetU32(); err != nil {
			return err
		}
		s.RemoteMRPConfig.IdleRetransTimeout = time.Duration(tlvElementValue) * time.Millisecond
		if err = tlvDecode.Next(); err != nil {
			if err != tlv.EndOfTLVError {
				return err
			}
			if err = tlvDecode.ExitContainer(container); err != nil {
				return err
			}
		}
	}

	if tlvDecode.GetTag().Number() != 2 {
		return tlvDecode.TagError(tlvDecode.GetTag())
	}
	if tlvElementValue, err = tlvDecode.GetU32(); err != nil {
		return err
	}
	s.RemoteMRPConfig.ActiveRetransTimeout = time.Duration(tlvElementValue) * time.Millisecond

	return tlvDecode.ExitContainer(container)
}

func (s *pairingSession) allocateSecureSession(sessionManager *transport.SessionManager, sessionEvictionHint lib.ScopedNodeId) error {
	handler := sessionManager.AllocateSession(s.secureSessionType, sessionEvictionHint)
	if handler == nil {
		return ErrorNoMemory
	}
	if !s.secureSessionHolder.GrabPairingSession(handler) {
		return fmt.Errorf("ERROR_INTERNAL")
	}
	s.sessionManager = sessionManager
	return nil
}

func (s *pairingSession) activateSecureSession(address raw.PeerAddress) (err error) {
	secureSession := s.secureSessionHolder.Session.(*session.Secure)
	context := secureSession.GetCryptoContext()
	err = s.base.DeriveSecureSession(context)

	secureSession.SetPeerAddress(address)
	secureSession.SessionMessageCounter().PeerMessageCounter.SetCounter(session.InitialSyncValue)

	secureSession.Activate(s.base.LocalScopedNodeId(), s.base.Peer(), s.base.PeerCATs(), s.peerSessionId.Unwrap(), s.RemoteMRPConfig)
	log.Debug("New secure session activated for device", "LSID", s.base.Peer(), "PSID", s.peerSessionId.Unwrap())
	return err
}

func (s *pairingSession) finish() {
	address := s.exchangeContext.SessionHandle().Session.(*session.Unauthenticated).PeerAddress()
	s.discardExchange()
	err := s.activateSecureSession(address)
	if err != nil {
		s.notifySessionEstablishmentError(err)
	} else {
		s.delegate.OnSessionEstablished(s.secureSessionHolder.SessionHandler())
	}
}

func (s *pairingSession) clear() {

}

func (s *pairingSession) notifySessionEstablishmentError(err error) {

}

func (s *pairingSession) discardExchange() {
	if s.exchangeContext != nil {
		s.exchangeContext.SetDelegate(nil)
		s.exchangeContext = nil
	}
}

func (s *pairingSession) handleStatusReport(msg *system.PacketBufferHandle, successExpected bool) (err error) {

	var report = &StatusReport{}
	if err = report.Decode(msg.Buffer); err != nil {
		return
	}
	if report.GeneralCode == kSuccess &&
		report.ProtocolCode == protocolCodeSuccess &&
		successExpected {
		if s.base != nil {
			s.base.OnSuccessStatusReport()
		}
	}

	return err
}

func (s *pairingSession) sendStatusReport(ctxt *messageing.ExchangeContext, protocolCode uint16) {
	var code = kFailure
	if protocolCode == protocolCodeSuccess {
		code = kSuccess
	}
	log.Debug("SecureChannel", "msg", "Sending status report", "ProtocolCode", protocolCode, "exchange", s.exchangeContext.ExchangeId())
	statusReport := StatusReport{
		ProtocolCode: protocolCode,
		ProtocolId:   protocols.New(protocolId, nil),
		GeneralCode:  code,
	}
	buf := bytes.NewBuffer(nil)
	if err := statusReport.Encode(buf); err != nil {
		log.Error("SecureChannel", err, "msg", "Failed to allocate status report message")
		return
	}
	if err := s.exchangeContext.SendMessage(CASEStatusReport, buf.Bytes(), messageing.None); err != nil {
		log.Error("SecureChannel", err, "msg", "Failed to send status report message")
	}
}

func (s *pairingSession) IsSessionEstablishmentInProgress() bool {
	if s.secureSessionHolder == nil {
		return false
	}
	secureSession := s.secureSessionHolder.Session.(*session.Secure)
	return secureSession.IsEstablishing()
}

func (s *pairingSession) OnSessionReleased() {

}

func (s *pairingSession) NotifySessionEstablishmentError(err error) {
	if s.delegate == nil {
		return
	}
	s.delegate.OnSessionEstablishmentError(err)
}
