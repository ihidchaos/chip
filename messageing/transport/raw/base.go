package raw

import (
	"github.com/galenliu/chip/lib"
	"net/netip"
)

type TransportDelegate interface {
	HandleMessageReceived(srcAddr netip.AddrPort, buf *lib.PacketBuffer)
}

type TransportBase interface {
	TransportDelegate
	GetBoundPort() uint16
	SetDelegate(m TransportDelegate)
}
