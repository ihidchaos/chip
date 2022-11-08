package config

import (
	"github.com/galenliu/chip/system/platform"
)

var (
	ChipConfigEnableSessionResumption                         = false
	ChipDeviceConfigDeviceVendorId                     uint16 = 0xFFF1
	ChipDeviceConfigDeviceProductName                         = "TEST_PRODUCT"
	ChipDeviceConfigDeviceProductId                    uint16 = 0x8001
	ChipDeviceConfigDeviceVendorName                          = "TEST_VENDOR"
	ChipDeviceConfigTestSerialNumber                          = "TEST_SN"
	ChipDeviceConfigDefaultDeviceHardwareVersion       uint16 = 0
	ChipDeviceConfigDefaultDeviceHardwareVersionString        = "TEST_VERSION"

	ChipDeviceConfigPairingInitialInstruction      = ""
	ChipDeviceConfigPairingSecondaryInstruction    = ""
	ChipDeviceConfigDeviceName                     = "Test Kitchen"
	ChipDeviceConfigEnableCommissionableDeviceName = 0

	ChipDeviceConfigUseTestSetupPinCode uint32 = 20202021

	ChipDeviceConfigRotatingDeviceIdUniqueId = []byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff}
)

var (
	ChipDefaultFactoryPath = platform.DefaultFactoryPath
	ChipDefaultConfigPath  = platform.DefaultConfigPath
	ChipDefaultDataPath    = platform.DefaultDataPath
)
