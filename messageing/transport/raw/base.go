package raw

import "net/netip"

type TransportDelegate interface {
	HandleMessageReceived(srcAddr netip.AddrPort, buf *PacketBuffer)
}

type TransportBase interface {
	TransportDelegate
	GetBoundPort() uint16
	SetDelegate(m TransportDelegate)
}
