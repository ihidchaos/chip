package transport

import (
	"github.com/galenliu/chip/lib"
	"github.com/galenliu/chip/messageing"
	"net/netip"
	"time"
)

type UnauthenticatedSession struct {
	Session
	mSessionRole              uint8
	mEphemeralInitiatorNodeId lib.NodeId

	mPeerAddress          netip.AddrPort
	mLastActivityTime     time.Duration
	mLastPeerActivityTime time.Duration
	mRemoteMRPConfig      messageing.ReliableMessageProtocolConfig
	mPeerMessageCounter   PeerMessageCounter
}

func (s UnauthenticatedSession) GetSessionRole() uint8 {
	return s.mSessionRole
}

func (s UnauthenticatedSession) GetEphemeralInitiatorNodeID() lib.NodeId {
	return s.mEphemeralInitiatorNodeId
}

const (
	kSessionRoleInitiator uint8 = iota
	kSessionRoleResponder
)

type UnauthenticatedSessionTable struct {
	mEntries []*UnauthenticatedSession
}

func (t UnauthenticatedSessionTable) FindOrAllocateResponder(source lib.NodeId, config *messageing.ReliableMessageProtocolConfig) SessionHandle {
	result := t.FindEntry(kSessionRoleResponder, source)
	if result != nil {
		SessionHandle()
	}
	return nil
}

func (t UnauthenticatedSessionTable) FindEntry(sessionRole uint8, ephemeralInitiatorNodeID lib.NodeId) *UnauthenticatedSession {
	for _, entry := range t.mEntries {
		if entry.GetSessionRole() == sessionRole && entry.GetEphemeralInitiatorNodeID() == ephemeralInitiatorNodeID {
			return entry
		}
	}
	return nil
}
