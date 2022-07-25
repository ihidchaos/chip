package dnssd

import (
	"github.com/galenliu/chip/server/dnssd/responder"
)

type QueryResponderAllocator struct {
	mAllocatedResponders []responder.RecordResponder
	mResponderInfos      []*QueryResponderInfo
	mQueryResponder      *QueryResponder
}

func NewQueryResponderAllocator() *QueryResponderAllocator {
	return &QueryResponderAllocator{
		mAllocatedResponders: make([]responder.RecordResponder, 0),
		mResponderInfos:      make([]*QueryResponderInfo, 0),
		mQueryResponder:      NewQueryResponder(),
	}
}

func (a *QueryResponderAllocator) AddResponder(recordResponder responder.RecordResponder) *QueryResponderSettings {
	if len(a.mAllocatedResponders) >= MaxCommissionRecords {
		return &QueryResponderSettings{}
	}
	a.mAllocatedResponders = append(a.mAllocatedResponders, recordResponder)

	return a.mQueryResponder.AddResponder(recordResponder)

}

func (a *QueryResponderAllocator) GetResponder(typ uint16, name string) responder.RecordResponder {
	for _, r := range a.mAllocatedResponders {
		if r.GetType() == typ && r.GetName() == name {
			return r
		}
	}
	return nil
}

func (a *QueryResponderAllocator) GetQueryResponder() *QueryResponder {
	return a.mQueryResponder
}
