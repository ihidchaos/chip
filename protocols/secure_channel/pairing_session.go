package secure_channel

import (
	"errors"
	"github.com/galenliu/chip/lib"
	"github.com/galenliu/chip/lib/tlv"
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
	EncodeMRPParameters(tag tlv.Tag, mrpLocalConfig *transport.ReliableMessageProtocolConfig)
}

type PairingSessionImpl struct {
	mRole              uint8
	mSecureSessionType transport.TSecureSessionType
	mPeerSessionId     uint16

	mDelegate SessionEstablishmentDelegate

	SessionManager       transport.SessionManager
	mExchangeCtxt        messageing.ExchangeContext
	mSecureSessionHolder transport.SessionHolderWithDelegate
	mLocalMRPConfig      *transport.ReliableMessageProtocolConfig
}

func (p PairingSessionImpl) GetSecureSessionType() uint8 {
	return uint8(p.mSecureSessionType)
}

func (p PairingSessionImpl) GetPeerCATs() lib.CATValues {
	//TODO implement me
	panic("implement me")
}

func (p PairingSessionImpl) GetLocalSessionId() (uint16, error) {
	ss := p.mSecureSessionHolder.AsSecureSession()
	if ss != nil {
		return ss.GetLocalSessionId(), nil
	}
	return 0, errors.New("secure session err")
}

func (p PairingSessionImpl) GetNewSessionHandlingPolicy() transport.NewSessionHandlingPolicy {
	return transport.KStayAtOldSession
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

func (p PairingSessionImpl) EncodeMRPParameters(tag tlv.Tag, mrpLocalConfig *transport.ReliableMessageProtocolConfig) {
	//TODO implement me
	panic("implement me")
}

func (p PairingSessionImpl) DecodeMRPParametersIfPresent(tag tlv.Tag, reader *tlv.ReaderImpl) error {
	if reader.GetTag() != tag {
		return nil
	}
	return nil
}

func NewPairingSessionImpl() *PairingSessionImpl {
	return &PairingSessionImpl{
		mRole:                0,
		mSecureSessionType:   transport.K_CASE,
		SessionManager:       nil,
		mExchangeCtxt:        messageing.ExchangeContext{},
		mSecureSessionHolder: nil,
		mLocalMRPConfig:      nil,
	}
}
