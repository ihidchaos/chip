package raw

import "net/netip"

type TransportBaseDelegate interface {
	HandleMessageReceived(addrPort netip.AddrPort, msg []byte)
}

type TransportBase interface {
	TransportBaseDelegate
	GetBoundPort() uint16
	SetDelegate(m TransportBaseDelegate)
}
