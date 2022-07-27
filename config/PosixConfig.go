package config

import (
	"github.com/galenliu/chip/lib"
	"github.com/galenliu/chip/storage"
	log "github.com/sirupsen/logrus"
	"sync"
)

const KConfigNamespace_ChipFactory = "chip-factory"
const KConfigNamespace_ChipConfig = "chip-config"
const KConfigNamespace_ChipCounters = "chip-counters"

type ConfigDelegate interface {
	ReadConfigValueBool(k Key) (bool, error)
	ReadConfigValueUint16(k Key) (uint16, error)
	ReadConfigValueUint32(k Key) (uint32, error)
	ReadConfigValueUint64(k Key) (uint64, error)
	ReadConfigValueStr(k Key) (string, error)
	ReadConfigValueBin(k Key) ([]byte, error)

	WriteConfigValueBool(k Key, val bool) error
	WriteConfigValueUint16(k Key, val uint16) error
	WriteConfigValueUint32(k Key, val uint32) error
	WriteConfigValueUint64(k Key, val uint64) error
	WriteConfigValueStr(k Key, val string) error
	WriteConfigValueBin(k Key, data []byte) error
}

type ConfigProvider interface {
	ConfigDelegate
	ClearConfigValue(k Key) error
	ConfigValueExists(k Key) bool
	FactoryResetConfig() error
	FactoryResetCounters() error
	RunConfigUnitTest()
	EnsureNamespace(k string) error
}

var _ConfigProviderInstance *ConfigProviderImpl
var _ConfigProviderInstanceOnce sync.Once

func GetConfigProviderInstance() *ConfigProviderImpl {
	_ConfigProviderInstanceOnce.Do(func() {
		if _ConfigProviderInstance == nil {
			_ConfigProviderInstance = &ConfigProviderImpl{}
		}
	})
	return _ConfigProviderInstance
}

type ConfigProviderImpl struct {
	mChipFactoryStorage  storage.PersistentStorageDelegate
	mChipConfigStorage   storage.PersistentStorageDelegate
	mChipCountersStorage storage.PersistentStorageDelegate
}

func NewConfigProviderImpl() *ConfigProviderImpl {
	return GetConfigProviderInstance()
}

func (conf *ConfigProviderImpl) ReadConfigValueBool(k Key) (bool, error) {
	store := conf.GetStorageForNamespace(k)
	return store.ReadBoolValue(k.Name)
}

func (conf *ConfigProviderImpl) ReadConfigValueUint16(k Key) (uint16, error) {
	store := conf.GetStorageForNamespace(k)
	return store.ReadValueUint16(k.Name)
}

func (conf *ConfigProviderImpl) ReadConfigValueUint32(k Key) (uint32, error) {
	store := conf.GetStorageForNamespace(k)
	return store.ReadValueUint32(k.Name)
}

func (conf *ConfigProviderImpl) ReadConfigValueUint64(k Key) (uint64, error) {
	store := conf.GetStorageForNamespace(k)
	return store.ReadValueUint64(k.Name)
}

func (conf *ConfigProviderImpl) ReadConfigValueStr(k Key) (string, error) {
	store := conf.GetStorageForNamespace(k)
	return store.ReadValueStr(k.Name)
}

func (conf *ConfigProviderImpl) ReadConfigValueBin(k Key) ([]byte, error) {
	store := conf.GetStorageForNamespace(k)
	return store.ReadValueBin(k.Name)
}

func (conf *ConfigProviderImpl) WriteConfigValueBool(k Key, val bool) error {
	store := conf.GetStorageForNamespace(k)
	err := store.WriteValueBool(k.Name, val)
	if err != nil {
		return err
	}
	return store.Commit()
}

func (conf *ConfigProviderImpl) WriteConfigValueUint16(k Key, val uint16) error {
	store := conf.GetStorageForNamespace(k)
	err := store.WriteValueUint16(k.Name, val)
	if err != nil {
		return err
	}
	return store.Commit()

}

func (conf *ConfigProviderImpl) WriteConfigValueUint32(k Key, val uint32) error {
	store := conf.GetStorageForNamespace(k)
	err := store.WriteValueUint32(k.Name, val)
	if err != nil {
		return err
	}
	return store.Commit()
}

func (conf *ConfigProviderImpl) WriteConfigValueUint64(k Key, val uint64) error {
	store := conf.GetStorageForNamespace(k)
	err := store.WriteValueUint64(k.Name, val)
	if err != nil {
		return err
	}
	return store.Commit()
}

func (conf *ConfigProviderImpl) WriteConfigValueStr(k Key, val string) error {
	store := conf.GetStorageForNamespace(k)
	err := store.WriteValueStr(k.Name, val)
	if err != nil {
		return err
	}
	return store.Commit()
}

func (conf *ConfigProviderImpl) WriteConfigValueBin(k Key, val []byte) error {
	store := conf.GetStorageForNamespace(k)
	err := store.WriteValueBin(k.Name, val)
	if err != nil {
		return err
	}
	return store.Commit()
}

func (conf *ConfigProviderImpl) ClearConfigValue(k Key) error {
	store := conf.GetStorageForNamespace(k)
	return store.ClearValue(k.Name)
}

func (conf *ConfigProviderImpl) ConfigValueExists(k Key) bool {
	store := conf.GetStorageForNamespace(k)
	return store.HasValue(k.Name)
}

func (conf *ConfigProviderImpl) FactoryResetConfig() error {
	if conf.mChipFactoryStorage == nil {
		log.Info("Storage get failed")
		return lib.ChipDeviceErrorConfigNotFound
	}
	err := conf.mChipFactoryStorage.ClearAll()
	if err != nil {
		log.Info("Storage ClearAll failed: %s", err.Error())
		return err
	}
	err = conf.mChipFactoryStorage.Commit()
	if err != nil {
		log.Info("Storage Commit failed: %s", err.Error())
		return err
	}
	return nil
}

func (conf *ConfigProviderImpl) FactoryResetCounters() error {
	if conf.mChipCountersStorage == nil {
		log.Info("Storage get failed")
		return lib.ChipDeviceErrorConfigNotFound
	}
	err := conf.mChipCountersStorage.ClearAll()
	if err != nil {
		log.Info("Storage ClearAll failed: %s", err.Error())
		return err
	}
	err = conf.mChipCountersStorage.Commit()
	if err != nil {
		log.Info("Storage Commit failed: %s", err.Error())
		return err
	}
	return nil
}

func (conf *ConfigProviderImpl) RunConfigUnitTest() {
	//TODO implement me
	panic("implement me")
}

type Key struct {
	Namespace string
	Name      string
}

func (conf *ConfigProviderImpl) GetStorageForNamespace(k Key) storage.PersistentStorageDelegate {
	if k.Namespace == KConfigNamespace_ChipFactory {
		if conf.mChipFactoryStorage == nil {
			conf.mChipFactoryStorage = storage.NewPersistentStorageImpl()
		}
		return conf.mChipFactoryStorage
	}
	if k.Namespace == KConfigNamespace_ChipConfig {
		if conf.mChipConfigStorage == nil {
			conf.mChipConfigStorage = storage.NewPersistentStorageImpl()
		}
		return conf.mChipConfigStorage
	}
	if k.Namespace == KConfigNamespace_ChipCounters {
		if conf.mChipCountersStorage == nil {
			conf.mChipCountersStorage = storage.NewPersistentStorageImpl()
		}
		return conf.mChipCountersStorage
	}
	return nil
}

func (conf *ConfigProviderImpl) EnsureNamespace(k string) error {
	if k == KConfigNamespace_ChipFactory {
		if conf.mChipFactoryStorage == nil {
			conf.mChipFactoryStorage = storage.NewPersistentStorageImpl()
		}
		return conf.mChipFactoryStorage.Init(ChipDefaultFactoryPath)
	}
	if k == KConfigNamespace_ChipConfig {
		if conf.mChipConfigStorage == nil {
			conf.mChipConfigStorage = storage.NewPersistentStorageImpl()
		}
		return conf.mChipConfigStorage.Init(ChipDefaultConfigPath)
	}
	if k == KConfigNamespace_ChipFactory {
		if conf.mChipCountersStorage == nil {
			conf.mChipCountersStorage = storage.NewPersistentStorageImpl()
		}
		return conf.mChipCountersStorage.Init(ChipDefaultDataPath)
	}
	return nil
}

var (
	KConfigKey_SerialNum             = Key{KConfigNamespace_ChipFactory, "serial-num"}
	kConfigKey_MfrDeviceId           = Key{KConfigNamespace_ChipFactory, "device-id"}
	kConfigKey_MfrDeviceCert         = Key{KConfigNamespace_ChipFactory, "device-cert"}
	kConfigKey_MfrDeviceICACerts     = Key{KConfigNamespace_ChipFactory, "device-ca-certs"}
	kConfigKey_MfrDevicePrivateKey   = Key{KConfigNamespace_ChipFactory, "device-key"}
	KConfigKey_HardwareVersion       = Key{KConfigNamespace_ChipFactory, "hardware-ver"}
	KConfigKey_ManufacturingDate     = Key{KConfigNamespace_ChipFactory, "mfg-date"}
	kConfigKey_SetupPinCode          = Key{KConfigNamespace_ChipFactory, "pin-code"}
	kConfigKey_SetupDiscriminator    = Key{KConfigNamespace_ChipFactory, "discriminator"}
	kConfigKey_Spake2pIterationCount = Key{KConfigNamespace_ChipFactory, "iteration-count"}
	kConfigKey_Spake2pSalt           = Key{KConfigNamespace_ChipFactory, "salt"}
	kConfigKey_Spake2pVerifier       = Key{KConfigNamespace_ChipFactory, "verifier"}
	KConfigKey_VendorId              = Key{KConfigNamespace_ChipFactory, "vendor-id"}
	KConfigKey_ProductId             = Key{KConfigNamespace_ChipFactory, "product-id"}
	kConfigKey_ServiceConfig         = Key{KConfigNamespace_ChipConfig, "service-config"}
	kConfigKey_PairedAccountId       = Key{KConfigNamespace_ChipConfig, "account-id"}
	kConfigKey_ServiceId             = Key{KConfigNamespace_ChipConfig, "service-id"}
	kConfigKey_LastUsedEpochKeyId    = Key{KConfigNamespace_ChipConfig, "last-ek-id"}
	kConfigKey_FailSafeArmed         = Key{KConfigNamespace_ChipConfig, "fail-safe-armed"}
	kConfigKey_WiFiStationSecType    = Key{KConfigNamespace_ChipConfig, "sta-sec-type"}
	KConfigKey_RegulatoryLocation    = Key{KConfigNamespace_ChipConfig, "regulatory-location"}
	kConfigKey_CountryCode           = Key{KConfigNamespace_ChipConfig, "country-code"}
	KConfigKey_LocationCapability    = Key{KConfigNamespace_ChipConfig, "location-capability"}
	kConfigKey_UniqueId              = Key{KConfigNamespace_ChipConfig, "unique-id"}

	KCounterKey_RebootCount           = Key{KConfigNamespace_ChipCounters, "reboot-count"}
	kCounterKey_UpTime                = Key{KConfigNamespace_ChipCounters, "up-time"}
	KCounterKey_TotalOperationalHours = Key{KConfigNamespace_ChipCounters, "total-operational-hours"}
	KCounterKey_BootReason            = Key{KConfigNamespace_ChipCounters, "boot-reason"}
)
