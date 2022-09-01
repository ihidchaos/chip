package transport

import (
	"github.com/galenliu/chip/config"
	"github.com/galenliu/chip/lib"
	"time"
)

type TypeSessionRole uint8

const (
	RoleInitiator TypeSessionRole = iota
	RoleResponder
)

func (t TypeSessionRole) Uint8() uint8 {
	return uint8(t)
}

type UnauthenticatedSessionTable struct {
	mEntries []*UnauthenticatedSession
	mMaxSize int
}

func NewUnauthenticatedSessionTable() *UnauthenticatedSessionTable {
	return &UnauthenticatedSessionTable{
		mEntries: make([]*UnauthenticatedSession, 0),
		mMaxSize: config.ChipConfigSecureSessionPoolSize,
	}
}

func (s *UnauthenticatedSession) SessionReleased() {
	//TODO implement me
	panic("implement me")
}

func (s *UnauthenticatedSession) DispatchSessionEvent(delegate SessionDelegate) {
	//TODO implement me
	panic("implement me")
}

func (s *UnauthenticatedSession) SessionRole() TypeSessionRole {
	return s.mSessionRole
}

func (s *UnauthenticatedSession) GetEphemeralInitiatorNodeID() lib.NodeId {
	return s.mEphemeralInitiatorNodeId
}

func (s *UnauthenticatedSession) MarkActiveRx() {
	s.mLastPeerActivityTime = time.Now()
	s.MarkActive()
}

func (s *UnauthenticatedSession) MarkActive() {
	s.mLastActivityTime = time.Now()
}

func (s *UnauthenticatedSession) GetPeerMessageCounter() *PeerMessageCounter {
	return s.mPeerMessageCounter
}

func (t *UnauthenticatedSessionTable) FindOrAllocateResponder(ephemeralInitiatorNodeID lib.NodeId, config *ReliableMessageProtocolConfig) *UnauthenticatedSession {
	result := t.FindEntry(RoleResponder, ephemeralInitiatorNodeID)
	if result != nil {
		return result
	}
	return t.AllocEntry(RoleResponder, ephemeralInitiatorNodeID, config)
}

func (t *UnauthenticatedSessionTable) AllocEntry(sessionRole TypeSessionRole, ephemeralInitiatorNodeID lib.NodeId, config *ReliableMessageProtocolConfig) *UnauthenticatedSession {
	entry := NewUnauthenticatedSessionImpl(sessionRole, ephemeralInitiatorNodeID, config)
	if len(t.mEntries) < t.mMaxSize {
		t.mEntries = append(t.mEntries, entry)
	}
	return t.FindLeastRecentUsedEntry(entry)
}

func (t *UnauthenticatedSessionTable) FindEntry(sessionRole TypeSessionRole, ephemeralInitiatorNodeID lib.NodeId) *UnauthenticatedSession {
	for _, entry := range t.mEntries {
		if entry.SessionRole() == sessionRole && entry.GetEphemeralInitiatorNodeID() == ephemeralInitiatorNodeID {
			return entry
		}
	}
	return nil
}

func (t *UnauthenticatedSessionTable) FindInitiator(ephemeralInitiatorNodeID lib.NodeId) *UnauthenticatedSession {
	result := t.FindEntry(RoleInitiator, ephemeralInitiatorNodeID)
	return result
}

func (t *UnauthenticatedSessionTable) FindLeastRecentUsedEntry(entry *UnauthenticatedSession) *UnauthenticatedSession {
	var oldTime time.Time
	var result *UnauthenticatedSession
	for _, e := range t.mEntries {
		if oldTime.IsZero() {
			result = e
			oldTime = e.LastActivityTime()
		}
		if e.LastActivityTime().After(oldTime) {
			oldTime = e.LastActivityTime()
			result = e
		}
	}
	result = entry
	return result
}
