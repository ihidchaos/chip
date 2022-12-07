package session

import (
	"github.com/galenliu/chip/lib"
	"github.com/galenliu/chip/messageing"
	"time"
)

type IncomingGroupSession struct {
	*BaseImpl
	mGroupId    lib.GroupId
	mPeerNodeId lib.NodeId
}

func NewIncomingGroupSession(groupId lib.GroupId, index lib.FabricIndex, nodeId lib.NodeId) *IncomingGroupSession {
	session := &IncomingGroupSession{
		mGroupId:    groupId,
		mPeerNodeId: nodeId,
	}
	session.SetFabricIndex(index)
	session.BaseImpl = NewBaseImpl(1, kGroupIncoming, session)
	return session
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

func (i *IncomingGroupSession) GetPeer() lib.ScopedNodeId {
	return lib.NewScopedNodeId(i.mPeerNodeId, i.FabricIndex())
}

func (i *IncomingGroupSession) ComputeRoundTripTimeout(upperlayerProcessingTimeout time.Duration) time.Duration {
	if i.IsGroupSession() {
		return time.Duration(0)
	}
	return i.AckTimeout() + upperlayerProcessingTimeout
}

func (i *IncomingGroupSession) RemoteMRPConfig() *messageing.ReliableMessageProtocolConfig {
	config := messageing.DefaultMRPConfig()
	return &config
}

func (i *IncomingGroupSession) Released() {
	//TODO implement me
	panic("implement me")
}

func (i *IncomingGroupSession) ClearValue() {
	//TODO implement me
	panic("implement me")
}

func (i *IncomingGroupSession) AckTimeout() time.Duration {
	return time.Duration(0)
}
