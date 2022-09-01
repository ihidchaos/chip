package transport

import (
	"github.com/galenliu/chip/lib"
	"time"
)

type SecureSessionBase interface {
	Session
}

type SecureSession struct {
	*SessionBaseImpl
	mState             StateSecureSession
	mTable             *SecureSessionTable
	mSecureSessionType TypeSecureSession
	mLocalSessionId    uint16
	mPeerSessionId     uint16
	mLocalNodeId       lib.NodeId
	mPeerNodeId        lib.NodeId
	mRemoteMRPConfig   *ReliableMessageProtocolConfig
	mCryptoContext     *CryptoContext
	mPeerCATs          *lib.CATValues
}

func NewSecureSession(
	table *SecureSessionTable,
	secureSessionType TypeSecureSession,
	localSessionId uint16,
) *SecureSession {
	return &SecureSession{
		SessionBaseImpl:    NewSessionBaseImpl(),
		mTable:             table,
		mState:             Establishing,
		mSecureSessionType: secureSessionType,
		mLocalSessionId:    localSessionId,
	}
}

func NewSecureSessionImplWithNodeId(
	table *SecureSessionTable,
	secureSessionType TypeSecureSession,
	localSessionId uint16,
	localNodeId lib.NodeId,
	peerNodeId lib.NodeId,
	peerCATs *lib.CATValues,
	peerSessionId uint16,
	fabric lib.FabricIndex,
	config *ReliableMessageProtocolConfig,
) *SecureSession {
	impl := &SecureSession{
		SessionBaseImpl:    NewSessionBaseImpl(),
		mTable:             table,
		mState:             Establishing,
		mSecureSessionType: secureSessionType,
		mLocalSessionId:    localSessionId,
		mLocalNodeId:       localNodeId,
		mPeerNodeId:        peerNodeId,
		mPeerSessionId:     peerSessionId,
		mRemoteMRPConfig:   config,
		mPeerCATs:          peerCATs,
	}
	impl.SetFabricIndex(fabric)
	return impl
}

func (s *SecureSession) Release() {
	//TODO implement me
	panic("implement me")
}

func (s *SecureSession) Retain() {
	//TODO implement me
	panic("implement me")
}

func (s *SecureSession) IsActiveSession() bool {
	return s.mState == Active
}

func (s *SecureSession) SessionType() uint8 {
	return uint8(Secure)
}

func (s *SecureSession) SessionTypeString() string {
	return "secure"
}

func (s *SecureSession) IsGroupSession() bool {
	return s.SessionType() == Secure.Uint8()
}

func (s *SecureSession) IsEstablishing() bool {
	return s.mState == Establishing
}

func (s *SecureSession) IsSecureSession() bool {
	return s.SessionType() == Secure.Uint8()
}

func (s *SecureSession) DispatchSessionEvent(delegate SessionDelegate) {
	//TODO implement me
	panic("implement me")
}

func (s *SecureSession) ComputeRoundTripTimeout(duration time.Duration) time.Duration {
	//TODO implement me
	panic("implement me")
}

func (s *SecureSession) SessionReleased() {
	//TODO implement me
	panic("implement me")
}

func (s *SecureSession) ClearValue() {
	//TODO implement me
	panic("implement me")
}

func (s *SecureSession) GetLocalSessionId() uint16 {
	return s.mLocalSessionId
}

func (s *SecureSession) IsDefunct() bool {
	return s.mState == Defunct
}

func (s *SecureSession) IsPendingEviction() bool {
	return s.mState == PendingEviction
}

func (s *SecureSession) GetStateStr() string {
	return s.mState.Str()
}

func (s *SecureSession) GetSecureSessionType() TypeSecureSession {
	return s.mSecureSessionType
}

func (s *SecureSession) GetCryptoContext() *CryptoContext {
	return s.mCryptoContext
}

func (s *SecureSession) GetPeerNodeId() lib.NodeId {
	return s.mPeerNodeId
}
