package protocols

import "github.com/galenliu/chip/lib"

type Id struct {
	mVendorId   lib.VendorId
	mProtocolId uint16
}

var StandardProtocolId = &Id{
	mVendorId:   lib.UnspecifiedVendorId,
	mProtocolId: 0xFFFF,
}
