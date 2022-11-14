package lib

import "fmt"

type VendorId uint16

const (
	VidCommon       VendorId = 0x0000
	VidUnspecified  VendorId = 0x0000
	VidApple        VendorId = 0x1349
	VidGoogle       VendorId = 0x6006
	VidTest1        VendorId = 0xFFF1
	VidTest2        VendorId = 0xFFF2
	VidTest3        VendorId = 0xFFF3
	VidTest4        VendorId = 0xFFF4
	VidNotSpecified VendorId = 0xFFFF
)

func (id VendorId) String() string {
	var value = uint16(id)
	return fmt.Sprintf("%04X", value)
}

func (id VendorId) IsTest() bool {
	switch id {
	case VidTest1, VidTest2, VidTest3, VidTest4:
		return true
	default:
		return false
	}
}

func (id VendorId) IsValidOperationally() bool {
	return id != VidCommon && id <= VidTest4
}
