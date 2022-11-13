package transport

import (
	"github.com/galenliu/chip/config"
	"github.com/galenliu/chip/lib"
	"time"
)

type SessionRole uint8

const (
	RoleInitiator SessionRole = iota
	RoleResponder
)

func (t SessionRole) Uint8() uint8 {
	return uint8(t)
}

func (t SessionRole) String() string {
	switch t {
	case RoleResponder:
		return "Responder"
	case RoleInitiator:
		return "Initiator"
	default:
		return "unknown"
	}
}

type UnauthenticatedSessionTable struct {
	mEntries []*UnauthenticatedSession
	mMaxSize int
}

func NewUnauthenticatedSessionTable() *UnauthenticatedSessionTable {
	return &UnauthenticatedSessionTable{
		mEntries: make([]*UnauthenticatedSession, 0),
		mMaxSize: config.SecureSessionPoolSize,
	}
}

func (t *UnauthenticatedSessionTable) FindOrAllocateResponder(ephemeralInitiatorNodeId lib.NodeId, config *ReliableMessageProtocolConfig) (*SessionHandle, error) {
	var err error = nil
	result := t.findEntry(RoleResponder, ephemeralInitiatorNodeId)
	if result != nil {
		return NewSessionHandle(result), nil
	}
	result, err = t.allocEntry(RoleResponder, ephemeralInitiatorNodeId, config)
	if err != nil {
		return nil, err
	}
	return NewSessionHandle(result), nil
}

func (t *UnauthenticatedSessionTable) FindInitiator(ephemeralInitiatorNodeID lib.NodeId) *SessionHandle {
	result := t.findEntry(RoleInitiator, ephemeralInitiatorNodeID)
	if result != nil {
		return NewSessionHandle(result)
	}
	return nil
}

func (t *UnauthenticatedSessionTable) findLeastRecentUsedEntry() *UnauthenticatedSession {
	var result *UnauthenticatedSession = nil
	var oldTime = time.Now()
	for _, e := range t.mEntries {
		if e.ReferenceCount() == 0 && e.LastPeerActivityTime().After(oldTime) {
			oldTime = e.LastActivityTime()
			result = e
		}
	}
	return result
}

func (t *UnauthenticatedSessionTable) allocEntry(sessionRole SessionRole, ephemeralInitiatorNodeID lib.NodeId, config *ReliableMessageProtocolConfig) (*UnauthenticatedSession, error) {
	var entry *UnauthenticatedSession
	if len(t.mEntries) < t.mMaxSize {
		entry = NewUnauthenticatedSession(sessionRole, ephemeralInitiatorNodeID, config)
		t.mEntries = append(t.mEntries, entry)
		return entry, nil
	}
	entry = t.findLeastRecentUsedEntry()
	if entry == nil {
		return nil, lib.NotMemory
	}
	entry = NewUnauthenticatedSession(sessionRole, ephemeralInitiatorNodeID, config)
	return entry, nil
}

func (t *UnauthenticatedSessionTable) findEntry(sessionRole SessionRole, ephemeralInitiatorNodeId lib.NodeId) *UnauthenticatedSession {
	for _, entry := range t.mEntries {
		if entry.SessionRole() == sessionRole && entry.EphemeralInitiatorNodeId() == ephemeralInitiatorNodeId {
			return entry
		}
	}
	return nil
}
