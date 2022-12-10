package session

import (
	"github.com/galenliu/chip/lib"
	"github.com/galenliu/chip/messageing"
	"github.com/galenliu/chip/messageing/transport/raw"
	"sync"
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
	session.BaseImpl = &BaseImpl{
		locker:           sync.Mutex{},
		mFabricIndex:     index,
		mHolders:         nil,
		mSessionType:     kGroupOutgoing,
		mPeerAddress:     raw.PeerAddress{},
		base:             session,
		ReferenceCounted: lib.NewReferenceCounted(1, session),
	}
	return session
}

func (o *OutgoingGroup) IsActive() bool {
	//TODO implement me
	panic("implement me")
}

func (o *OutgoingGroup) IsEstablishing() bool {
	//TODO implement me
	panic("implement me")
}

func (o *OutgoingGroup) ComputeRoundTripTimeout(upperlayerProcessingTimeout time.Duration) time.Duration {
	if o.IsGroup() {
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

func (o *OutgoingGroup) GetPeer() lib.ScopedNodeId {
	return lib.NewScopedNodeId(lib.UndefinedNodeId(), o.FabricIndex())
}

func (o *OutgoingGroup) GroupId() lib.GroupId {
	return o.mGroupId
}

func (o *OutgoingGroup) RemoteMRPConfig() *messageing.ReliableMessageProtocolConfig {
	return messageing.DefaultMRPConfig()
}

func (o *OutgoingGroup) AckTimeout() time.Duration {
	return time.Duration(0)
}
