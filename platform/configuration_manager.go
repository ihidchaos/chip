package platform

import (
	"github.com/galenliu/chip/device"
	"net"
	"strings"
	"sync"
)

type Config struct {
	ChipDeviceConfigPairingSecondaryHint        int
	ChipDeviceConfigDeviceVendorName            string
	ChipDeviceConfigDeviceType                  int
	ChipDeviceConfigDeviceProductName           string
	ChipDeviceConfigPairingInitialHint          uint16
	ChipDeviceConfigPairingInitialInstruction   string
	ChipDeviceConfigPairingSecondaryInstruction string
}

var cmInstance *ConfigurationManager
var cmOnce sync.Once

func ConfigurationMgr() *ConfigurationManager {
	cmOnce.Do(func() {
		cmInstance = newConfigurationManager(Config{
			ChipDeviceConfigPairingSecondaryHint:        0,
			ChipDeviceConfigDeviceVendorName:            "",
			ChipDeviceConfigDeviceType:                  0,
			ChipDeviceConfigDeviceProductName:           "",
			ChipDeviceConfigPairingInitialHint:          0,
			ChipDeviceConfigPairingInitialInstruction:   "",
			ChipDeviceConfigPairingSecondaryInstruction: "",
		})
	})
	return cmInstance
}

type ConfigurationManager struct {
	mVendorId                                  uint16
	mVendorName                                string
	mProductName                               string
	mProductId                                 int16
	mDeviceType                                device.MatterDeviceType
	mDeviceName                                string
	mTcpSupported                              bool
	mDevicePairingHint                         uint16
	mDevicePairingSecondaryHint                uint16
	mDeviceSecondaryPairingHint                uint16
	mDeviceConfigPairingInitialInstruction     string
	mDeviceConfigPairingSecondaryInstruction   string
	deviceConfigEnableCommissionableDeviceType bool
}

func newConfigurationManager(conf Config) *ConfigurationManager {
	return &ConfigurationManager{
		mVendorId:                                0,
		mVendorName:                              "",
		mProductName:                             "",
		mProductId:                               0,
		mDeviceType:                              conf.ChipDeviceConfigDeviceType,
		mDeviceName:                              conf.ChipDeviceConfigDeviceProductName,
		mDevicePairingHint:                       conf.ChipDeviceConfigPairingInitialHint,
		mDeviceConfigPairingInitialInstruction:   conf.ChipDeviceConfigPairingInitialInstruction,
		mDeviceConfigPairingSecondaryInstruction: conf.ChipDeviceConfigPairingSecondaryInstruction,
		deviceConfigEnableCommissionableDeviceType: false,
	}
}

func (c ConfigurationManager) GetVendorId() (uint16, error) {
	return c.mVendorId, nil
}

func (c ConfigurationManager) GetVendorName() string {
	return c.mVendorName
}

func (c ConfigurationManager) GetProductId() (int16, error) {
	return c.mProductId, nil
}

func (c ConfigurationManager) GetProductName() string {
	return c.mProductName
}

func (c ConfigurationManager) GetPrimaryMACAddress() (mac net.HardwareAddr, err error) {
	return c.GetPrimaryWiFiMACAddress(), nil
}

func (c ConfigurationManager) GetPrimaryWiFiMACAddress() (mac net.HardwareAddr) {
	ifs, _ := net.Interfaces()
	for _, i := range ifs {
		if !strings.Contains(i.Name, "lo") && i.HardwareAddr != nil {
			mac = i.HardwareAddr
		}
	}
	return
}

func (c ConfigurationManager) IsCommissionableDeviceTypeEnabled() bool {
	return c.deviceConfigEnableCommissionableDeviceType
}

func (c ConfigurationManager) GetDeviceTypeId() (device.MatterDeviceType, error) {
	return c.mDeviceType, nil
}

func (c ConfigurationManager) SetDeviceTypeId(t device.MatterDeviceType) {
	c.mDeviceType = t
}

func (c ConfigurationManager) IsCommissionableDeviceNameEnabled() bool {
	return true
}

func (c ConfigurationManager) GetCommissionableDeviceName() (string, error) {
	return c.mDeviceName, nil
}

func (c ConfigurationManager) GetInitialPairingHint() uint16 {
	return c.mDevicePairingHint
}

func (c ConfigurationManager) GetInitialPairingInstruction() string {
	return c.mDeviceConfigPairingInitialInstruction
}

func (c ConfigurationManager) GetSecondaryPairingHint() uint16 {
	return c.mDeviceSecondaryPairingHint
}

func (c ConfigurationManager) GetSecondaryPairingInstruction() string {
	return c.mDeviceConfigPairingSecondaryInstruction
}
