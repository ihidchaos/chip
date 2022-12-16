package raw

import (
	"net"
	"net/netip"
)

type TransportType uint8

const (
	Undefined TransportType = iota
	Udp
	Ble
	Tcp
)

type PeerAddress struct {
	TransportType TransportType
	Interface     net.Interface
	AddrPort      netip.AddrPort
}
