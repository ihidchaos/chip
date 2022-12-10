package dnssd

import (
	"github.com/miekg/dns"
	log "golang.org/x/exp/slog"
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

func (rs *ResponseSender) ServeMdns(writer ResponseWriter, data *QueryData) error {
	msg, err := rs.OnQuery(data)
	if err != nil {
		return err
	}
	return writer.WriteMsg(msg)
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

func (rs *ResponseSender) Respond(query *QueryData, interfaceId net.Interface) error {

	for _, res := range rs.mResponders {
		res.ResetAdditionals()
	}
	msg, err := rs.OnQuery(query)
	if err != nil {
		err := rs.mServer.SendTo(msg, query.mDestAddr, interfaceId)
		if err != nil {
			return err
		}
	}
	return nil
}

func (rs *ResponseSender) OnQuery(query *QueryData) (*dns.Msg, error) {
	msg := new(dns.Msg)
	msg.SetReply(query.Msg)
	log.Info("query responder:\n mssage ->%s \n srcAddr-> %s \n destAddr-> %s", query.Msg.String(), query.mSrcAddr.String(), query.mDestAddr.String())
	{
		queryReplyFilter := NewQueryReplyFilter(query)
		responseFilter := NewQueryResponderRecordFilter().
			SetReplyFilter(queryReplyFilter)
		if !query.SendUnicast() {
			responseFilter.SetIncludeOnlyMulticastBeforeMS(time.Now().Add(-1 * time.Second))
		}
		for _, resp := range rs.mResponders {
			for _, info := range resp.ResponderInfos {
				if !responseFilter.Accept(info) {
					continue
				}
				record := info.Responder.GetRecord()
				if query.configuration != nil {
					query.configuration.Adjust(record)
				}
				msg.Answer = append(msg.Answer, info.Responder.GetRecord())
			}
		}

	}
	// send all 'Additional' replies
	{
		//rs.mSendState.SetResourceType(ResourceType_Additional)

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
				if query.configuration != nil {
					query.configuration.Adjust(record)
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

func (rs *ResponseSender) AddQueryResponder(queryResponder *QueryResponder) {
	rs.mResponders = append(rs.mResponders, queryResponder)
}
