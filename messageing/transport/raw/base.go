package raw

import (
	"net/netip"
)

type TransportDelegate interface {
	HandleMessageReceived(srcAddr netip.AddrPort, buf *PacketBuffer)
}

type TransportBase interface {
	SetDelegate(m TransportDelegate)
	HandleMessageReceived(srcAddr netip.AddrPort, buf *PacketBuffer)
}
