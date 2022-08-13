package transport

import (
	"github.com/galenliu/chip/lib"
	"net/netip"
	"time"
)

type UnauthenticatedSession interface {
	Session
}

type UnauthenticatedSessionImpl struct {
	*SessionBaseImpl
	mSessionRole              TSessionRole
	mEphemeralInitiatorNodeId lib.NodeId
	mPeerAddress              netip.AddrPort
	mLastActivityTime         time.Time
	mLastPeerActivityTime     time.Time
	mRemoteMRPConfig          *ReliableMessageProtocolConfig
	mPeerMessageCounter       *PeerMessageCounter
}

func NewUnauthenticatedSessionImpl(roleResponder TSessionRole, id lib.NodeId, config *ReliableMessageProtocolConfig) *UnauthenticatedSessionImpl {
	return &UnauthenticatedSessionImpl{
		SessionBaseImpl:           NewSessionBaseImpl(),
		mSessionRole:              roleResponder,
		mEphemeralInitiatorNodeId: id,
		mPeerAddress:              netip.AddrPort{},
		mRemoteMRPConfig:          config,
		mPeerMessageCounter:       NewPeerMessageCounter(),
	}
}

func (s *UnauthenticatedSessionImpl) Release() {
	//TODO implement me
	panic("implement me")
}

func (s *UnauthenticatedSessionImpl) IsActiveSession() bool {
	//TODO implement me
	panic("implement me")
}

func (s *UnauthenticatedSessionImpl) IsEstablishing() bool {
	//TODO implement me
	panic("implement me")
}

func (s *UnauthenticatedSessionImpl) ClearValue() {
	//TODO implement me
	panic("implement me")
}

func (s *UnauthenticatedSessionImpl) SetPeerAddress(addr netip.AddrPort) {
	s.mPeerAddress = addr
}

func (s *UnauthenticatedSessionImpl) AsUnauthenticatedSession() *UnauthenticatedSessionImpl {
	//TODO implement me
	panic("implement me")
}

func (s *UnauthenticatedSessionImpl) Retain() {
	//TODO implement me
	panic("implement me")
}

func (s *UnauthenticatedSessionImpl) IsGroupSession() bool {
	//TODO implement me
	panic("implement me")
}

func (s *UnauthenticatedSessionImpl) IsSecureSession() bool {
	return s.GetSessionType() == kSecure.Uint8()
}

func (s *UnauthenticatedSessionImpl) ComputeRoundTripTimeout(duration time.Duration) time.Duration {
	//TODO implement me
	panic("implement me")
}

func (s *UnauthenticatedSessionImpl) GetSessionType() uint8 {
	return kUnauthenticated.Uint8()
}

func (s *UnauthenticatedSessionImpl) GetSessionTypeString() string {
	return "unauthenticated"
}

func (s *UnauthenticatedSessionImpl) GetLastActivityTime() time.Time {
	return s.mLastActivityTime
}
