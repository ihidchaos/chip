package secure_channel

import (
	"github.com/galenliu/chip/lib"
	"github.com/galenliu/chip/messageing"
	"github.com/galenliu/chip/messageing/transport"
)

type PairingSession interface {
	transport.SessionDelegate
	GetSecureSessionType() uint8
	GetPeer() lib.ScopedNodeId

	GetLocalScopedNodeId() lib.ScopedNodeId
	GetPeerCATs() lib.CATValues
	GetNewSessionHandlingPolicy() transport.NewSessionHandlingPolicy

	CopySecureSession() transport.SessionHandle

	IsValidPeerSessionId() bool
	DeriveSecureSession(ctx transport.CryptoContext) error
	GetRemoteMRPConfig() *transport.ReliableMessageProtocolConfig
	SetRemoteMRPConfig(mrpLocalConfig *transport.ReliableMessageProtocolConfig)
	EncodeMRPParameters(tag lib.Tag, mrpLocalConfig *transport.ReliableMessageProtocolConfig)
}

type PairingSessionImpl struct {
	mRole              uint8
	mSecureSessionType uint8

	SessionManager       transport.SessionManager
	mExchangeCtxt        messageing.ExchangeContext
	mSecureSessionHolder transport.SessionHolderWithDelegate
	mLocalMRPConfig      *transport.ReliableMessageProtocolConfig
}

func (p PairingSessionImpl) GetSecureSessionType() uint8 {
	return p.mSecureSessionType
}

func (p PairingSessionImpl) GetPeerCATs() lib.CATValues {
	//TODO implement me
	panic("implement me")
}

func (p PairingSessionImpl) GetNewSessionHandlingPolicy() transport.NewSessionHandlingPolicy {
	return transport.KStayAtOldSession
}

func (p PairingSessionImpl) GetLocalSessionId() uint16 {
	//TODO implement me
	panic("implement me")
}

func (p PairingSessionImpl) CopySecureSession() transport.SessionHandle {
	//TODO implement me
	panic("implement me")
}

func (p PairingSessionImpl) GetPeerSessionId() uint16 {
	//TODO implement me
	panic("implement me")
}

func (p PairingSessionImpl) IsValidPeerSessionId() bool {
	//TODO implement me
	panic("implement me")
}

func (p PairingSessionImpl) DeriveSecureSession(ctx transport.CryptoContext) error {
	//TODO implement me
	panic("implement me")
}

func (p PairingSessionImpl) GetRemoteMRPConfig() *transport.ReliableMessageProtocolConfig {
	//TODO implement me
	panic("implement me")
}

func (p PairingSessionImpl) SetRemoteMRPConfig(mrpLocalConfig *transport.ReliableMessageProtocolConfig) {
	//TODO implement me
	panic("implement me")
}

func (p PairingSessionImpl) EncodeMRPParameters(tag lib.Tag, mrpLocalConfig *transport.ReliableMessageProtocolConfig) {
	//TODO implement me
	panic("implement me")
}

func NewPairingSessionImpl() *PairingSessionImpl {
	return &PairingSessionImpl{
		mRole:                0,
		mSecureSessionType:   transport.KSecureSessionCASE,
		SessionManager:       nil,
		mExchangeCtxt:        messageing.ExchangeContext{},
		mSecureSessionHolder: nil,
		mLocalMRPConfig:      nil,
	}
}
