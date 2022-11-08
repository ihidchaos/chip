package secure_channel

import (
	"github.com/galenliu/chip/messageing"
	"github.com/galenliu/chip/messageing/transport"
	"github.com/galenliu/chip/messageing/transport/raw"
	"net/netip"
)

type MessageCounterManagerBase interface {
	messageing.UnsolicitedMessageHandler
	messageing.ExchangeDelegate
	transport.MessageCounterManager
}

type MessageCounterManager struct {
}

func (m *MessageCounterManager) OnMessageReceived(context *messageing.ExchangeContext, header *raw.PayloadHeader, data *raw.PacketBuffer) error {
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
		ec.GetSessionHandle()
	}
}

func (m *MessageCounterManager) OnExchangeClosing(ec *messageing.ExchangeContext) {
	//TODO implement me
	panic("implement me")
}

func (m *MessageCounterManager) OnUnsolicitedMessageReceived(header *raw.PayloadHeader, delegate messageing.ExchangeDelegate) error {
	//TODO implement me
	panic("implement me")
}

func (m *MessageCounterManager) OnExchangeCreationFailed(delegate messageing.ExchangeDelegate) {
	//TODO implement me
	panic("implement me")
}

func (m *MessageCounterManager) StartSync(handle transport.SessionHandle, session *transport.SecureSessionBase) error {
	//TODO implement me
	panic("implement me")
}

func (m *MessageCounterManager) QueueReceivedMessageAndStartSync(header *raw.PacketHeader, handle *transport.SessionHandle, state uint8, peerAdders netip.AddrPort, buf *raw.PacketBuffer) error {
	//TODO implement me
	panic("implement me")
}

func (m *MessageCounterManager) Init(mgr messageing.ExchangeManagerBase) error {
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
