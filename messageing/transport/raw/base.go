package raw

import (
	"github.com/galenliu/chip/lib/buffer"
	"net/netip"
)

type TransportDelegate interface {
	HandleMessageReceived(srcAddr netip.AddrPort, buf *buffer.PacketBuffer)
}

type TransportBase interface {
	TransportDelegate
	GetBoundPort() uint16
	SetDelegate(m TransportDelegate)
}
