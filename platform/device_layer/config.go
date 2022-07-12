package DeviceLayer

const (
	CHIP_DEVICE_CONFIG_DEVICE_VENDOR_ID                       uint16 = 0xFFF1
	CHIP_DEVICE_CONFIG_DEVICE_PRODUCT_NAME                           = "TEST_PRODUCT"
	CHIP_DEVICE_CONFIG_DEVICE_PRODUCT_ID                             = 0x8001
	CHIP_DEVICE_CONFIG_DEVICE_VENDOR_NAME                            = "TEST_VENDOR"
	CHIP_DEVICE_CONFIG_TEST_SERIAL_NUMBER                            = "TEST_SN"
	CHIP_DEVICE_CONFIG_DEFAULT_DEVICE_HARDWARE_VERSION               = 0
	CHIP_DEVICE_CONFIG_DEFAULT_DEVICE_HARDWARE_VERSION_STRING        = "TEST_VERSION"
)

var CHIP_DEVICE_CONFIG_ROTATING_DEVICE_ID_UNIQUE_ID = []byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff}

const (
	kConfigKey_SerialNum         = "serial-num"
	kConfigKey_ManufacturingDate = "mfg-date"
	kConfigKey_HardwareVersion   = "hardware-ver"
)
