package dnssd

import (
	"github.com/galenliu/chip/pkg/dnssd/responder"
	"time"
)

type QueryResponder struct {
	ResponderInfos []*QueryResponderInfo // TODO 数量需要做限定
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
	info := &QueryResponderInfo{
		Responder: res,
	}
	r.ResponderInfos = append(r.ResponderInfos, info)
	return NewQueryResponderSettings(info)
}
