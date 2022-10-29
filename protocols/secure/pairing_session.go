package secure

import (
	"errors"
	"github.com/galenliu/chip/lib"
	"github.com/galenliu/chip/messageing"
	"github.com/galenliu/chip/messageing/transport"
	tlv2 "github.com/galenliu/chip/pkg/tlv"
)

type PairingSession interface {
	transport.SessionDelegate
	SecureSessionType() uint8
	Peer() lib.ScopedNodeId

	LocalScopedNodeId() lib.ScopedNodeId
	PeerCATs() lib.CATValues
	GetNewSessionHandlingPolicy() transport.NewSessionHandlingPolicy

	CopySecureSession() transport.SessionHandleBase

	IsValidPeerSessionId() bool
	DeriveSecureSession(ctx transport.CryptoContext) error
	GetRemoteMRPConfig() *transport.ReliableMessageProtocolConfig
	SetRemoteMRPConfig(mrpLocalConfig *transport.ReliableMessageProtocolConfig)
	EncodeMRPParameters(tag tlv2.Tag, mrpLocalConfig *transport.ReliableMessageProtocolConfig)
}

type PairingSessionImpl struct {
	mRole              uint8
	mSecureSessionType transport.TypeSecureSession
	mPeerSessionId     uint16

	mDelegate SessionEstablishmentDelegate

	SessionManager       transport.SessionManager
	mExchangeCtxt        messageing.ExchangeContext
	mSecureSessionHolder transport.SessionHolderWithDelegate
	mLocalMRPConfig      *transport.ReliableMessageProtocolConfig
}

func (p PairingSessionImpl) SecureSessionType() uint8 {
	return uint8(p.mSecureSessionType)
}

func (p PairingSessionImpl) PeerCATs() lib.CATValues {
	//TODO implement me
	panic("implement me")
}

func (p PairingSessionImpl) LocalSessionId() (uint16, error) {
	//ss := p.mSecureSessionHolder.SessionHandle()
	//if ss != nil {
	//	return ss.
	//}
	return 0, errors.New("secure session is not available")
}

func (p PairingSessionImpl) GetNewSessionHandlingPolicy() transport.NewSessionHandlingPolicy {
	return transport.KStayAtOldSession
}

func (p PairingSessionImpl) CopySecureSession() transport.SessionHandleBase {
	//TODO implement me
	panic("implement me")
}

func (p PairingSessionImpl) PeerSessionId() uint16 {
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

func (p PairingSessionImpl) RemoteMRPConfig() *transport.ReliableMessageProtocolConfig {
	//TODO implement me
	panic("implement me")
}

func (p PairingSessionImpl) SetRemoteMRPConfig(mrpLocalConfig *transport.ReliableMessageProtocolConfig) {
	//TODO implement me
	panic("implement me")
}

func (p PairingSessionImpl) EncodeMRPParameters(tag tlv2.Tag, mrpLocalConfig *transport.ReliableMessageProtocolConfig) {
	//TODO implement me
	panic("implement me")
}

func (p PairingSessionImpl) DecodeMRPParametersIfPresent(tag tlv2.Tag, reader *tlv2.ReaderImpl) error {
	if reader.GetTag() != tag {
		return nil
	}
	return nil
}

func NewPairingSessionImpl() *PairingSessionImpl {
	return &PairingSessionImpl{
		mRole:                0,
		mSecureSessionType:   transport.CASE,
		SessionManager:       nil,
		mExchangeCtxt:        messageing.ExchangeContext{},
		mSecureSessionHolder: nil,
		mLocalMRPConfig:      nil,
	}
}
