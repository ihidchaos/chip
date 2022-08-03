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

func (t UnauthenticatedSessionTable) FindOrAllocateResponder(ephemeralInitiatorNodeID lib.NodeId, config *ReliableMessageProtocolConfig) SessionHandle {
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
