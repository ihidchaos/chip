package lib

import "fmt"

type VendorId uint16

const (
	VendorIdCommon         VendorId = 0x0000
	VendorIdMatterStandard VendorId = 0x0000
	VendorIdUnspecified    VendorId = 0x0000
	VendorIdApple          VendorId = 0x1349
	VendorIdGoogle         VendorId = 0x6006
	VendorTest1            VendorId = 0xFFF1
	VendorTest2            VendorId = 0xFFF2
	VendorTest3            VendorId = 0xFFF3
	VendorTest4            VendorId = 0xFFF4
	VendorIdNotSpecified   VendorId = 0xFFFF
)

func (id VendorId) String() string {
	var value = uint16(id)
	return fmt.Sprintf("%04X", value)
}
