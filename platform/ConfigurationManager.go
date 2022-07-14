package platform

import (
	"fmt"
	"github.com/galenliu/chip/config"
	"net"
	"sync"
)

var instance *ConfigurationManager
var _once sync.Once

func ConfigurationMgr() *ConfigurationManager {
	_once.Do(func() {
		instance = &ConfigurationManager{}
	})
	return instance
}

type ConfigurationManager struct {
	mVendorId                                  uint16
	mVendorName                                string
	mProductName                               string
	mProductId                                 uint16
	mDeviceType                                uint32
	mDeviceName                                string
	mTcpSupported                              bool
	mDevicePairingHint                         string
	mDevicePairingSecondaryHint                uint16
	mDeviceSecondaryPairingHint                string
	mDeviceConfigPairingInitialInstruction     string
	mDeviceConfigPairingSecondaryInstruction   string
	deviceConfigEnableCommissionableDeviceType bool
	mDiscriminator                             uint16 //_L<dddd>, where <dddd> provides the full 12-bit discriminator, encoded as a variable-length decimal number in ASCII text, omitting any leading zeroes.

}

func NewConfigurationManager(options *config.DeviceOptions) *ConfigurationManager {
	mgr := &ConfigurationManager{}
	if options.Payload.VendorID != 0 {
		mgr.StoreVendorId(options.Payload.VendorID)
	}
	if options.Payload.ProductID != 0 {
		mgr.StoreProductId(options.Payload.ProductID)
	}
	return mgr
}

func (c ConfigurationManager) GetVendorId() (uint16, error) {
	return c.mVendorId, nil
}

func (c ConfigurationManager) GetSetupDiscriminator() (uint16, error) {
	return c.mDiscriminator, nil
}

func (c ConfigurationManager) GetVendorName() (string, error) {
	return c.mVendorName, nil
}

func (c ConfigurationManager) GetProductId() (uint16, error) {
	return c.mProductId, nil
}

func (c ConfigurationManager) GetProductName() string {
	return c.mProductName
}

func (c ConfigurationManager) GetPrimaryMACAddress() string {
	return c.GetPrimaryWiFiMACAddress()
}

func (c ConfigurationManager) GetPrimaryWiFiMACAddress() string {
	ifs, _ := net.Interfaces()
	var value = ""
	for _, i := range ifs {
		if len(i.HardwareAddr) >= 6 {
			value = fmt.Sprintf("%X%X%X%X%X%X", i.HardwareAddr[0], i.HardwareAddr[1], i.HardwareAddr[2], i.HardwareAddr[3], i.HardwareAddr[4], i.HardwareAddr[5])
		}
	}
	return value
}

func (c ConfigurationManager) IsCommissionableDeviceTypeEnabled() bool {
	return c.deviceConfigEnableCommissionableDeviceType
}

func (c ConfigurationManager) GetDeviceTypeId() (uint32, error) {
	return c.mDeviceType, nil
}

func (c ConfigurationManager) SetDeviceTypeId(t uint32) {
	c.mDeviceType = t
}

func (c ConfigurationManager) IsCommissionableDeviceNameEnabled() bool {
	return true
}

func (c ConfigurationManager) GetCommissionableDeviceName() (string, error) {
	return c.mDeviceName, nil
}

func (c ConfigurationManager) GetInitialPairingHint() string {
	return c.mDevicePairingHint
}

func (c ConfigurationManager) GetInitialPairingInstruction() string {
	return c.mDeviceConfigPairingInitialInstruction
}

func (c ConfigurationManager) GetSecondaryPairingHint() string {
	return c.mDeviceSecondaryPairingHint
}

func (c ConfigurationManager) GetSecondaryPairingInstruction() string {
	return c.mDeviceConfigPairingSecondaryInstruction
}

func (c ConfigurationManager) StoreVendorId(id uint16) {
	c.mVendorId = id
}

func (c ConfigurationManager) StoreProductId(id uint16) {
	c.mProductId = id
}

func (c *ConfigurationManager) ReadConfigValueStr(num string) (string, error) {
	return "", nil
}

func (c ConfigurationManager) ReadConfigValue(version string) (uint32, error) {
	return 0, nil
}
