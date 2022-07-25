package dnssd

import "github.com/galenliu/chip/server/dnssd/responder"

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

func (r *QueryResponder) AddResponder(res responder.RecordResponder) *QueryResponderSettings {
	info := &QueryResponderInfo{
		Responder: res,
	}
	r.ResponderInfos = append(r.ResponderInfos, info)
	return NewQueryResponderSettings(info)
}

//func (r *QueryResponder) AddAllResponses(source *IPPacket.Info, delegate Delegate, configuration *ResponseConfiguration) {
//for _, m := range r.ResponderInfos {
//	if !m.reportService {
//		continue
//	}
//	if m.RecordResponder == nil {
//		continue
//	}
//	r := record.NewPtrResourceRecord("", m.GetName())
//	configuration.Adjust(r)
//	delegate.AddResponse(r)
//
//}
//}
//
//func (r *QueryResponder) MarkAdditionalRepliesFor(info *QueryResponderInfo) {
//	if !info.alsoReportAdditionalQName {
//		return
//	}
//	if r.markAdditional(info.additionalQName) == 0 {
//		return
//	}
//	var keepAdding = true
//	for keepAdding {
//		keepAdding = false
//		var filter = QueryResponderRecordFilter{}
//		filter.SetIncludeAdditionalRepliesOnly(true)
//		for _, i := range r.ResponderInfos {
//			if i.alsoReportAdditionalQName {
//				keepAdding = keepAdding || r.markAdditional(i.additionalQName) != 0
//			}
//		}
//
//	}
//
//}
//
//func (r *QueryResponder) markAdditional(name core.FullQName) int {
//	var count = 0
//	for _, info := range r.ResponderInfos {
//		if info.reportNowAsAdditional {
//			continue
//		}
//		if info.RecordResponder == nil {
//			continue
//		}
//		if info.GetQName() == r.Msg.MsgHdr.Id {
//			info.reportNowAsAdditional = true
//			count++
//		}
//	}
//	return count
//}
