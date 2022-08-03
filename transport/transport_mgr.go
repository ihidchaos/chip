package transport

import "net/netip"

type RawTransportDelegate interface {
	HandleMessageReceived(peerAddress netip.AddrPort, msg []byte)
}

type TransportMgrBase interface {
	RawTransportDelegate
	MulticastGroupJoinLeave(addr netip.Addr, join bool) error
	SendMessage(port netip.AddrPort, msg []byte) error
	Close()
	Disconnect(addr netip.Addr)
}
