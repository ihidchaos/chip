package dnssd

//import (
//	"github.com/galenliu/chip/server/dnssd/responder"
//	"github.com/miekg/dns"
//)
//
//type ResponseBuilder struct {
//	mBuildOk bool
//	mMessage dns.Msg
//}
//
//func (b *ResponseBuilder) Reset() {
//	b.mMessage = dns.Msg{}
//	b.mBuildOk = true
//}
//
//func NewResponseBuilder() *ResponseBuilder {
//	return &ResponseBuilder{}
//}
//
//func (b *ResponseBuilder) AddQuery(query *Dnssd.QueryData) *ResponseBuilder {
//	if !b.mBuildOk {
//		return b
//	}
//	b.mMessage.SetReply(query.Msg)
//	return b
//}
//
//func (b *ResponseBuilder) BuildMessage() dns.Msg {
//	return b.mMessage
//}
//
//func (b *ResponseBuilder) AddRecord(resourceType int, record responder.RecordProvider) {
//	switch resourceType {
//	case ResourceType_Answer:
//		b.mMessage.Answer = append(b.mMessage.Answer, record.GetRecord())
//	case ResourceType_Additional:
//		b.mMessage.Extra = append(b.mMessage.Extra, record.GetRecord())
//	case ResourceType_Authority:
//		b.mMessage.Ns = append(b.mMessage.Ns, record.GetRecord())
//	}
//}
