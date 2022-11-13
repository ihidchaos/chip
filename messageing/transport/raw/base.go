package raw

import (
	"github.com/galenliu/chip/platform/system"
	"net/netip"
)

type Delegate interface {
	HandleMessageReceived(srcAddr netip.AddrPort, msg *system.PacketBufferHandle)
}

type TransportBase interface {
	SetDelegate(m Delegate)
	SendMessage(peerAddr netip.AddrPort, msg *system.PacketBufferHandle) error
	Disconnect(peerAddr netip.AddrPort)
	Close()
}
