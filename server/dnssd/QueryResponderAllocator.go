package DnssdServer

import "github.com/galenliu/dnssd/responder"

type QueryResponderAllocator struct {
	mAllocatedResponders []responder.RecordResponder
	mResponderInfos      []*responder.QueryResponderInfo
	mQueryResponder      *responder.QueryResponder
}

func NewQueryResponderAllocator() *QueryResponderAllocator {
	return &QueryResponderAllocator{
		mAllocatedResponders: make([]responder.RecordResponder, 0),
		mResponderInfos:      make([]*responder.QueryResponderInfo, 0),
		mQueryResponder:      responder.NewQueryResponder(),
	}
}

func (a *QueryResponderAllocator) AddResponder(recordResponder responder.RecordResponder) *responder.QueryResponderSettings {
	if len(a.mAllocatedResponders) >= kMaxCommissionRecords {
		return &responder.QueryResponderSettings{}
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

func (a *QueryResponderAllocator) GetQueryResponder() *responder.QueryResponder {
	return a.mQueryResponder
}
