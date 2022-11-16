package config

import (
	"fmt"
	"github.com/galenliu/chip/ble"
	"github.com/galenliu/chip/clusters"
	log "golang.org/x/exp/slog"
	"math/rand"
	"sync/atomic"
	"time"
)

// ConfigurationManager Defines the public interface for the Device Layer ConfigurationManager object.
type ConfigurationManager interface {
	StorageDelegate
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
	GetBootReason() (uint32, error)
	StoreBootReason(bootReason uint32) error
	GetPartNumber() (string, error)
	GetProductURL() (string, error)
	GetProductLabel() (string, error)
	GetUniqueId() (string, error)
	StoreUniqueId(uniqueId string) error
	GenerateUniqueId() error
	GetFailSafeArmed() bool
	SetFailSafeArmed(val bool) error

	GetBLEDeviceIdentificationInfo() (ble.DeviceIdentificationInfo, error)

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
	Provider
}

var defaultManager atomic.Value

func init() {
	defaultManager.Store(NewConfigurationManagerImpl())
}

func DefaultManager() *ConfigurationManagerImpl {
	return defaultManager.Load().(*ConfigurationManagerImpl)
}

func NewConfigurationManagerImpl() *ConfigurationManagerImpl {
	return &ConfigurationManagerImpl{}
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

func (c *ConfigurationManagerImpl) GetBootReason() (uint32, error) {
	return c.Provider.ReadConfigValueUint32(KCounterKey_BootReason)
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

func (c *ConfigurationManagerImpl) GetBLEDeviceIdentificationInfo() (ble.DeviceIdentificationInfo, error) {
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
	return 0, nil
}

func (c *ConfigurationManagerImpl) GetInitialPairingInstruction() (string, error) {
	return PairingInitialInstruction, nil
}

func (c *ConfigurationManagerImpl) GetSecondaryPairingHint() (uint16, error) {
	//TODO implement me
	panic("implement me")
}

func (c *ConfigurationManagerImpl) GetSecondaryPairingInstruction() (string, error) {
	return PairingSecondaryInstruction, nil
}

func (c *ConfigurationManagerImpl) GetLocationCapability() (uint8, error) {
	//TODO implement me
	panic("implement me")
}

func (c *ConfigurationManagerImpl) Init(configProvider Provider, options *DeviceOptions) error {
	err := configProvider.EnsureNamespace(KConfigNamespaceChipConfig)
	if err != nil {
		log.Warn(err.Error())
	}
	err = configProvider.EnsureNamespace(KConfigNamespaceChipCounters)
	if err != nil {
		log.Warn(err.Error())
	}
	err = configProvider.EnsureNamespace(KConfigNamespaceChipFactory)
	if err != nil {
		log.Warn(err.Error())
	}

	c.Provider = configProvider

	// if  vendorId set, stored!
	if options.Payload.VendorID != 0 {
		err := c.StoreVendorId(options.Payload.VendorID)
		if err != nil {
			return err
		}
	} else {
		if !configProvider.ConfigValueExists(KConfigKey_VendorId) {
			err := c.StoreVendorId(DeviceVendorId)
			if err != nil {
				return err
			}
		}
	}

	if options.Payload.ProductID != 0 {
		err := c.StoreVendorId(options.Payload.VendorID)
		if err != nil {
			return err
		}
	} else {
		if !configProvider.ConfigValueExists(KConfigKey_ProductId) {
			err := c.StoreProductId(DeviceProductId)
			if err != nil {
				log.Info(err.Error())
			}
		}
	}

	if configProvider.ConfigValueExists(KCounterKey_RebootCount) {
		rebootCount, err := c.GetRebootCount()
		if err != nil {
			return err
		}
		err = c.StoreRebootCount(rebootCount + 1)
		if err != nil {
			return err
		}
	} else {
		err := c.StoreRebootCount(1)
		if err != nil {
			return err
		}
	}

	if !configProvider.ConfigValueExists(KCounterKey_TotalOperationalHours) {
		err := c.StoreTotalOperationalHours(0)
		if err != nil {
			return err
		}
	}

	if !configProvider.ConfigValueExists(KCounterKey_BootReason) {
		err := c.StoreBootReason(clusters.BootreasontypeKunspecified)
		if err != nil {
			return err
		}
	}

	if !configProvider.ConfigValueExists(KConfigKey_RegulatoryLocation) {
		err := configProvider.WriteConfigValueUint32(KConfigKey_RegulatoryLocation, clusters.RegulatorylocationtypeKindoor)
		if err != nil {
			return err
		}
	}

	if !configProvider.ConfigValueExists(KConfigKey_LocationCapability) {
		err := configProvider.WriteConfigValueUint32(KConfigKey_LocationCapability, clusters.RegulatorylocationtypeKindooroutdoor)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *ConfigurationManagerImpl) StoreBootReason(bootReason uint32) error {
	return c.WriteConfigValueUint32(KCounterKey_BootReason, bootReason)
}

func (c *ConfigurationManagerImpl) StoreTotalOperationalHours(totalOperationalHours uint32) error {
	return c.WriteConfigValueUint32(KCounterKey_TotalOperationalHours, totalOperationalHours)
}

func (c *ConfigurationManagerImpl) GetRebootCount() (uint32, error) {
	return c.Provider.ReadConfigValueUint32(KCounterKey_RebootCount)
}

func (c *ConfigurationManagerImpl) StoreRebootCount(rebootCount uint32) error {
	return c.WriteConfigValueUint32(KCounterKey_RebootCount, rebootCount)
}

func (c *ConfigurationManagerImpl) GetPrimaryMACAddress() (string, error) {
	return fmt.Sprintf("%016X", rand.Uint64()), nil
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
	return c.WriteConfigValueUint16(KConfigKey_VendorId, vendorId)
}

func (c *ConfigurationManagerImpl) StoreProductId(productId uint16) error {
	return c.WriteConfigValueUint16(KConfigKey_ProductId, productId)
}

func (c *ConfigurationManagerImpl) RunUnitTests() error {
	//TODO implement me
	panic("implement me")
}

func (c *ConfigurationManagerImpl) GetVendorId() (uint16, error) {
	return c.Provider.ReadConfigValueUint16(KConfigKey_VendorId)
}

func (c *ConfigurationManagerImpl) GetSetupDiscriminator() (uint16, error) {
	return c.mDiscriminator, nil
}

func (c *ConfigurationManagerImpl) GetVendorName() (string, error) {
	return c.mVendorName, nil
}

func (c *ConfigurationManagerImpl) GetProductId() (uint16, error) {
	return c.Provider.ReadConfigValueUint16(KConfigKey_ProductId)
}

func (c *ConfigurationManagerImpl) GetProductName() string {
	return c.mProductName
}

func (c *ConfigurationManagerImpl) IsCommissionableDeviceTypeEnabled() bool {
	return c.deviceConfigEnableCommissionableDeviceType
}

func (c *ConfigurationManagerImpl) GetDeviceTypeId() (uint32, error) {
	return c.mDeviceType, nil
}

func (c *ConfigurationManagerImpl) SetDeviceTypeId(t uint32) {
	c.mDeviceType = t
}

func (c *ConfigurationManagerImpl) IsCommissionableDeviceNameEnabled() bool {
	return EnableCommissionableDeviceName == 1
}

func (c *ConfigurationManagerImpl) GetCommissionableDeviceName() (string, error) {
	return DeviceName, nil
}
