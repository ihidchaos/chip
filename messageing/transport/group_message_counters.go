package transport

import "github.com/galenliu/chip/pkg/storage"

type GroupOutgoingCounters struct {
	mStorage             storage.KvsPersistentStorageDelegate
	mGroupDataCounter    uint32
	mGroupControlCounter uint32
}

func NewGroupOutgoingCounters() *GroupOutgoingCounters {
	return &GroupOutgoingCounters{}
}

func (g *GroupOutgoingCounters) Init(storage storage.KvsPersistentStorageDelegate) error {
	g.mStorage = storage
	return nil
}

func (g *GroupOutgoingCounters) GetCounter(isControl bool) uint32 {
	if isControl {
		return g.mGroupControlCounter
	}
	return g.mGroupDataCounter
}
