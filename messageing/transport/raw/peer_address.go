package raw

import (
	"net"
	"net/netip"
)

type Type uint8

const (
	Undefined Type = iota
	Udp
	Ble
	Tcp
)

type PeerAddress struct {
	TransportType Type
	Interface     net.Interface
	AddrPort      netip.AddrPort
}
