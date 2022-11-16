package raw

import "net"

type Type uint8

const (
	Undefined Type = iota
	Udp
	Ble
	Tcp
)

type PeerAddress struct {
	mTransportType Type
	mInterface     net.Interface
	mPort          uint16
}

func (a PeerAddress) TransportType() Type {
	return a.mTransportType
}
