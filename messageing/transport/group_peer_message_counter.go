package transport

import (
	"github.com/galenliu/chip/lib"
	"github.com/galenliu/chip/lib/store"
	"github.com/galenliu/chip/messageing/transport/session"
)

type GroupSender struct {
}

type GroupFabric struct {
}

type GroupPeerTable struct {
	mGroupFabrics []*GroupFabric
}

func NewGroupPeerTable(size int) *GroupPeerTable {
	return &GroupPeerTable{mGroupFabrics: make([]*GroupFabric, size)}
}

func (g *GroupPeerTable) FindOrAddPeer(fabricIndex lib.FabricIndex, nodeId lib.NodeId, isControl bool) (*session.PeerMessageCounter, error) {
	return nil, nil
}

func (g *GroupPeerTable) RemovePeer(fabricIndex lib.FabricIndex, nodeId lib.NodeId, isControl bool) error {
	return nil
}

func (g *GroupPeerTable) FabricRemoved(fabricIndex lib.FabricIndex) error {
	return nil
}

func (g *GroupPeerTable) removeSpecificPeer(list *GroupSender, nodeId lib.NodeId, size int) bool {
	return false
}

func (g *GroupPeerTable) compactPeers(list *GroupSender, size int) {

}

func (g *GroupPeerTable) removeAndCompactFabric(tableIndex int) {

}

type GroupOutgoingCounters struct {
	mStorage             store.PersistentStorageDelegate
	mGroupDataCounter    uint32
	mGroupControlCounter uint32
}

func NewGroupOutgoingCounters() *GroupOutgoingCounters {
	return &GroupOutgoingCounters{}
}

func (g *GroupOutgoingCounters) Init(storage store.PersistentStorageDelegate) error {
	g.mStorage = storage
	return nil
}

func (g *GroupOutgoingCounters) GetCounter(isControl bool) uint32 {
	if isControl {
		return g.mGroupControlCounter
	}
	return g.mGroupDataCounter
}

func (g *GroupOutgoingCounters) IncrementCounter(isControl bool) error {
	//key := lib.Uninitialized()
	//var value uint32 = 0
	//if isControl {
	//	g.mGroupControlCounter++
	//	value = g.mGroupControlCounter
	//}
	return nil
}
