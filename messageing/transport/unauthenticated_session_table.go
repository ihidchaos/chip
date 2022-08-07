package transport

import (
	"github.com/galenliu/chip/lib"
	"time"
)

const (
	kSessionRoleInitiator uint8 = iota
	kSessionRoleResponder
)

type UnauthenticatedSessionTable struct {
	mEntries []*UnauthenticatedSessionImpl
}

func (s *UnauthenticatedSessionImpl) SessionReleased() {
	//TODO implement me
	panic("implement me")
}

func (s *UnauthenticatedSessionImpl) DispatchSessionEvent(delegate SessionDelegate) {
	//TODO implement me
	panic("implement me")
}

func (s UnauthenticatedSessionImpl) GetSessionRole() uint8 {
	return s.mSessionRole
}

func (s *UnauthenticatedSessionImpl) GetEphemeralInitiatorNodeID() lib.NodeId {
	return s.mEphemeralInitiatorNodeId
}

func (s *UnauthenticatedSessionImpl) MarkActiveRx() {
	s.mLastPeerActivityTime = time.Now()
	s.MarkActive()
}

func (s UnauthenticatedSessionImpl) MarkActive() {
	s.mLastActivityTime = time.Now()
}

func (s UnauthenticatedSessionImpl) GetPeerMessageCounter() PeerMessageCounter {
	return s.mPeerMessageCounter
}

func (t UnauthenticatedSessionTable) FindOrAllocateResponder(ephemeralInitiatorNodeID lib.NodeId, config *ReliableMessageProtocolConfig) *UnauthenticatedSessionImpl {
	result := t.FindEntry(kSessionRoleResponder, ephemeralInitiatorNodeID)
	if result != nil {
		return result
	}
	entry := NewUnauthenticatedSessionImpl(kSessionRoleResponder, ephemeralInitiatorNodeID, config)
	t.mEntries = append(t.mEntries, entry)
	return entry
}

func (t UnauthenticatedSessionTable) FindEntry(sessionRole uint8, ephemeralInitiatorNodeID lib.NodeId) *UnauthenticatedSessionImpl {
	for _, entry := range t.mEntries {
		if entry.GetSessionRole() == sessionRole && entry.GetEphemeralInitiatorNodeID() == ephemeralInitiatorNodeID {
			return entry
		}
	}
	return nil
}

func (t UnauthenticatedSessionTable) FindInitiator(ephemeralInitiatorNodeID lib.NodeId) *UnauthenticatedSessionImpl {
	result := t.FindEntry(kSessionRoleInitiator, ephemeralInitiatorNodeID)
	if result != nil {
		return result
	}
	return nil
}
