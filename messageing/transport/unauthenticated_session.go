package transport

import (
	"github.com/galenliu/chip/lib"
	"net/netip"
	"time"
)

type UnauthenticatedSessionBase interface {
	Session
}

type UnauthenticatedSession struct {
	*SessionBaseImpl
	mSessionRole              TypeSessionRole
	mEphemeralInitiatorNodeId lib.NodeId
	mPeerAddress              netip.AddrPort
	mLastActivityTime         time.Time
	mLastPeerActivityTime     time.Time
	mRemoteMRPConfig          *ReliableMessageProtocolConfig
	mPeerMessageCounter       *PeerMessageCounter
}

func NewUnauthenticatedSessionImpl(roleResponder TypeSessionRole, id lib.NodeId, config *ReliableMessageProtocolConfig) *UnauthenticatedSession {
	return &UnauthenticatedSession{
		SessionBaseImpl:           NewSessionBaseImpl(),
		mSessionRole:              roleResponder,
		mEphemeralInitiatorNodeId: id,
		mPeerAddress:              netip.AddrPort{},
		mRemoteMRPConfig:          config,
		mPeerMessageCounter:       NewPeerMessageCounter(),
	}
}

func (s *UnauthenticatedSession) Release() {
	//TODO implement me
	panic("implement me")
}

func (s *UnauthenticatedSession) IsActiveSession() bool {
	//TODO implement me
	panic("implement me")
}

func (s *UnauthenticatedSession) IsEstablishing() bool {
	//TODO implement me
	panic("implement me")
}

func (s *UnauthenticatedSession) ClearValue() {
	//TODO implement me
	panic("implement me")
}

func (s *UnauthenticatedSession) SetPeerAddress(addr netip.AddrPort) {
	s.mPeerAddress = addr
}

func (s *UnauthenticatedSession) AsUnauthenticatedSession() *UnauthenticatedSession {
	//TODO implement me
	panic("implement me")
}

func (s *UnauthenticatedSession) Retain() {
	//TODO implement me
	panic("implement me")
}

func (s *UnauthenticatedSession) IsGroupSession() bool {
	//TODO implement me
	panic("implement me")
}

func (s *UnauthenticatedSession) IsSecureSession() bool {
	return s.SessionType() == Secure
}

func (s *UnauthenticatedSession) ComputeRoundTripTimeout(duration time.Duration) time.Duration {
	//TODO implement me
	panic("implement me")
}

func (s *UnauthenticatedSession) SessionType() SessionType {
	return Unauthenticated
}

func (s *UnauthenticatedSession) SessionTypeString() string {
	return "unauthenticated"
}

func (s *UnauthenticatedSession) LastActivityTime() time.Time {
	return s.mLastActivityTime
}
