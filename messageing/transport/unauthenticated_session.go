package transport

import (
	"github.com/galenliu/chip/lib"
	log "golang.org/x/exp/slog"
	"net/netip"
	"time"
)

type UnauthenticatedSessionBase interface {
	Session
	MarkActiveRx()
	MarkActive()
	PeerAddress() netip.AddrPort
	SetPeerAddress(address netip.AddrPort)
}

type UnauthenticatedSession struct {
	*SessionBaseImpl
	mSessionRole              SessionRole
	mEphemeralInitiatorNodeId lib.NodeId
	mPeerAddress              netip.AddrPort
	mLastActivityTime         time.Time
	mLastPeerActivityTime     time.Time
	mRemoteMRPConfig          *ReliableMessageProtocolConfig
	mPeerMessageCounter       *PeerMessageCounter
}

func NewUnauthenticatedSession(roleResponder SessionRole, ephemeralInitiatorNodeID lib.NodeId, config *ReliableMessageProtocolConfig) *UnauthenticatedSession {
	session := &UnauthenticatedSession{
		mSessionRole:              roleResponder,
		mEphemeralInitiatorNodeId: ephemeralInitiatorNodeID,
		mPeerAddress:              netip.AddrPort{},
		mRemoteMRPConfig:          config,
		mPeerMessageCounter:       NewPeerMessageCounter(),
		mLastActivityTime:         time.Now(),
		mLastPeerActivityTime:     time.Time{},
	}
	session.SessionBaseImpl = NewSessionBaseImpl(1, kUnauthenticated, session)
	return session
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

func (s *UnauthenticatedSession) PeerAddress() netip.AddrPort {
	return s.mPeerAddress
}

func (s *UnauthenticatedSession) IsGroupSession() bool {
	//TODO implement me
	panic("implement me")
}

func (s *UnauthenticatedSession) ComputeRoundTripTimeout(duration time.Duration) time.Duration {
	//TODO implement me
	panic("implement me")
}

func (s *UnauthenticatedSession) SessionReleased() {
	//TODO implement me
	panic("implement me")
}

func (s *UnauthenticatedSession) DispatchSessionEvent(delegate SessionDelegate) {
	//TODO implement me
	panic("implement me")
}

func (s *UnauthenticatedSession) SessionRole() SessionRole {
	return s.mSessionRole
}

func (s *UnauthenticatedSession) EphemeralInitiatorNodeId() lib.NodeId {
	return s.mEphemeralInitiatorNodeId
}

func (s *UnauthenticatedSession) MarkActiveRx() {
	s.mLastPeerActivityTime = time.Now()
	s.MarkActive()
}

func (s *UnauthenticatedSession) MarkActive() {
	s.mLastActivityTime = time.Now()
}

func (s *UnauthenticatedSession) PeerMessageCounter() *PeerMessageCounter {
	return s.mPeerMessageCounter
}

func (s *UnauthenticatedSession) LastActivityTime() time.Time {
	return s.mLastActivityTime
}

func (s *UnauthenticatedSession) LastPeerActivityTime() time.Time {
	return s.mLastPeerActivityTime
}

func (s *UnauthenticatedSession) Released() {
	//TODO implement me
	panic("implement me")
}

func (s *UnauthenticatedSession) LogValue() log.Value {
	return log.GroupValue(
		log.String("SessionRole", s.mSessionRole.String()),
	)
}
