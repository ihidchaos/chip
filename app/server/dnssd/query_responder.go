package dnssd

import (
	"github.com/galenliu/chip/app/server/dnssd/responder"
	"time"
)

type QueryResponder struct {
	ResponderInfos []*QueryResponderInfo // TODO
}

func NewQueryResponder() *QueryResponder {
	return &QueryResponder{ResponderInfos: make([]*QueryResponderInfo, 0)}
}

func (r *QueryResponder) ResetAdditionals() {
	for _, r := range r.ResponderInfos {
		r.reportNowAsAdditional = false
	}
}

func (r *QueryResponder) ClearBroadcastThrottle() {
	for _, r := range r.ResponderInfos {
		r.lastMulticastTime = time.Time{}
	}
}

func (r *QueryResponder) AddResponder(res responder.RecordResponder) *QueryResponderSettings {
	info := NewQueryResponderInfo(res)
	r.ResponderInfos = append(r.ResponderInfos, info)
	return NewQueryResponderSettings(info)
}
