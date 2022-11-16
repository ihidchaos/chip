package device

import (
	"fmt"
	"github.com/galenliu/chip/config"
	"log"
	"sync/atomic"
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

type InstanceInfoImpl struct {
	mConfigManager config.ConfigurationManager
}

var defaultInstanceInfoProvider atomic.Value

func init() {
	defaultInstanceInfoProvider.Store(NewInstanceInfo())
}

func DefaultInstanceInfo() *InstanceInfoImpl {
	return defaultInstanceInfoProvider.Load().(*InstanceInfoImpl)
}

func NewInstanceInfo() *InstanceInfoImpl {
	return &InstanceInfoImpl{}
}

func (info *InstanceInfoImpl) Init(configMgr config.ConfigurationManager) error {
	info.mConfigManager = configMgr
	return nil
}

func (info *InstanceInfoImpl) GetVendorId() (uint16, error) {
	return config.DeviceVendorId, nil
}

func (info *InstanceInfoImpl) GetProductId() (uint16, error) {

	return config.DeviceProductId, nil

}

func (info *InstanceInfoImpl) ProductName() (string, error) {
	return config.DeviceProductName, nil
}

func (info *InstanceInfoImpl) VendorName() (string, error) {
	return config.DeviceVendorName, nil
}

func (info *InstanceInfoImpl) SerialNumber() (string, error) {
	sn, err := info.mConfigManager.ReadConfigValueStr(config.KConfigKey_SerialNum)
	if sn == "" || err != nil {
		return config.TestSerialNumber, nil
	}
	return sn, nil
}

func (info *InstanceInfoImpl) GetManufacturingDate() (time.Time, error) {
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

func (info *InstanceInfoImpl) GetHardwareVersion() (uint16, error) {
	version, err := info.mConfigManager.ReadConfigValueUint16(config.KConfigKey_HardwareVersion)
	if err != nil {
		return config.DefaultDeviceHardwareVersion, nil
	}
	return version, nil
}

func (info *InstanceInfoImpl) HardwareVersionString() (string, error) {
	return config.DefaultDeviceHardwareVersionString, nil
}

func (info *InstanceInfoImpl) GetRotatingDeviceIdUniqueId() ([]byte, error) {
	return config.RotatingDeviceIdUniqueId, nil
}
