package DeviceLayer

import (
	"fmt"
	"github.com/galenliu/chip/platform"
	"time"
)

type DeviceInfo struct {
	VendorId        uint16
	ProductId       uint16
	ProductName     string
	VendorName      string
	HardwareVersion string
}

type DeviceInstanceInfo struct {
	configManager *platform.ConfigurationManager
	info          *DeviceInfo
}

func NewDeviceInstanceInfo(configManager *platform.ConfigurationManager, info *DeviceInfo) *DeviceInstanceInfo {
	return &DeviceInstanceInfo{
		configManager: configManager,
		info:          info,
	}
}

func (d DeviceInstanceInfo) GetVendorId() (uint16, error) {
	if d.info.VendorId == 0 {
		return CHIP_DEVICE_CONFIG_DEVICE_VENDOR_ID, nil
	}
	return d.info.VendorId, nil
}

func (d DeviceInstanceInfo) GetProductId() (uint16, error) {
	if d.info.ProductId == 0 {
		return CHIP_DEVICE_CONFIG_DEVICE_PRODUCT_ID, nil
	}
	return d.info.ProductId, nil
}

func (d DeviceInstanceInfo) GetProductName() (string, error) {
	if d.info.ProductName == "" {
		return CHIP_DEVICE_CONFIG_DEVICE_PRODUCT_NAME, nil
	}
	return d.info.ProductName, nil
}

func (d DeviceInstanceInfo) GetVendorName() (string, error) {
	if d.info.VendorName == "" {
		return CHIP_DEVICE_CONFIG_DEVICE_VENDOR_NAME, nil
	}
	return d.info.VendorName, nil
}

func (d DeviceInstanceInfo) GetSerialNumber() (string, error) {
	sn, err := d.configManager.ReadConfigValueStr(kConfigKey_SerialNum)
	if sn == "" || err != nil {
		return CHIP_DEVICE_CONFIG_TEST_SERIAL_NUMBER, nil
	}
	return sn, nil
}

func (d DeviceInstanceInfo) GetManufacturingDate() (time.Time, error) {
	data, err := d.configManager.ReadConfigValueStr(kConfigKey_ManufacturingDate)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid manufacturing date: %s", err.Error())
	}
	t, err := time.Parse("2006-Jan-02", data)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid manufacturing date: %s", err.Error())
	}
	return t, nil
}

func (d DeviceInstanceInfo) GetHardwareVersion() (uint16, error) {
	data, err := d.configManager.ReadConfigValue(kConfigKey_HardwareVersion)
	if err != nil {
		return CHIP_DEVICE_CONFIG_DEFAULT_DEVICE_HARDWARE_VERSION, nil
	}
	return uint16(data), nil
}

func (d DeviceInstanceInfo) GetHardwareVersionString() (string, error) {
	if d.info.HardwareVersion != "" {
		return d.info.HardwareVersion, nil
	}
	return CHIP_DEVICE_CONFIG_DEFAULT_DEVICE_HARDWARE_VERSION_STRING, nil
}

func (d DeviceInstanceInfo) GetRotatingDeviceIdUniqueId() ([]byte, error) {
	return CHIP_DEVICE_CONFIG_ROTATING_DEVICE_ID_UNIQUE_ID, nil
}

type DeviceInstanceInfoProvider interface {
	GetVendorId() (uint16, error)
	GetProductId() (uint16, error)
	GetProductName() (string, error)
	GetVendorName() (string, error)
	GetSerialNumber() (string, error)
	GetManufacturingDate() (time.Time, error)
	GetHardwareVersion() (uint16, error)
	GetHardwareVersionString() (string, error)
	GetRotatingDeviceIdUniqueId() ([]byte, error)
}

type DeviceInfoProvider interface {
}

func GetDeviceInfoProvider() DeviceInfoProvider {
	return platform.ConfigurationMgr()
}

func GetDeviceInstanceInfoProvider() DeviceInstanceInfoProvider {
	return nil
}
