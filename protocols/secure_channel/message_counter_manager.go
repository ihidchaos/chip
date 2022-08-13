package secure_channel

import (
	"github.com/galenliu/chip/lib"
	"github.com/galenliu/chip/messageing"
	"github.com/galenliu/chip/messageing/transport"
	"github.com/galenliu/chip/messageing/transport/raw"
	"net/netip"
)

type MessageCounterManager struct {
}

func (m MessageCounterManager) StartSync(handle transport.SessionHandle, session transport.SecureSessionBase) error {
	//TODO implement me
	panic("implement me")
}

func (m MessageCounterManager) QueueReceivedMessageAndStartSync(header *raw.PacketHeader, handle transport.SessionHandle, state uint8, peerAdders netip.AddrPort, buf *lib.PacketBuffer) error {
	//TODO implement me
	panic("implement me")
}

func (m MessageCounterManager) Init(mgr messageing.ExchangeManager) error {
	return nil
}

func NewMessageCounterManager() *MessageCounterManager {
	return &MessageCounterManager{}
}
