package config

var (
	ChipDeviceConfigDeviceVendorId                     uint16 = 0xFFF1
	ChipDeviceConfigDeviceProductName                         = "TEST_PRODUCT"
	ChipDeviceConfigDeviceProductId                    uint16 = 0x8001
	ChipDeviceConfigDeviceVendorName                          = "TEST_VENDOR"
	ChipDeviceConfigTestSerialNumber                          = "TEST_SN"
	ChipDeviceConfigDefaultDeviceHardwareVersion       uint16 = 0
	ChipDeviceConfigDefaultDeviceHardwareVersionString        = "TEST_VERSION"

	ChipDeviceConfigRotatingDeviceIdUniqueId = []byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff}
)

const (
	kConfigKey_SerialNum         = "serial-num"
	kConfigKey_ManufacturingDate = "mfg-date"
	kConfigKey_HardwareVersion   = "hardware-ver"
)
