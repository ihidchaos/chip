package transport

import (
	"github.com/galenliu/chip/config"
	"github.com/galenliu/chip/lib"
	"github.com/galenliu/chip/messageing/transport/session"
	"time"
)

type UnauthenticatedSessionTable struct {
	mEntries [config.SecureSessionPoolSize]*session.Unauthenticated
}

func NewUnauthenticatedSessionTable() *UnauthenticatedSessionTable {
	return &UnauthenticatedSessionTable{
		mEntries: [config.SecureSessionPoolSize]*session.Unauthenticated{},
	}
}

func (t *UnauthenticatedSessionTable) FindOrAllocateResponder(ephemeralInitiatorNodeId lib.NodeId, config *session.ReliableMessageProtocolConfig) (*SessionHandle, error) {
	var err error = nil
	result := t.findEntry(session.Responder, ephemeralInitiatorNodeId)
	if result != nil {
		return NewSessionHandle(result), nil
	}
	result, err = t.allocEntry(session.Responder, ephemeralInitiatorNodeId, config)
	if err != nil {
		return nil, err
	}
	return NewSessionHandle(result), nil
}

func (t *UnauthenticatedSessionTable) FindInitiator(ephemeralInitiatorNodeID lib.NodeId) *SessionHandle {
	result := t.findEntry(session.Initiator, ephemeralInitiatorNodeID)
	if result != nil {
		return NewSessionHandle(result)
	}
	return nil
}

func (t *UnauthenticatedSessionTable) findLeastRecentUsedEntry() *session.Unauthenticated {
	var result *session.Unauthenticated = nil
	var oldTime = time.Now()
	for _, e := range t.mEntries {
		if e.ReferenceCount() == 0 && e.LastPeerActivityTime().After(oldTime) {
			oldTime = e.LastActivityTime()
			result = e
		}
	}
	return result
}

func (t *UnauthenticatedSessionTable) allocEntry(sessionRole session.Role, ephemeralInitiatorNodeID lib.NodeId, config *session.ReliableMessageProtocolConfig) (*session.Unauthenticated, error) {
	var entry *session.Unauthenticated
	entry = t.findLeastRecentUsedEntry()
	if entry == nil {
		return nil, lib.NotMemory
	}
	entry = session.NewUnauthenticated(sessionRole, ephemeralInitiatorNodeID, config)
	return entry, nil
}

func (t *UnauthenticatedSessionTable) findEntry(sessionRole session.Role, ephemeralInitiatorNodeId lib.NodeId) *session.Unauthenticated {
	for _, entry := range t.mEntries {
		if entry.SessionRole() == sessionRole && entry.EphemeralInitiatorNodeId() == ephemeralInitiatorNodeId {
			return entry
		}
	}
	return nil
}
