package transport

import (
	"github.com/galenliu/chip/lib"
	"time"
)

type OutgoingGroupSession struct {
	*SessionBaseImpl
	mGroupId lib.GroupId
}

type IncomingGroupSession struct {
	*SessionBaseImpl
	mGroupId    lib.GroupId
	mPeerNodeId lib.NodeId
}

func (o *OutgoingGroupSession) IsActiveSession() bool {
	//TODO implement me
	panic("implement me")
}

func (o *OutgoingGroupSession) IsGroupSession() bool {
	//TODO implement me
	panic("implement me")
}

func (o *OutgoingGroupSession) IsEstablishing() bool {
	//TODO implement me
	panic("implement me")
}

func (o *OutgoingGroupSession) ComputeRoundTripTimeout(duration time.Duration) time.Duration {
	//TODO implement me
	panic("implement me")
}

func (o *OutgoingGroupSession) Released() {
	//TODO implement me
	panic("implement me")
}

func (o *OutgoingGroupSession) ClearValue() {
	//TODO implement me
	panic("implement me")
}

func (s *OutgoingGroupSession) GetPeer() lib.ScopedNodeId {
	return lib.NewScopedNodeId(lib.UndefinedNodeId, s.FabricIndex())
}

func (o *OutgoingGroupSession) GroupId() lib.GroupId {
	return o.mGroupId
}

func (i *IncomingGroupSession) IsActiveSession() bool {
	//TODO implement me
	panic("implement me")
}

func (i *IncomingGroupSession) IsGroupSession() bool {
	//TODO implement me
	panic("implement me")
}

func (i *IncomingGroupSession) IsEstablishing() bool {
	//TODO implement me
	panic("implement me")
}

func (s *IncomingGroupSession) GetPeer() lib.ScopedNodeId {
	return lib.NewScopedNodeId(s.mPeerNodeId, s.FabricIndex())
}

func (i *IncomingGroupSession) ComputeRoundTripTimeout(duration time.Duration) time.Duration {
	//TODO implement me
	panic("implement me")
}

func (i *IncomingGroupSession) Released() {
	//TODO implement me
	panic("implement me")
}

func (i *IncomingGroupSession) ClearValue() {
	//TODO implement me
	panic("implement me")
}

func NewIncomingGroupSession(groupId lib.GroupId, index lib.FabricIndex, nodeId lib.NodeId) *IncomingGroupSession {
	session := &IncomingGroupSession{
		mGroupId:    groupId,
		mPeerNodeId: nodeId,
	}
	session.SetFabricIndex(index)

	session.SessionBaseImpl = NewSessionBaseImpl(1, kGroupIncoming, session)
	return session
}

func NewOutgoingGroupSession(groupId lib.GroupId, index lib.FabricIndex) *OutgoingGroupSession {
	session := &OutgoingGroupSession{
		mGroupId: groupId,
	}
	session.SetFabricIndex(index)
	session.SessionBaseImpl = NewSessionBaseImpl(1, kGroupIncoming, session)
	return session
}
