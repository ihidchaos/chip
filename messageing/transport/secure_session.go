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
	mState             TSecureSessionState
	mTable             *SecureSessionTable
	mSecureSessionType TSecureSessionType
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
	secureSessionType TSecureSessionType,
	localSessionId uint16,
) *SecureSession {
	return &SecureSession{
		SessionBaseImpl:    NewSessionBaseImpl(),
		mTable:             table,
		mState:             KEstablishing,
		mSecureSessionType: secureSessionType,
		mLocalSessionId:    localSessionId,
	}
}

func NewSecureSessionImplWithNodeId(
	table *SecureSessionTable,
	secureSessionType TSecureSessionType,
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
		mState:             KEstablishing,
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
	return s.mState == KActive
}

func (s *SecureSession) GetSessionType() uint8 {
	return uint8(kSecure)
}

func (s *SecureSession) GetSessionTypeString() string {
	return "secure"
}

func (s *SecureSession) IsGroupSession() bool {
	return s.GetSessionType() == kSecure.Uint8()
}

func (s *SecureSession) IsEstablishing() bool {
	return s.mState == KEstablishing
}

func (s *SecureSession) IsSecureSession() bool {
	return s.GetSessionType() == kSecure.Uint8()
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

func (s *SecureSession) AsUnauthenticatedSession() *UnauthenticatedSessionImpl {
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
	return s.mState == KDefunct
}

func (s *SecureSession) IsPendingEviction() bool {
	return s.mState == KPendingEviction
}

func (s *SecureSession) GetStateStr() string {
	return s.mState.Str()
}

func (s *SecureSession) GetSecureSessionType() TSecureSessionType {
	return s.mSecureSessionType
}

func (s *SecureSession) GetCryptoContext() *CryptoContext {
	return s.mCryptoContext
}

func (s *SecureSession) GetPeerNodeId() lib.NodeId {
	return s.mPeerNodeId
}
