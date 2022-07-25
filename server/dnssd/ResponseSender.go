package dnssd

import (
	"github.com/galenliu/chip/server/dnssd/responder"
	"github.com/miekg/dns"
	log "github.com/sirupsen/logrus"
	"net"
	"time"
)

type Server interface {
	Broadcast(msg *dns.Msg) error
}

// ResponseSender 实现 ResponderDelegate接口
type ResponseSender struct {
	mSendState  *ResponseSendingState
	mResponders []*QueryResponder
	mServer     MdnsServer
}

func NewResponseSender() *ResponseSender {
	return &ResponseSender{
		mResponders: nil,
		mServer:     NewMdnsServerImpl(),
	}
}

func (rs *ResponseSender) Int() *ResponseSender {
	rs.mResponders = make([]*QueryResponder, 0)
	rs.mServer = NewMdnsServerImpl()

	rs.mSendState = &ResponseSendingState{
		mQuery:        nil,
		mMessageId:    0,
		mResourceType: 0,
		mSendError:    nil,
	}
	return rs
}

//func (rs *ResponseSender) BroadcastRecords(query *QueryData, client *dns.Client, address string) error {
//
//	msg, err := rs.OnQuery(query, nil)
//	if err != nil {
//		return err
//	}
//	if client == nil {
//		return fmt.Errorf("dns clint empty")
//	}
//	log.Printf("mDns broadcast Msg: %s", msg.String())
//	log.Printf("mDns broadcast Net: %s ,Local Addr: %s Des Addr: %s", client.Net, client.Dialer.LocalAddr.String(), address)
//	_, _, err = client.Exchange(msg, address)
//	if err != nil {
//		return err
//	}
//	return nil
//}

func (rs *ResponseSender) Respond(query *QueryData, configuration *responder.ResponseConfiguration, interfaceId net.Interface) error {

	for _, res := range rs.mResponders {
		res.ResetAdditionals()
	}
	log.Infof("Query Message RespondResponse:\t\n %s ", query.Msg.String())

	msg, err := rs.OnQuery(query, configuration)
	if err != nil {
		err := rs.mServer.SendTo(msg, query.mDestAddr, interfaceId)
		if err != nil {
			return err
		}
	}
	return nil
}

func (rs *ResponseSender) OnQuery(query *QueryData, configuration *responder.ResponseConfiguration) (*dns.Msg, error) {
	msg := new(dns.Msg)
	msg.SetReply(query.Msg)

	{
		queryReplyFilter := NewQueryReplyFilter(query)
		responseFilter := NewQueryResponderRecordFilter().
			SetReplyFilter(queryReplyFilter)
		if !rs.mSendState.SendUnicast() {
			responseFilter.SetIncludeOnlyMulticastBeforeMS(time.Now().Add(-1 * time.Second))
		}
		for _, resp := range rs.mResponders {
			for _, info := range resp.ResponderInfos {
				if !responseFilter.Accept(info) {
					continue
				}
				record := info.Responder.GetRecord()
				if configuration != nil {
					configuration.Adjust(record)
				}
				msg.Answer = append(msg.Answer, info.Responder.GetRecord())
			}
		}

	}
	// send all 'Additional' replies
	{
		rs.mSendState.SetResourceType(ResourceType_Additional)

		queryReplyFilter := NewQueryReplyFilter(query)
		queryReplyFilter.SetIgnoreNameMatch(true).
			SetSendingAdditionalItems(true)

		responseFilter := NewQueryResponderRecordFilter().
			SetReplyFilter(queryReplyFilter).
			SetIncludeAdditionalRepliesOnly(true)

		for _, resp := range rs.mResponders {
			for _, info := range resp.ResponderInfos {
				if !responseFilter.Accept(info) {
					continue
				}
				record := info.Responder.GetRecord()
				if configuration != nil {
					configuration.Adjust(record)
				}
				msg.Extra = append(msg.Answer, record)
			}
		}
	}
	return msg, nil
}

func (rs *ResponseSender) SetServer(server MdnsServer) {
	rs.mServer = server
}

func (rs *ResponseSender) AddQueryResponder(queryResponder *QueryResponder) error {
	rs.mResponders = append(rs.mResponders, queryResponder)
	return nil
}
