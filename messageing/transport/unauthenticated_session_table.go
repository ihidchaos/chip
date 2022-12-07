package transport

import (
	"github.com/galenliu/chip/config"
	"github.com/galenliu/chip/lib"
	"github.com/galenliu/chip/messageing"
	"github.com/galenliu/chip/messageing/transport/session"
	"time"
)

type UnauthenticatedSessionTable struct {
	mEntries [config.UnauthenticatedConnectionPoolSize]*session.Unauthenticated
}

func NewUnauthenticatedSessionTable() *UnauthenticatedSessionTable {
	return &UnauthenticatedSessionTable{
		mEntries: [config.UnauthenticatedConnectionPoolSize]*session.Unauthenticated{},
	}
}

func (t *UnauthenticatedSessionTable) FindOrAllocateResponder(ephemeralInitiatorNodeId lib.NodeId, config *messageing.ReliableMessageProtocolConfig) (*SessionHandle, error) {
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
	var index = -1
	var oldTime = time.Now()
	for i, e := range t.mEntries {
		if e != nil && e.ReferenceCount() == 0 && e.LastPeerActivityTime().After(oldTime) {
			oldTime = e.LastActivityTime()
			index = i
		}
	}
	if index >= 0 && index < len(t.mEntries) {
		return t.mEntries[index]
	}
	return nil
}

func (t *UnauthenticatedSessionTable) allocEntry(sessionRole session.Role, ephemeralInitiatorNodeID lib.NodeId, config *messageing.ReliableMessageProtocolConfig) (*session.Unauthenticated, error) {

	entry := t.createEntry(sessionRole, ephemeralInitiatorNodeID, config)
	if entry != nil {
		return entry, nil
	}
	entry = t.findLeastRecentUsedEntry()
	if entry == nil {
		return nil, lib.NotMemory
	}
	*entry = *session.NewUnauthenticated(sessionRole, ephemeralInitiatorNodeID, config)
	return entry, nil
}

func (t *UnauthenticatedSessionTable) createEntry(role session.Role, ephemeralInitiatorNodeID lib.NodeId, config *messageing.ReliableMessageProtocolConfig) *session.Unauthenticated {
	for i, e := range t.mEntries {
		if e == nil {
			t.mEntries[i] = session.NewUnauthenticated(role, ephemeralInitiatorNodeID, config)
			return t.mEntries[i]
		}
	}
	return nil
}

func (t *UnauthenticatedSessionTable) findEntry(sessionRole session.Role, ephemeralInitiatorNodeId lib.NodeId) *session.Unauthenticated {
	for i, entry := range t.mEntries {
		if entry != nil && entry.Role() == sessionRole && entry.EphemeralInitiatorNodeId() == ephemeralInitiatorNodeId {
			return t.mEntries[i]
		}
	}
	return nil
}
