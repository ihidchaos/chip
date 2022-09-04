package secure

import (
	"github.com/galenliu/chip/lib/buffer"
	"github.com/galenliu/chip/messageing"
	"github.com/galenliu/chip/messageing/transport"
	"github.com/galenliu/chip/messageing/transport/raw"
	"net/netip"
)

type MessageCounterManager struct {
}

func (m MessageCounterManager) StartSync(handle transport.SessionHandleBase, session transport.SecureSessionBase) error {
	//TODO implement me
	panic("implement me")
}

func (m MessageCounterManager) QueueReceivedMessageAndStartSync(header *raw.PacketHeader, handle transport.SessionHandleBase, state uint8, peerAdders netip.AddrPort, buf *buffer.PacketBuffer) error {
	//TODO implement me
	panic("implement me")
}

func (m MessageCounterManager) Init(mgr messageing.ExchangeManager) error {
	return nil
}

func NewMessageCounterManager() *MessageCounterManager {
	return &MessageCounterManager{}
}
