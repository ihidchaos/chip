package config

import (
	"github.com/galenliu/chip/system"
)

var (
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

	ChipDeviceConfigRotatingDeviceIdUniqueId = []byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff}
)

const (
	kconfigkeySerialnum         = "serial-num"
	kconfigkeyManufacturingdate = "mfg-date"
	kconfigkeyHardwareversion   = "hardware-ver"
)

var (
	ChipDefaultFactoryPath = system.GetFatConFile()
	ChipDefaultConfigPath  = system.GetSysConFile()
	ChipDefaultDataPath    = system.GetLocalConFile()
)
