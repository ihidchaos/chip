package transport

import (
	"github.com/galenliu/chip/lib"
	"time"
)

type TSessionRole uint8

const (
	KInitiator TSessionRole = iota
	KResponder
)

func (t TSessionRole) Uint8() uint8 {
	return uint8(t)
}

type UnauthenticatedSessionTable struct {
	mEntries []*UnauthenticatedSessionImpl
}

func NewUnauthenticatedSessionTable() *UnauthenticatedSessionTable {
	return &UnauthenticatedSessionTable{
		mEntries: make([]*UnauthenticatedSessionImpl, 0),
	}
}

func (s *UnauthenticatedSessionImpl) SessionReleased() {
	//TODO implement me
	panic("implement me")
}

func (s *UnauthenticatedSessionImpl) DispatchSessionEvent(delegate SessionDelegate) {
	//TODO implement me
	panic("implement me")
}

func (s *UnauthenticatedSessionImpl) GetSessionRole() TSessionRole {
	return s.mSessionRole
}

func (s *UnauthenticatedSessionImpl) GetEphemeralInitiatorNodeID() lib.NodeId {
	return s.mEphemeralInitiatorNodeId
}

func (s *UnauthenticatedSessionImpl) MarkActiveRx() {
	s.mLastPeerActivityTime = time.Now()
	s.MarkActive()
}

func (s *UnauthenticatedSessionImpl) MarkActive() {
	s.mLastActivityTime = time.Now()
}

func (s *UnauthenticatedSessionImpl) GetPeerMessageCounter() PeerMessageCounter {
	return s.mPeerMessageCounter
}

func (t *UnauthenticatedSessionTable) FindOrAllocateResponder(ephemeralInitiatorNodeID lib.NodeId, config *ReliableMessageProtocolConfig) *UnauthenticatedSessionImpl {
	result := t.FindEntry(KResponder, ephemeralInitiatorNodeID)
	if result != nil {
		return result
	}
	return t.AllocEntry(KResponder, ephemeralInitiatorNodeID, config)
}

func (t *UnauthenticatedSessionTable) AllocEntry(sessionRole TSessionRole, ephemeralInitiatorNodeID lib.NodeId, config *ReliableMessageProtocolConfig) *UnauthenticatedSessionImpl {

	entry := NewUnauthenticatedSessionImpl(sessionRole, ephemeralInitiatorNodeID, config)
	t.mEntries = append(t.mEntries, entry)
	return entry
}

func (t *UnauthenticatedSessionTable) FindEntry(sessionRole TSessionRole, ephemeralInitiatorNodeID lib.NodeId) *UnauthenticatedSessionImpl {
	for _, entry := range t.mEntries {
		if entry.GetSessionRole() == sessionRole && entry.GetEphemeralInitiatorNodeID() == ephemeralInitiatorNodeID {
			return entry
		}
	}
	return nil
}

func (t *UnauthenticatedSessionTable) FindInitiator(ephemeralInitiatorNodeID lib.NodeId) *UnauthenticatedSessionImpl {
	result := t.FindEntry(KInitiator, ephemeralInitiatorNodeID)
	if result != nil {
		return result
	}
	return nil
}
