package platform

import (
	"github.com/galenliu/chip/ble"
	"github.com/galenliu/chip/clusters"
	"github.com/galenliu/chip/config"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

// ConfigurationManager Defines the public interface for the Device Layer ConfigurationManager object.
type ConfigurationManager interface {
	ConfigDelegate
	StoreVendorId(vendorId uint16) error
	StoreProductId(productId uint16) error

	GetPrimaryMACAddress() (string, error)
	GetPrimaryWiFiMACAddress() ([]byte, error)
	GetPrimary802154MACAddress() ([]byte, error)

	GetSoftwareVersionString()
	GetSoftwareVersion() (uint32, error)
	GetFirmwareBuildChipEpochTime() (time.Duration, error)
	SetFirmwareBuildChipEpochTime() (time.Duration, error)

	GetLifetimeCounter() (uint16, error)

	//IncrementLifetimeCounter()error

	//SetRotatingDeviceIdUniqueId(const ByteSpan & uniqueIdSpan) = 0;

	GetRegulatoryLocation() (location uint8, err error)
	GetCountryCode() (string, error)
	StoreSerialNumber(serialNum string) error
	StoreManufacturingDate(mfgDate string) error
	StoreSoftwareVersion(softwareVer uint32) error
	StoreHardwareVersion(hardwareVer uint16) error
	StoreRegulatoryLocation(location uint8) error
	StoreCountryCode(code string) error
	GetRebootCount() (uint32, error)
	StoreRebootCount(rebootCount uint32) error
	GetTotalOperationalHours(totalOperationalHours uint32) error
	StoreTotalOperationalHours(totalOperationalHours uint32) error
	GetBootReason(bootReason uint32) error
	StoreBootReason(bootReason uint32) error
	GetPartNumber() (string, error)
	GetProductURL() (string, error)
	GetProductLabel() (string, error)
	GetUniqueId() (string, error)
	StoreUniqueId(uniqueId string) error
	GenerateUniqueId() error
	GetFailSafeArmed() bool
	SetFailSafeArmed(val bool) error
	//
	GetBLEDeviceIdentificationInfo() (ble.ChipBLEDeviceIdentificationInfo, error)

	IsFullyProvisioned() bool
	InitiateFactoryReset()

	RunUnitTests() error

	LogDeviceConfig()
	IsCommissionableDeviceTypeEnabled() bool
	GetDeviceTypeId() (deviceType uint32, err error)
	IsCommissionableDeviceNameEnabled() bool
	GetCommissionableDeviceName() (string, error)
	GetInitialPairingHint() (uint16, error)
	GetInitialPairingInstruction() (string, error)
	GetSecondaryPairingHint() (uint16, error)
	GetSecondaryPairingInstruction() (string, error)

	GetLocationCapability() (uint8, error)
}

type ConfigurationManagerImpl struct {
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
	ConfigProvider
}

var instance *ConfigurationManagerImpl
var _once sync.Once

func ConfigurationMgr() *ConfigurationManagerImpl {
	_once.Do(func() {
		instance = &ConfigurationManagerImpl{}
		instance.ConfigProvider = NewConfigProviderImpl()
	})
	return instance
}

func (c *ConfigurationManagerImpl) GetRegulatoryLocation() (location uint8, err error) {
	//TODO implement me
	panic("implement me")
}

func (c *ConfigurationManagerImpl) GetCountryCode() (string, error) {
	//TODO implement me
	panic("implement me")
}

func (c *ConfigurationManagerImpl) StoreSerialNumber(serialNum string) error {
	//TODO implement me
	panic("implement me")
}

func (c *ConfigurationManagerImpl) StoreManufacturingDate(mfgDate string) error {
	//TODO implement me
	panic("implement me")
}

func (c *ConfigurationManagerImpl) StoreSoftwareVersion(softwareVer uint32) error {
	//TODO implement me
	panic("implement me")
}

func (c *ConfigurationManagerImpl) StoreHardwareVersion(hardwareVer uint16) error {
	//TODO implement me
	panic("implement me")
}

func (c *ConfigurationManagerImpl) StoreRegulatoryLocation(location uint8) error {
	//TODO implement me
	panic("implement me")
}

func (c *ConfigurationManagerImpl) StoreCountryCode(code string) error {
	//TODO implement me
	panic("implement me")
}

func (c *ConfigurationManagerImpl) GetTotalOperationalHours(totalOperationalHours uint32) error {
	//TODO implement me
	panic("implement me")
}

func (c *ConfigurationManagerImpl) GetBootReason(bootReason uint32) error {
	//TODO implement me
	panic("implement me")
}

func (c *ConfigurationManagerImpl) GetPartNumber() (string, error) {
	//TODO implement me
	panic("implement me")
}

func (c *ConfigurationManagerImpl) GetProductURL() (string, error) {
	//TODO implement me
	panic("implement me")
}

func (c *ConfigurationManagerImpl) GetProductLabel() (string, error) {
	//TODO implement me
	panic("implement me")
}

func (c *ConfigurationManagerImpl) GetUniqueId() (string, error) {
	//TODO implement me
	panic("implement me")
}

func (c *ConfigurationManagerImpl) StoreUniqueId(uniqueId string) error {
	//TODO implement me
	panic("implement me")
}

func (c *ConfigurationManagerImpl) GenerateUniqueId() error {
	//TODO implement me
	panic("implement me")
}

func (c *ConfigurationManagerImpl) GetFailSafeArmed() bool {
	//TODO implement me
	panic("implement me")
}

func (c *ConfigurationManagerImpl) SetFailSafeArmed(val bool) error {
	//TODO implement me
	panic("implement me")
}

func (c *ConfigurationManagerImpl) GetBLEDeviceIdentificationInfo() (ble.ChipBLEDeviceIdentificationInfo, error) {
	//TODO implement me
	panic("implement me")
}

func (c *ConfigurationManagerImpl) IsFullyProvisioned() bool {
	//TODO implement me
	panic("implement me")
}

func (c *ConfigurationManagerImpl) InitiateFactoryReset() {
	//TODO implement me
	panic("implement me")
}

func (c *ConfigurationManagerImpl) LogDeviceConfig() {
	//TODO implement me
	panic("implement me")
}

func (c *ConfigurationManagerImpl) GetInitialPairingHint() (uint16, error) {
	//TODO implement me
	panic("implement me")
}

func (c *ConfigurationManagerImpl) GetInitialPairingInstruction() (string, error) {
	//TODO implement me
	panic("implement me")
}

func (c *ConfigurationManagerImpl) GetSecondaryPairingHint() (uint16, error) {
	//TODO implement me
	panic("implement me")
}

func (c *ConfigurationManagerImpl) GetSecondaryPairingInstruction() (string, error) {
	//TODO implement me
	panic("implement me")
}

func (c *ConfigurationManagerImpl) GetLocationCapability() (uint8, error) {
	//TODO implement me
	panic("implement me")
}

func (c *ConfigurationManagerImpl) Init(provider ConfigProvider) (*ConfigurationManagerImpl, error) {
	err := provider.EnsureNamespace(KConfigNamespace_ChipConfig)
	if err != nil {
		log.Panic(err.Error())
	}
	err = provider.EnsureNamespace(KConfigNamespace_ChipCounters)
	if err != nil {
		log.Panic(err.Error())
	}
	err = provider.EnsureNamespace(KConfigNamespace_ChipFactory)
	if err != nil {
		log.Panic(err.Error())
	}

	c.ConfigProvider = provider

	if !provider.ConfigValueExists(kConfigKey_VendorId) {
		err := c.StoreProductId(config.ChipDeviceConfigDeviceVendorId)
		if err != nil {
			log.Panic(err.Error())
		}
	}

	if !provider.ConfigValueExists(kConfigKey_ProductId) {
		err := c.StoreProductId(config.ChipDeviceConfigDeviceProductId)
		if err != nil {
			log.Panic(err.Error())
		}
	}

	if provider.ConfigValueExists(kCounterKey_RebootCount) {
		rebootCount, err := c.GetRebootCount()
		if err != nil {
			log.Panic(err.Error())
		}
		err = c.StoreRebootCount(rebootCount + 1)
		if err != nil {
			log.Panic(err.Error())
		}
	} else {
		err := c.StoreRebootCount(1)
		if err != nil {
			log.Panic(err.Error())
		}
	}

	if !provider.ConfigValueExists(kCounterKey_TotalOperationalHours) {
		err := c.StoreTotalOperationalHours(0)
		if err != nil {
			log.Panic(err.Error())
		}
	}

	if !provider.ConfigValueExists(kCounterKey_BootReason) {
		err := c.StoreBootReason(clusters.BootReasonType_kUnspecified)
		if err != nil {
			log.Panic(err.Error())
		}
	}

	if !provider.ConfigValueExists(kConfigKey_RegulatoryLocation) {
		err := provider.WriteConfigValueUint32(kConfigKey_RegulatoryLocation, clusters.RegulatoryLocationType_kIndoor)
		if err != nil {
			log.Panic(err.Error())
		}
	}

	if !provider.ConfigValueExists(kConfigKey_LocationCapability) {
		err := provider.WriteConfigValueUint32(kConfigKey_LocationCapability, clusters.RegulatoryLocationType_kIndoorOutdoor)
		if err != nil {
			log.Panic(err.Error())
		}
	}
	return c, nil
}

func (c *ConfigurationManagerImpl) StoreBootReason(bootReason uint32) error {
	return c.WriteConfigValueUint32(kCounterKey_RebootCount, bootReason)
}

func (c *ConfigurationManagerImpl) StoreTotalOperationalHours(totalOperationalHours uint32) error {
	return c.WriteConfigValueUint32(kCounterKey_TotalOperationalHours, totalOperationalHours)
}

func (c *ConfigurationManagerImpl) GetRebootCount() (uint32, error) {
	//TODO implement me
	panic("implement me")
}

func (c *ConfigurationManagerImpl) StoreRebootCount(rebootCount uint32) error {
	return c.WriteConfigValueUint32(kCounterKey_RebootCount, rebootCount)
}

func (c *ConfigurationManagerImpl) GetPrimaryMACAddress() (string, error) {
	//TODO implement me
	panic("implement me")
}

func (c *ConfigurationManagerImpl) GetPrimaryWiFiMACAddress() ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (c *ConfigurationManagerImpl) GetPrimary802154MACAddress() ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (c *ConfigurationManagerImpl) GetSoftwareVersionString() {
	//TODO implement me
	panic("implement me")
}

func (c *ConfigurationManagerImpl) GetSoftwareVersion() (uint32, error) {
	//TODO implement me
	panic("implement me")
}

func (c *ConfigurationManagerImpl) GetFirmwareBuildChipEpochTime() (time.Duration, error) {
	//TODO implement me
	panic("implement me")
}

func (c *ConfigurationManagerImpl) SetFirmwareBuildChipEpochTime() (time.Duration, error) {
	//TODO implement me
	panic("implement me")
}

func (c *ConfigurationManagerImpl) GetLifetimeCounter() (uint16, error) {
	//TODO implement me
	panic("implement me")
}

func (c *ConfigurationManagerImpl) StoreVendorId(vendorId uint16) error {
	return c.WriteConfigValueUint16(kConfigKey_VendorId, vendorId)
}

func (c *ConfigurationManagerImpl) StoreProductId(productId uint16) error {
	return c.WriteConfigValueUint16(kConfigKey_ProductId, productId)
}

func (c *ConfigurationManagerImpl) RunUnitTests() error {
	//TODO implement me
	panic("implement me")
}

func (c ConfigurationManagerImpl) GetVendorId() (uint16, error) {
	return c.mVendorId, nil
}

func (c ConfigurationManagerImpl) GetSetupDiscriminator() (uint16, error) {
	return c.mDiscriminator, nil
}

func (c ConfigurationManagerImpl) GetVendorName() (string, error) {
	return c.mVendorName, nil
}

func (c ConfigurationManagerImpl) GetProductId() (uint16, error) {
	return c.mProductId, nil
}

func (c ConfigurationManagerImpl) GetProductName() string {
	return c.mProductName
}

func (c ConfigurationManagerImpl) IsCommissionableDeviceTypeEnabled() bool {
	return c.deviceConfigEnableCommissionableDeviceType
}

func (c ConfigurationManagerImpl) GetDeviceTypeId() (uint32, error) {
	return c.mDeviceType, nil
}

func (c ConfigurationManagerImpl) SetDeviceTypeId(t uint32) {
	c.mDeviceType = t
}

func (c ConfigurationManagerImpl) IsCommissionableDeviceNameEnabled() bool {
	return true
}

func (c ConfigurationManagerImpl) GetCommissionableDeviceName() (string, error) {
	return c.mDeviceName, nil
}
