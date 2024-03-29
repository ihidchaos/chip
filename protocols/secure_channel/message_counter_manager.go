package secure_channel

import (
	"github.com/galenliu/chip/messageing"
	"github.com/galenliu/chip/messageing/transport"
	"github.com/galenliu/chip/messageing/transport/raw"
	"github.com/galenliu/chip/messageing/transport/session"
	"github.com/galenliu/chip/platform/system"
	"net/netip"
)

type MessageCounterManagerBase interface {
	messageing.UnsolicitedMessageHandler
	messageing.ExchangeDelegate
	transport.MessageCounterManagerBase
}

type MessageCounterManager struct {
}

func (m *MessageCounterManager) OnMessageReceived(context *messageing.ExchangeContext, header *raw.PayloadHeader, data *system.PacketBufferHandle) error {
	if header.HasMessageType(uint8(MsgCounterSyncReq)) {
		return m.handleMsgCounterSyncReq(context, header)
	}
	if header.HasMessageType(uint8(MsgCounterSyncRsp)) {
		return m.handleMsgCounterSyncResp(context, header)
	}
	return nil
}

func (m *MessageCounterManager) OnResponseTimeout(ec *messageing.ExchangeContext) {
	if ec.HasSessionHandle() {
		ec.SessionHandle()
	}
}

func (m *MessageCounterManager) OnExchangeClosing(ec *messageing.ExchangeContext) {
	//TODO implement me
	panic("implement me")
}

func (m *MessageCounterManager) OnUnsolicitedMessageReceived(header *raw.PayloadHeader) (messageing.ExchangeDelegate, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MessageCounterManager) OnExchangeCreationFailed(delegate messageing.ExchangeDelegate) {
	//TODO implement me
	panic("implement me")
}

func (m *MessageCounterManager) StartSync(handle *transport.SessionHandle, session *session.Secure) error {
	//TODO implement me
	panic("implement me")
}

func (m *MessageCounterManager) QueueReceivedMessageAndStartSync(header *raw.PacketHeader,
	handle *transport.SessionHandle,
	state *session.Secure,
	peerAdders netip.AddrPort,
	buf *system.PacketBufferHandle) error {
	//TODO implement me
	panic("implement me")
}

func (m *MessageCounterManager) Init(mgr *messageing.ExchangeManager) error {
	return nil
}

func (m *MessageCounterManager) handleMsgCounterSyncReq(context *messageing.ExchangeContext, header *raw.PayloadHeader) error {
	//TODO implement me
	panic("implement me")
}

func (m *MessageCounterManager) handleMsgCounterSyncResp(context *messageing.ExchangeContext, header *raw.PayloadHeader) error {
	//TODO implement me
	panic("implement me")
}

func NewMessageCounterManager() *MessageCounterManager {
	return &MessageCounterManager{}
}
