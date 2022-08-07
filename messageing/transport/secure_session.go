package transport

import (
	"github.com/galenliu/chip/lib"
	"time"
)

type TSecureSession uint8

const (
	PASE TSecureSession = 1
	CASE TSecureSession = 2
)

const (
	KStateEstablishing = iota
	KStateActive
	KStateDefunct
	KStatePendingEviction
)

type SecureSession interface {
	Session
}

type SecureSessionImpl struct {
	*SessionBaseImpl
	mState             uint8
	mTable             *SecureSessionTable
	mSecureSessionType TSecureSession
	mLocalSessionId    uint16
	mPeerSessionId     uint16
	mLocalNodeId       lib.NodeId
	mPeerNodeId        lib.NodeId
	mRemoteMRPConfig   *ReliableMessageProtocolConfig
	mPeerCATs          *lib.CATValues
}

func NewSecureSessionImpl(
	table *SecureSessionTable,
	secureSessionType TSecureSession,
	localSessionId uint16,
) *SecureSessionImpl {
	return &SecureSessionImpl{
		SessionBaseImpl:    NewSessionBaseImpl(),
		mTable:             table,
		mState:             KStateEstablishing,
		mSecureSessionType: secureSessionType,
		mLocalSessionId:    localSessionId,
	}
}

func NewSecureSessionImplWithNodeId(
	table *SecureSessionTable,
	secureSessionType TSecureSession,
	localSessionId uint16,
	localNodeId lib.NodeId,
	peerNodeId lib.NodeId,
	peerCATs *lib.CATValues,
	peerSessionId uint16,
	fabric lib.FabricIndex,
	config *ReliableMessageProtocolConfig,
) *SecureSessionImpl {
	impl := &SecureSessionImpl{
		SessionBaseImpl:    NewSessionBaseImpl(),
		mTable:             table,
		mState:             KStateEstablishing,
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

func (s *SecureSessionImpl) Release() {
	//TODO implement me
	panic("implement me")
}

func (s *SecureSessionImpl) Retain() {
	//TODO implement me
	panic("implement me")
}

func (s *SecureSessionImpl) IsActiveSession() bool {
	return s.mState == KStateActive
}

func (s *SecureSessionImpl) GetSessionType() uint8 {
	return kSessionTypeSecure
}

func (s *SecureSessionImpl) GetSessionTypeString() string {
	return "secure"
}

func (s *SecureSessionImpl) IsGroupSession() bool {
	return s.GetSessionType() == kSessionTypeSecure
}

func (s *SecureSessionImpl) IsEstablishing() bool {
	return s.mState == KStateEstablishing
}

func (s *SecureSessionImpl) IsSecureSession() bool {
	return s.GetSessionType() == kSessionTypeSecure
}

func (s *SecureSessionImpl) DispatchSessionEvent(delegate SessionDelegate) {
	//TODO implement me
	panic("implement me")
}

func (s *SecureSessionImpl) ComputeRoundTripTimeout(duration time.Duration) time.Duration {
	//TODO implement me
	panic("implement me")
}

func (s *SecureSessionImpl) SessionReleased() {
	//TODO implement me
	panic("implement me")
}

func (s *SecureSessionImpl) AsUnauthenticatedSession() *UnauthenticatedSessionImpl {
	//TODO implement me
	panic("implement me")
}

func (s *SecureSessionImpl) ClearValue() {
	//TODO implement me
	panic("implement me")
}
