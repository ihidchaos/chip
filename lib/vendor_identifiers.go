package lib

type VendorId = uint16

const (
	CommonVendorId      VendorId = 0x0000
	UnspecifiedVendorId VendorId = 0x0000
	AppleVendorId       VendorId = 0x1349
	GoogleVendorId      VendorId = 0x6006
	TestVendor1         VendorId = 0xFFF1
	TestVendor2         VendorId = 0xFFF2
	TestVendor3         VendorId = 0xFFF3
	TestVendor4         VendorId = 0xFFF4
	NotSpecified        VendorId = 0xFFFF
)
