package protocols

import "github.com/galenliu/chip/lib"

type Id struct {
	mVendorId   uint16
	mProtocolId uint16
}

var NotSpecifiedId = Id{
	mVendorId:   lib.VendorIdNotSpecified,
	mProtocolId: 0xFFFF,
}
