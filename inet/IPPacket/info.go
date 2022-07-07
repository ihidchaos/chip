package IPPacket

import (
	"github.com/galenliu/chip/inet/IP"
	"github.com/galenliu/chip/inet/Interface"
)

type Info struct {
	SrcAddress  IP.Address
	DestAddress IP.Address
	InterfaceId Interface.Id
	SrcPort     uint16
	DestPort    uint16
}
