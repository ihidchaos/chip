package secure_channel

import (
	"github.com/galenliu/chip/lib"
	"github.com/galenliu/chip/lib/tlv"
	"github.com/galenliu/chip/messageing"
	"github.com/galenliu/chip/messageing/transport"
	"github.com/galenliu/chip/messageing/transport/session"
	"github.com/moznion/go-optional"
	log "golang.org/x/exp/slog"
	"time"
)

type PairingSession struct {
	role              session.Role
	secureSessionType session.SecureSessionType
	peerSessionId     optional.Option[uint16]

	delegate SessionEstablishmentDelegate

	sessionManager      *transport.SessionManager
	exchangeContext     *messageing.ExchangeContext
	secureSessionHolder *transport.SessionHolderWithDelegate
	localMRPConfig      *messageing.ReliableMessageProtocolConfig
	remoteMRPConfig     *messageing.ReliableMessageProtocolConfig
}

func NewPairingSessionImpl() PairingSession {
	return PairingSession{
		role:                0,
		secureSessionType:   session.SecureSessionTypeCASE,
		sessionManager:      nil,
		exchangeContext:     &messageing.ExchangeContext{},
		secureSessionHolder: nil,
		localMRPConfig:      nil,
		remoteMRPConfig:     messageing.GetLocalMRPConfig(),
	}
}

func (s *PairingSession) PeerCATs() lib.CATValues {
	//TODO implement me
	panic("implement me")
}

func (s *PairingSession) LocalSessionId() *uint16 {
	if s.secureSessionHolder == nil {
		return nil
	}
	if s.secureSessionHolder.SessionType() == session.SecureType {
		sessionId := s.secureSessionHolder.Session.(*session.Secure).LocalSessionId()
		return &sessionId
	}
	return nil
}

func (s *PairingSession) GetNewSessionHandlingPolicy() transport.NewSessionHandlingPolicy {
	return transport.StayAtOldSession
}

func (s *PairingSession) CopySecureSession() *transport.SessionHandle {
	//TODO implement me
	panic("implement me")
}

func (s *PairingSession) PeerSessionId() uint16 {
	return s.peerSessionId.Unwrap()
}

func (s *PairingSession) IsValidPeerSessionId() bool {
	//TODO implement me
	panic("implement me")
}

func (s *PairingSession) DeriveSecureSession(ctx session.CryptoContext) error {
	//TODO implement me
	panic("implement me")
}

func (s *PairingSession) RemoteMRPConfig() *messageing.ReliableMessageProtocolConfig {
	return s.remoteMRPConfig
}

func (s *PairingSession) SetRemoteMRPConfig(mrpLocalConfig *messageing.ReliableMessageProtocolConfig) {
	//TODO implement me
	panic("implement me")
}

func (s *PairingSession) encodeMRPParameters(tlvEncode *tlv.Encoder, tag tlv.Tag, mrpLocalConfig *messageing.ReliableMessageProtocolConfig) error {
	mrpParamsContainer, err := tlvEncode.StartContainer(tag, tlv.TypeStructure)
	if err != nil {
		return err
	}
	if err = tlvEncode.PutUint(tlv.ContextTag(1), uint64(mrpLocalConfig.IdleRetransTimeout.Milliseconds())); err != nil {
		return err
	}

	if err = tlvEncode.PutUint(tlv.ContextTag(2), uint64(mrpLocalConfig.ActiveRetransTimeout.Milliseconds())); err != nil {
		return err
	}
	return tlvEncode.EndContainer(mrpParamsContainer)
}

func (s *PairingSession) decodeMRPParametersIfPresent(tlvDecode *tlv.Decoder, expectedTag tlv.Tag) (err error) {
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
		s.remoteMRPConfig.IdleRetransTimeout = time.Duration(tlvElementValue) * time.Millisecond
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
	s.remoteMRPConfig.ActiveRetransTimeout = time.Duration(tlvElementValue) * time.Millisecond

	return tlvDecode.ExitContainer(container)
}

func (s *PairingSession) AllocateSecureSession(manager *transport.SessionManager) error {
	return nil
}

func (s *PairingSession) SendStatusReport(ctxt *messageing.ExchangeContext, param uint16) {

}
