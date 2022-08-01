package protocols

import "github.com/galenliu/chip/lib"

type Id struct {
	mVendorId   uint16
	mProtocolId uint16
}

func NotSpecifiedId() Id {
	return Id{
		lib.VendorIdNotSpecified,
		0xFFFF,
	}
}
