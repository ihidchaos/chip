package device

import (
	"fmt"
	"github.com/galenliu/chip/config"
	"log"
	"sync"
	"time"
)

type InstanceInfoProvider interface {
	VendorName() (string, error)
	GetVendorId() (uint16, error)

	ProductName() (string, error)
	GetProductId() (uint16, error)

	SerialNumber() (string, error)

	GetManufacturingDate() (time.Time, error)

	GetHardwareVersion() (uint16, error)
	HardwareVersionString() (string, error)

	GetRotatingDeviceIdUniqueId() ([]byte, error)
}

type InstanceInfo struct {
	mConfigManager config.ConfigurationManager
}

func NewDeviceInstanceInfoImpl() *InstanceInfo {
	return GetDeviceInstanceInfoProvider()
}

var _deviceInstanceInfo *InstanceInfo
var _deviceInstanceInfoOnce sync.Once

func GetDeviceInstanceInfoProvider() *InstanceInfo {
	_deviceInstanceInfoOnce.Do(func() {
		if _deviceInstanceInfo == nil {
			_deviceInstanceInfo = &InstanceInfo{}
		}
	})
	return _deviceInstanceInfo
}

func (info *InstanceInfo) Init(configMgr config.ConfigurationManager) (*InstanceInfo, error) {
	info.mConfigManager = configMgr
	return info, nil
}

func NewInstanceInfo() *InstanceInfo {
	return GetDeviceInstanceInfoProvider()
}

func (info *InstanceInfo) GetVendorId() (uint16, error) {
	return config.DeviceVendorId, nil
}

func (info *InstanceInfo) GetProductId() (uint16, error) {

	return config.DeviceProductId, nil

}

func (info *InstanceInfo) ProductName() (string, error) {
	return config.DeviceProductName, nil
}

func (info *InstanceInfo) VendorName() (string, error) {
	return config.DeviceVendorName, nil
}

func (info *InstanceInfo) SerialNumber() (string, error) {
	sn, err := info.mConfigManager.ReadConfigValueStr(config.KConfigKey_SerialNum)
	if sn == "" || err != nil {
		return config.TestSerialNumber, nil
	}
	return sn, nil
}

func (info *InstanceInfo) GetManufacturingDate() (time.Time, error) {
	data, err := info.mConfigManager.ReadConfigValueStr(config.KConfigKey_ManufacturingDate)
	if err != nil {
		log.Panicf("invalid manufacturing date: %s", err.Error())
	}
	t, err := time.Parse("2006-Jan-02", data)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid manufacturing date: %s", err.Error())
	}
	return t, nil
}

func (info *InstanceInfo) GetHardwareVersion() (uint16, error) {
	version, err := info.mConfigManager.ReadConfigValueUint16(config.KConfigKey_HardwareVersion)
	if err != nil {
		return config.DefaultDeviceHardwareVersion, nil
	}
	return version, nil
}

func (info *InstanceInfo) HardwareVersionString() (string, error) {
	return config.DefaultDeviceHardwareVersionString, nil
}

func (info *InstanceInfo) GetRotatingDeviceIdUniqueId() ([]byte, error) {
	return config.RotatingDeviceIdUniqueId, nil
}
