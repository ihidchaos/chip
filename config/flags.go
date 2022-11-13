package config

import (
	"github.com/galenliu/chip/platform/path"
)

var (
	EnableSessionResumption                   = false
	DeviceVendorId                     uint16 = 0xFFF1
	DeviceProductName                         = "TEST_PRODUCT"
	DeviceProductId                    uint16 = 0x8001
	DeviceVendorName                          = "TEST_VENDOR"
	TestSerialNumber                          = "TEST_SN"
	DefaultDeviceHardwareVersion       uint16 = 0
	DefaultDeviceHardwareVersionString        = "TEST_VERSION"

	PairingInitialInstruction      = ""
	PairingSecondaryInstruction    = ""
	DeviceName                     = "Test Kitchen"
	EnableCommissionableDeviceName = 0

	UseTestSetupPinCode      uint32 = 20202021
	RotatingDeviceIdUniqueId        = []byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff}
)

var (
	ChipDefaultFactoryPath = path.DefaultFactoryPath
	ChipDefaultConfigPath  = path.DefaultConfigPath
	ChipDefaultDataPath    = path.DefaultDataPath
)
