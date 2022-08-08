package protocols

import "github.com/galenliu/chip/lib"

type Id struct {
	VendorId   lib.VendorId
	ProtocolId uint16
}

var StandardProtocolId = &Id{
	VendorId:   lib.VidUnspecified,
	ProtocolId: 0xFFFF,
}

func (id *Id) Equal(other *Id) bool {
	return id.VendorId == other.VendorId && id.ProtocolId == other.ProtocolId
}
