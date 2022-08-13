package transport

import (
	"github.com/galenliu/chip/lib"
	"github.com/galenliu/chip/messageing/transport/raw"
	"net/netip"
)

type MessageCounterManagerInterface interface {
	StartSync(handle SessionHandle, session SecureSessionBase) error
	QueueReceivedMessageAndStartSync(
		header *raw.PacketHeader,
		handle SessionHandle,
		state uint8,
		peerAdders netip.AddrPort,
		buf *lib.PacketBuffer,
	) error
}

type GlobalUnencryptedMessageCounterImpl struct {
}

func NewGlobalUnencryptedMessageCounterImpl() *GlobalUnencryptedMessageCounterImpl {
	return &GlobalUnencryptedMessageCounterImpl{}
}

func (g *GlobalUnencryptedMessageCounterImpl) Init() {

}

func (g *GlobalUnencryptedMessageCounterImpl) StartSync(handle SessionHandle, session SecureSessionBase) error {
	//TODO implement me
	panic("implement me")
}

func (g *GlobalUnencryptedMessageCounterImpl) QueueReceivedMessageAndStartSync(header *raw.PacketHeader, handle SessionHandle, state uint8, peerAdders netip.AddrPort, buf *lib.PacketBuffer) error {
	//TODO implement me
	panic("implement me")
}
