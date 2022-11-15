package secure_channel

import (
	"errors"
	"github.com/galenliu/chip/lib"
	"github.com/galenliu/chip/messageing"
	"github.com/galenliu/chip/messageing/transport"
	tlv2 "github.com/galenliu/chip/pkg/tlv"
)

type PairingSessionBase interface {
	transport.SessionDelegate
	SecureSessionType() uint8
	Peer() lib.ScopedNodeId

	LocalScopedNodeId() lib.ScopedNodeId
	PeerCATs() lib.CATValues
	GetNewSessionHandlingPolicy() transport.NewSessionHandlingPolicy

	CopySecureSession() *transport.SessionHandle

	IsValidPeerSessionId() bool
	DeriveSecureSession(ctx transport.CryptoContext) error
	GetRemoteMRPConfig() *transport.ReliableMessageProtocolConfig
	SetRemoteMRPConfig(mrpLocalConfig *transport.ReliableMessageProtocolConfig)
	EncodeMRPParameters(tag tlv2.Tag, mrpLocalConfig *transport.ReliableMessageProtocolConfig)
}

type PairingSession struct {
	mRole              uint8
	mSecureSessionType transport.SecureSessionType
	mPeerSessionId     uint16

	mDelegate SessionEstablishmentDelegate

	SessionManager       transport.SessionManagerBase
	mExchangeCtxt        messageing.ExchangeContext
	mSecureSessionHolder *transport.SessionHolderWithDelegate
	mLocalMRPConfig      *transport.ReliableMessageProtocolConfig
}

func (p PairingSession) SecureSessionType() uint8 {
	return uint8(p.mSecureSessionType)
}

func (p PairingSession) PeerCATs() lib.CATValues {
	//TODO implement me
	panic("implement me")
}

func (p PairingSession) LocalSessionId() (uint16, error) {
	//ss := p.mSecureSessionHolder.SessionHandle()
	//if ss != nil {
	//	return ss.
	//}
	return 0, errors.New("secure session is not available")
}

func (p PairingSession) GetNewSessionHandlingPolicy() transport.NewSessionHandlingPolicy {
	return transport.KStayAtOldSession
}

func (p PairingSession) CopySecureSession() *transport.SessionHandle {
	//TODO implement me
	panic("implement me")
}

func (p PairingSession) PeerSessionId() uint16 {
	//TODO implement me
	panic("implement me")
}

func (p PairingSession) IsValidPeerSessionId() bool {
	//TODO implement me
	panic("implement me")
}

func (p PairingSession) DeriveSecureSession(ctx transport.CryptoContext) error {
	//TODO implement me
	panic("implement me")
}

func (p PairingSession) RemoteMRPConfig() *transport.ReliableMessageProtocolConfig {
	//TODO implement me
	panic("implement me")
}

func (p PairingSession) SetRemoteMRPConfig(mrpLocalConfig *transport.ReliableMessageProtocolConfig) {
	//TODO implement me
	panic("implement me")
}

func (p PairingSession) EncodeMRPParameters(tag tlv2.Tag, mrpLocalConfig *transport.ReliableMessageProtocolConfig) {
	//TODO implement me
	panic("implement me")
}

func (p PairingSession) DecodeMRPParametersIfPresent(tag tlv2.Tag, reader *tlv2.ReaderImpl) error {
	if reader.GetTag() != tag {
		return nil
	}
	return nil
}

func NewPairingSessionImpl() *PairingSession {
	return &PairingSession{
		mRole:                0,
		mSecureSessionType:   transport.SecureSessionTypeCASE,
		SessionManager:       nil,
		mExchangeCtxt:        messageing.ExchangeContext{},
		mSecureSessionHolder: nil,
		mLocalMRPConfig:      nil,
	}
}
