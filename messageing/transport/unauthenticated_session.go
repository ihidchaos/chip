package transport

import (
	"github.com/galenliu/chip/lib"
	"net/netip"
	"time"
)

type UnauthenticatedSession struct {
	*SessionBaseImpl
	mSessionRole              uint8
	mEphemeralInitiatorNodeId lib.NodeId
	mPeerAddress              netip.AddrPort
	mLastActivityTime         time.Time
	mLastPeerActivityTime     time.Time
	mRemoteMRPConfig          *ReliableMessageProtocolConfig
	mPeerMessageCounter       PeerMessageCounter
}

func NewUnauthenticatedSession(roleResponder uint8, id lib.NodeId, config *ReliableMessageProtocolConfig) *UnauthenticatedSession {
	return &UnauthenticatedSession{
		SessionBaseImpl:           NewSessionBaseImpl(),
		mSessionRole:              roleResponder,
		mEphemeralInitiatorNodeId: id,
		mPeerAddress:              netip.AddrPort{},
		mRemoteMRPConfig:          config,
		mPeerMessageCounter:       PeerMessageCounter{},
	}
}

func (s UnauthenticatedSession) ClearValue() {
	//TODO implement me
	panic("implement me")
}

func (s UnauthenticatedSession) SetPeerAddress(addr netip.AddrPort) {
	s.mPeerAddress = addr
}

func (s UnauthenticatedSession) AsUnauthenticatedSession() *UnauthenticatedSession {
	//TODO implement me
	panic("implement me")
}

func (s UnauthenticatedSession) Retain() {
	//TODO implement me
	panic("implement me")
}

func (s *UnauthenticatedSession) IsGroupSession() bool {
	//TODO implement me
	panic("implement me")
}

func (s UnauthenticatedSession) IsSecureSession() bool {
	return s.GetSessionType() == kSessionTypeSecure
}

func (s *UnauthenticatedSession) ComputeRoundTripTimeout(duration time.Duration) time.Duration {
	//TODO implement me
	panic("implement me")
}

func (s UnauthenticatedSession) GetSessionType() uint8 {
	return kSessionTypeUnauthenticated
}

func (s *UnauthenticatedSession) GetSessionTypeString() string {
	return "unauthenticated"
}
