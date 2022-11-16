package session

import (
	"github.com/galenliu/chip/lib"
	"time"
)

type OutgoingGroup struct {
	*BaseImpl
	mGroupId lib.GroupId
}

func NewOutgoingGroupSession(groupId lib.GroupId, index lib.FabricIndex) *OutgoingGroup {
	session := &OutgoingGroup{
		mGroupId: groupId,
	}
	session.SetFabricIndex(index)
	session.BaseImpl = NewBaseImpl(1, kGroupOutgoing, session)
	return session
}

func (o *OutgoingGroup) IsActiveSession() bool {
	//TODO implement me
	panic("implement me")
}

func (o *OutgoingGroup) IsEstablishing() bool {
	//TODO implement me
	panic("implement me")
}

func (o *OutgoingGroup) ComputeRoundTripTimeout(upperlayerProcessingTimeout time.Duration) time.Duration {
	if o.IsGroupSession() {
		return time.Duration(0)
	}
	return o.AckTimeout() + upperlayerProcessingTimeout
}

func (o *OutgoingGroup) Released() {
	//TODO implement me
	panic("implement me")
}

func (o *OutgoingGroup) ClearValue() {
	//TODO implement me
	panic("implement me")
}

func (s *OutgoingGroup) GetPeer() lib.ScopedNodeId {
	return lib.NewScopedNodeId(lib.UndefinedNodeId(), s.FabricIndex())
}

func (o *OutgoingGroup) GroupId() lib.GroupId {
	return o.mGroupId
}

func (o *OutgoingGroup) RemoteMRPConfig() *ReliableMessageProtocolConfig {
	return DefaultMRPConfig()
}

func (o *OutgoingGroup) AckTimeout() time.Duration {
	return time.Duration(0)
}
