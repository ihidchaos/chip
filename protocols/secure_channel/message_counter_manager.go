package secure_channel

import (
	"github.com/galenliu/chip/messageing"
	"github.com/galenliu/chip/messageing/transport"
	"github.com/galenliu/chip/messageing/transport/raw"
	"net/netip"
)

type MessageCounterManager interface {
	messageing.UnsolicitedMessageHandler
	messageing.ExchangeDelegate
	transport.MessageCounterManager
}

type MessageCounterManagerImpl struct {
}

func (m *MessageCounterManagerImpl) OnMessageReceived(context *messageing.ExchangeContext, header *raw.PayloadHeader, data *raw.PacketBuffer) error {
	if header.HasMessageType(uint8(MsgCounterSyncReq)) {
		return m.handleMsgCounterSyncReq(context, header)
	}
	if header.HasMessageType(uint8(MsgCounterSyncRsp)) {
		return m.handleMsgCounterSyncResp(context, header)
	}
	return nil
}

func (m *MessageCounterManagerImpl) OnResponseTimeout(ec *messageing.ExchangeContext) {
	if ec.HasSessionHandle() {
		ec.GetSessionHandle()
	}
}

func (m *MessageCounterManagerImpl) OnExchangeClosing(ec *messageing.ExchangeContext) {
	//TODO implement me
	panic("implement me")
}

func (m *MessageCounterManagerImpl) OnUnsolicitedMessageReceived(header *raw.PayloadHeader, delegate messageing.ExchangeDelegate) error {
	//TODO implement me
	panic("implement me")
}

func (m *MessageCounterManagerImpl) OnExchangeCreationFailed(delegate messageing.ExchangeDelegate) {
	//TODO implement me
	panic("implement me")
}

func (m *MessageCounterManagerImpl) StartSync(handle transport.SessionHandleBase, session transport.SecureSessionBase) error {
	//TODO implement me
	panic("implement me")
}

func (m *MessageCounterManagerImpl) QueueReceivedMessageAndStartSync(header *raw.PacketHeader, handle transport.SessionHandleBase, state uint8, peerAdders netip.AddrPort, buf *raw.PacketBuffer) error {
	//TODO implement me
	panic("implement me")
}

func (m *MessageCounterManagerImpl) Init(mgr messageing.ExchangeManager) error {
	return nil
}

func (m *MessageCounterManagerImpl) handleMsgCounterSyncReq(context *messageing.ExchangeContext, header *raw.PayloadHeader) error {
	//TODO implement me
	panic("implement me")
}

func (m *MessageCounterManagerImpl) handleMsgCounterSyncResp(context *messageing.ExchangeContext, header *raw.PayloadHeader) error {
	//TODO implement me
	panic("implement me")
}

func NewMessageCounterManager() *MessageCounterManagerImpl {
	return &MessageCounterManagerImpl{}
}
