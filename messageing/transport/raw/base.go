package raw

import (
	"github.com/galenliu/chip/platform/system"
	"net/netip"
)

type TransportDelegate interface {
	HandleMessageReceived(srcAddr netip.AddrPort, buf *system.PacketBufferHandle)
}

type TransportBase interface {
	SetDelegate(m TransportDelegate)
	HandleMessageReceived(srcAddr netip.AddrPort, buf *PacketBuffer)
}
