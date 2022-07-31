package transport

import (
	"github.com/galenliu/chip/lib"
	"github.com/galenliu/chip/messageing"
	"net/netip"
	"time"
)

type UnauthenticatedSession struct {
	mSessionRole              uint8
	mEphemeralInitiatorNodeId lib.NodeId

	mPeerAddress          netip.AddrPort
	mLastActivityTime     time.Time
	mLastPeerActivityTime time.Time
	mRemoteMRPConfig      *messageing.ReliableMessageProtocolConfig
	mPeerMessageCounter   PeerMessageCounter
}

func (s UnauthenticatedSession) SetPeerAddress(addr netip.AddrPort) {
	//TODO implement me
	panic("implement me")
}

func (s UnauthenticatedSession) AsUnauthenticatedSession() *UnauthenticatedSession {
	//TODO implement me
	panic("implement me")
}

func (s UnauthenticatedSession) NotifySessionReleased() {
	//TODO implement me
	panic("implement me")
}

func (s UnauthenticatedSession) GetFabricIndex() {
	//TODO implement me
	panic("implement me")
}

func (s UnauthenticatedSession) SetFabricIndex(index lib.FabricIndex) {
	//TODO implement me
	panic("implement me")
}

func (s UnauthenticatedSession) Retain() {
	//TODO implement me
	panic("implement me")
}

func (s UnauthenticatedSession) IsGroupSession() bool {
	//TODO implement me
	panic("implement me")
}

func (s UnauthenticatedSession) IsSecureSession() bool {
	//TODO implement me
	panic("implement me")
}

func (s UnauthenticatedSession) ComputeRoundTripTimeout(duration time.Duration) time.Duration {
	//TODO implement me
	panic("implement me")
}

func (s UnauthenticatedSession) AddHolder(handle SessionHandle) {
	//TODO implement me
	panic("implement me")
}

func (s UnauthenticatedSession) RemoveHolder(handle SessionHandle) {
	//TODO implement me
	panic("implement me")
}

func NewUnauthenticatedSession(roleResponder uint8, id lib.NodeId, config *messageing.ReliableMessageProtocolConfig) SessionHandle {
	return &UnauthenticatedSession{
		mSessionRole:              roleResponder,
		mEphemeralInitiatorNodeId: id,
		mPeerAddress:              netip.AddrPort{},
		mRemoteMRPConfig:          config,
		mPeerMessageCounter:       PeerMessageCounter{},
	}
}

type UnauthenticatedSessionTable struct {
	mEntries []*UnauthenticatedSession
}

func (s *UnauthenticatedSession) SessionReleased() {
	//TODO implement me
	panic("implement me")
}

func (s *UnauthenticatedSession) DispatchSessionEvent(delegate SessionDelegate) {
	//TODO implement me
	panic("implement me")
}

func (s UnauthenticatedSession) GetSessionRole() uint8 {
	return s.mSessionRole
}

func (s *UnauthenticatedSession) GetEphemeralInitiatorNodeID() lib.NodeId {
	return s.mEphemeralInitiatorNodeId
}

func (s *UnauthenticatedSession) MarkActiveRx() {
	s.mLastPeerActivityTime = time.Now()
	s.MarkActive()
}

func (s UnauthenticatedSession) MarkActive() {
	s.mLastActivityTime = time.Now()
}

func (s UnauthenticatedSession) GetPeerMessageCounter() PeerMessageCounter {
	return s.mPeerMessageCounter
}

const (
	kSessionRoleInitiator uint8 = iota
	kSessionRoleResponder
)

func (t UnauthenticatedSessionTable) FindOrAllocateResponder(ephemeralInitiatorNodeID lib.NodeId, config *messageing.ReliableMessageProtocolConfig) SessionHandle {
	result := t.FindEntry(kSessionRoleResponder, ephemeralInitiatorNodeID)
	if result != nil {
		return result
	}
	return NewUnauthenticatedSession(kSessionRoleResponder, ephemeralInitiatorNodeID, config)
}

func (t UnauthenticatedSessionTable) FindEntry(sessionRole uint8, ephemeralInitiatorNodeID lib.NodeId) *UnauthenticatedSession {
	for _, entry := range t.mEntries {
		if entry.GetSessionRole() == sessionRole && entry.GetEphemeralInitiatorNodeID() == ephemeralInitiatorNodeID {
			return entry
		}
	}
	return nil
}

func (t UnauthenticatedSessionTable) FindInitiator(ephemeralInitiatorNodeID lib.NodeId) SessionHandle {
	result := t.FindEntry(kSessionRoleInitiator, ephemeralInitiatorNodeID)
	if result != nil {
		return result
	}
	return nil
}
