package secure_channel

import (
	"errors"
	"github.com/galenliu/chip/lib"
	"github.com/galenliu/chip/messageing"
	"github.com/galenliu/chip/messageing/transport"
	"github.com/galenliu/chip/messageing/transport/session"
	tlv2 "github.com/galenliu/chip/pkg/tlv"
)

type PairingSessionBase interface {
	session.Delegate
	SecureSessionType() uint8
	Peer() lib.ScopedNodeId
	LocalScopedNodeId() lib.ScopedNodeId
	PeerCATs() lib.CATValues
	GetNewSessionHandlingPolicy() session.NewSessionHandlingPolicy
	CopySecureSession() *transport.SessionHandle
	IsValidPeerSessionId() bool
	DeriveSecureSession(ctx session.CryptoContext) error
	GetRemoteMRPConfig() *session.ReliableMessageProtocolConfig
	SetRemoteMRPConfig(mrpLocalConfig *session.ReliableMessageProtocolConfig)
	EncodeMRPParameters(tag tlv2.Tag, mrpLocalConfig *session.ReliableMessageProtocolConfig)
}

type PairingSession struct {
	mRole              uint8
	mSecureSessionType session.SecureSessionType
	mPeerSessionId     uint16

	mDelegate SessionEstablishmentDelegate

	SessionManager       transport.SessionManagerBase
	mExchangeCtxt        *messageing.ExchangeContext
	mSecureSessionHolder *transport.SessionHolderWithDelegate
	mLocalMRPConfig      *session.ReliableMessageProtocolConfig
}

func NewPairingSessionImpl() *PairingSession {
	return &PairingSession{
		mRole:                0,
		mSecureSessionType:   session.SecureSessionTypeCASE,
		SessionManager:       nil,
		mExchangeCtxt:        &messageing.ExchangeContext{},
		mSecureSessionHolder: nil,
		mLocalMRPConfig:      nil,
	}
}

func (p PairingSession) SecureSessionType() session.SecureSessionType {
	return p.mSecureSessionType
}

func (p PairingSession) PeerCATs() lib.CATValues {
	//TODO implement me
	panic("implement me")
}

func (p PairingSession) LocalSessionId() (uint16, error) {
	return 0, errors.New("secure session is not available")
}

func (p PairingSession) GetNewSessionHandlingPolicy() session.NewSessionHandlingPolicy {
	return session.KStayAtOldSession
}

func (p PairingSession) CopySecureSession() *transport.SessionHandle {
	//TODO implement me
	panic("implement me")
}

func (p PairingSession) PeerSessionId() uint16 {
	return p.mPeerSessionId
}

func (p PairingSession) IsValidPeerSessionId() bool {
	//TODO implement me
	panic("implement me")
}

func (p PairingSession) DeriveSecureSession(ctx session.CryptoContext) error {
	//TODO implement me
	panic("implement me")
}

func (p PairingSession) RemoteMRPConfig() *session.ReliableMessageProtocolConfig {
	//TODO implement me
	panic("implement me")
}

func (p PairingSession) SetRemoteMRPConfig(mrpLocalConfig *session.ReliableMessageProtocolConfig) {
	//TODO implement me
	panic("implement me")
}

func (p PairingSession) EncodeMRPParameters(tag tlv2.Tag, mrpLocalConfig *session.ReliableMessageProtocolConfig) {
	//TODO implement me
	panic("implement me")
}

func (p PairingSession) DecodeMRPParametersIfPresent(tag tlv2.Tag, reader *tlv2.ReaderImpl) error {
	if reader.GetTag() != tag {
		return nil
	}
	return nil
}
