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

type TSecureState uint8

const (
	KEstablishing TSecureState = iota
	KActive
	KDefunct
	KPendingEviction
)

func (t TSecureState) Str() string {
	return [...]string{
		"Establishing", "Active", "Defunct", "PendingEviction",
	}[t]
}

type SecureSessionBase interface {
	Session
}

type SecureSession struct {
	*SessionBaseImpl
	mState             TSecureState
	mTable             *SecureSessionTable
	mSecureSessionType TSecureSession
	mLocalSessionId    uint16
	mPeerSessionId     uint16
	mLocalNodeId       lib.NodeId
	mPeerNodeId        lib.NodeId
	mRemoteMRPConfig   *ReliableMessageProtocolConfig
	mPeerCATs          *lib.CATValues
}

func NewSecureSession(
	table *SecureSessionTable,
	secureSessionType TSecureSession,
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
	secureSessionType TSecureSession,
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
	return kSessionTypeSecure
}

func (s *SecureSession) GetSessionTypeString() string {
	return "secure"
}

func (s *SecureSession) IsGroupSession() bool {
	return s.GetSessionType() == kSessionTypeSecure
}

func (s *SecureSession) IsEstablishing() bool {
	return s.mState == KEstablishing
}

func (s *SecureSession) IsSecureSession() bool {
	return s.GetSessionType() == kSessionTypeSecure
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
