package session

import (
	"github.com/galenliu/chip/lib"
	"github.com/galenliu/chip/messageing"
	"github.com/galenliu/chip/messageing/transport/raw"
	"net"
	"net/netip"
	"sync"
	"time"
)

type Unauthenticated struct {
	*BaseImpl
	mSessionRole              Role
	mEphemeralInitiatorNodeId lib.NodeId
	mPeerAddress              netip.AddrPort
	mLastActivityTime         time.Time
	mLastPeerActivityTime     time.Time

	mRemoteMRPConfig    *messageing.ReliableMessageProtocolConfig
	mPeerMessageCounter *PeerMessageCounter
}

func NewUnauthenticated(roleResponder Role, ephemeralInitiatorNodeID lib.NodeId, config *messageing.ReliableMessageProtocolConfig) *Unauthenticated {
	session := &Unauthenticated{
		mSessionRole:              roleResponder,
		mEphemeralInitiatorNodeId: ephemeralInitiatorNodeID,
		mPeerAddress:              netip.AddrPort{},
		mRemoteMRPConfig:          config,
		mPeerMessageCounter:       NewPeerMessageCounter(),
		mLastActivityTime:         time.Now(),
		mLastPeerActivityTime:     time.Time{},
	}
	session.BaseImpl = &BaseImpl{
		locker:           sync.Mutex{},
		mFabricIndex:     lib.UndefinedFabricIndex(),
		mHolders:         nil,
		mSessionType:     kUnauthenticated,
		mPeerAddress:     raw.PeerAddress{},
		base:             session,
		ReferenceCounted: lib.NewReferenceCounted(1, session),
	}
	return session
}

func (s *Unauthenticated) GetPeer() lib.ScopedNodeId {
	return lib.NewScopedNodeId(lib.UndefinedNodeId(), s.mFabricIndex)
}

func (s *Unauthenticated) IsActive() bool {
	//TODO implement me
	panic("implement me")
}

func (s *Unauthenticated) IsEstablishing() bool {
	//TODO implement me
	panic("implement me")
}

func (s *Unauthenticated) ClearValue() {
	//TODO implement me
	panic("implement me")
}

func (s *Unauthenticated) SetPeerAddress(addr netip.AddrPort) {
	s.mPeerAddress = addr
}

func (s *Unauthenticated) PeerAddress() raw.PeerAddress {
	return raw.PeerAddress{
		TransportType: 0,
		Interface:     net.Interface{},
		AddrPort:      s.mPeerAddress,
	}
}

func (s *Unauthenticated) PeerNodeId() lib.NodeId {
	if s.mSessionRole == RoleInitiator {
		return lib.UndefinedNodeId()
	}
	return s.mEphemeralInitiatorNodeId
}

func (s *Unauthenticated) AckTimeout() time.Duration {
	switch s.BaseImpl.mPeerAddress.TransportType {
	case raw.Udp:
		return messageing.GetRetransmissionTimeout(s.mRemoteMRPConfig.ActiveRetransTimeout,
			s.mRemoteMRPConfig.IdleRetransTimeout, s.mLastPeerActivityTime, 1)
	case raw.Tcp:
		return 30 * time.Second
	default:
		return time.Duration(0)
	}
}

func (s *Unauthenticated) ComputeRoundTripTimeout(upperlayerProcessingTimeout time.Duration) time.Duration {
	if s.IsGroup() {
		return time.Duration(0)
	}
	return s.AckTimeout() + upperlayerProcessingTimeout
}

func (s *Unauthenticated) SessionReleased() {
	//TODO implement me
	panic("implement me")
}

func (s *Unauthenticated) DispatchSessionEvent(delegate DelegateEvent) {
	//TODO implement me
	panic("implement me")
}

func (s *Unauthenticated) Role() Role {
	return s.mSessionRole
}

func (s *Unauthenticated) EphemeralInitiatorNodeId() lib.NodeId {
	return s.mEphemeralInitiatorNodeId
}

func (s *Unauthenticated) MarkActiveRx() {
	s.mLastPeerActivityTime = time.Now()
	s.MarkActive()
}

func (s *Unauthenticated) MarkActive() {
	s.mLastActivityTime = time.Now()
}

func (s *Unauthenticated) PeerMessageCounter() *PeerMessageCounter {
	return s.mPeerMessageCounter
}

func (s *Unauthenticated) LastActivityTime() time.Time {
	return s.mLastActivityTime
}

func (s *Unauthenticated) LastPeerActivityTime() time.Time {
	return s.mLastPeerActivityTime
}

func (s *Unauthenticated) Released() {
	//TODO implement me
	panic("implement me")
}

func (s *Unauthenticated) RemoteMRPConfig() *messageing.ReliableMessageProtocolConfig {
	return s.mRemoteMRPConfig
}

func (s *Unauthenticated) SetRemoteMRPConfig(config *messageing.ReliableMessageProtocolConfig) {
	s.mRemoteMRPConfig = config
}
