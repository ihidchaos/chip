package platform

import (
	"github.com/galenliu/chip/config"
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

var _ConfigProviderInstance *PosixConfigImpl
var _ConfigProviderInstanceOnce sync.Once

func GetConfigProviderInstance() *PosixConfigImpl {
	_ConfigProviderInstanceOnce.Do(func() {
		if _ConfigProviderInstance == nil {
			_ConfigProviderInstance = &PosixConfigImpl{}
		}
	})
	return _ConfigProviderInstance
}

type PosixConfigImpl struct {
	mChipFactoryStorage  storage.PersistentStorageDelegate
	mChipConfigStorage   storage.PersistentStorageDelegate
	mChipCountersStorage storage.PersistentStorageDelegate
}

func (conf *PosixConfigImpl) ReadConfigValueBool(k Key) (bool, error) {
	store := conf.GetStorageForNamespace(k)
	return store.ReadBoolValue(k.Name)
}

func (conf *PosixConfigImpl) ReadConfigValueUint16(k Key) (uint16, error) {
	store := conf.GetStorageForNamespace(k)
	return store.ReadValueUint16(k.Name)
}

func (conf *PosixConfigImpl) ReadConfigValueUint32(k Key) (uint32, error) {
	store := conf.GetStorageForNamespace(k)
	return store.ReadValueUint32(k.Name)
}

func (conf *PosixConfigImpl) ReadConfigValueUint64(k Key) (uint64, error) {
	store := conf.GetStorageForNamespace(k)
	return store.ReadValueUint64(k.Name)
}

func (conf *PosixConfigImpl) ReadConfigValueStr(k Key) (string, error) {
	store := conf.GetStorageForNamespace(k)
	return store.ReadValueStr(k.Name)
}

func (conf *PosixConfigImpl) ReadConfigValueBin(k Key) ([]byte, error) {
	store := conf.GetStorageForNamespace(k)
	return store.ReadValueBin(k.Name)
}

func (conf *PosixConfigImpl) WriteConfigValueBool(k Key, val bool) error {
	store := conf.GetStorageForNamespace(k)
	err := store.WriteValueBool(k.Name, val)
	if err != nil {
		return err
	}
	return store.Commit()
}

func (conf *PosixConfigImpl) WriteConfigValueUint16(k Key, val uint16) error {
	store := conf.GetStorageForNamespace(k)
	err := store.WriteValueUint16(k.Name, val)
	if err != nil {
		return err
	}
	return store.Commit()

}

func (conf *PosixConfigImpl) WriteConfigValueUint32(k Key, val uint32) error {
	store := conf.GetStorageForNamespace(k)
	err := store.WriteValueUint32(k.Name, val)
	if err != nil {
		return err
	}
	return store.Commit()
}

func (conf *PosixConfigImpl) WriteConfigValueUint64(k Key, val uint64) error {
	store := conf.GetStorageForNamespace(k)
	err := store.WriteValueUint64(k.Name, val)
	if err != nil {
		return err
	}
	return store.Commit()
}

func (conf *PosixConfigImpl) WriteConfigValueStr(k Key, val string) error {
	store := conf.GetStorageForNamespace(k)
	err := store.WriteValueStr(k.Name, val)
	if err != nil {
		return err
	}
	return store.Commit()
}

func (conf *PosixConfigImpl) WriteConfigValueBin(k Key, val []byte) error {
	store := conf.GetStorageForNamespace(k)
	err := store.WriteValueBin(k.Name, val)
	if err != nil {
		return err
	}
	return store.Commit()
}

func (conf *PosixConfigImpl) ClearConfigValue(k Key) error {
	store := conf.GetStorageForNamespace(k)
	return store.ClearValue(k.Name)
}

func (conf *PosixConfigImpl) ConfigValueExists(k Key) bool {
	store := conf.GetStorageForNamespace(k)
	return store.HasValue(k.Name)
}

func (conf *PosixConfigImpl) FactoryResetConfig() error {
	if conf.mChipFactoryStorage == nil {
		log.Info("Storage get failed")
		return lib.CHIP_DEVICE_ERROR_CONFIG_NOT_FOUND
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

func (conf *PosixConfigImpl) FactoryResetCounters() error {
	if conf.mChipCountersStorage == nil {
		log.Info("Storage get failed")
		return lib.CHIP_DEVICE_ERROR_CONFIG_NOT_FOUND
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

func (conf *PosixConfigImpl) RunConfigUnitTest() {
	//TODO implement me
	panic("implement me")
}

type Key struct {
	Namespace string
	Name      string
}

func (conf *PosixConfigImpl) GetStorageForNamespace(k Key) storage.PersistentStorageDelegate {
	if k.Namespace == KConfigNamespace_ChipFactory {
		if conf.mChipFactoryStorage == nil {
			conf.mChipFactoryStorage = &storage.PersistentStorageImpl{}
		}
		return conf.mChipFactoryStorage
	}
	if k.Namespace == KConfigNamespace_ChipConfig {
		if conf.mChipConfigStorage == nil {
			conf.mChipConfigStorage = &storage.PersistentStorageImpl{}
		}
		return conf.mChipConfigStorage
	}
	if k.Namespace == KConfigNamespace_ChipFactory {
		if conf.mChipCountersStorage == nil {
			conf.mChipCountersStorage = &storage.PersistentStorageImpl{}
		}
		return conf.mChipCountersStorage
	}
	return nil
}

func (conf *PosixConfigImpl) EnsureNamespace(k string) error {
	if k == KConfigNamespace_ChipFactory {
		if conf.mChipFactoryStorage == nil {
			conf.mChipFactoryStorage = &storage.PersistentStorageImpl{}
		}
		return conf.mChipFactoryStorage.Init(config.CHIP_DEFAULT_FACTORY_PATH)
	}
	if k == KConfigNamespace_ChipConfig {
		if conf.mChipConfigStorage == nil {
			conf.mChipConfigStorage = &storage.PersistentStorageImpl{}
		}
		return conf.mChipConfigStorage.Init(config.CHIP_DEFAULT_CONFIG_PATH)
	}
	if k == KConfigNamespace_ChipFactory {
		if conf.mChipCountersStorage == nil {
			conf.mChipCountersStorage = &storage.PersistentStorageImpl{}
		}
		return conf.mChipCountersStorage.Init(config.CHIP_DEFAULT_DATA_PATH)
	}
	return nil
}

var (
	kConfigKey_SerialNum             = Key{KConfigNamespace_ChipFactory, "serial-num"}
	kConfigKey_MfrDeviceId           = Key{KConfigNamespace_ChipFactory, "device-id"}
	kConfigKey_MfrDeviceCert         = Key{KConfigNamespace_ChipFactory, "device-cert"}
	kConfigKey_MfrDeviceICACerts     = Key{KConfigNamespace_ChipFactory, "device-ca-certs"}
	kConfigKey_MfrDevicePrivateKey   = Key{KConfigNamespace_ChipFactory, "device-key"}
	kConfigKey_HardwareVersion       = Key{KConfigNamespace_ChipFactory, "hardware-ver"}
	kConfigKey_ManufacturingDate     = Key{KConfigNamespace_ChipFactory, "mfg-date"}
	kConfigKey_SetupPinCode          = Key{KConfigNamespace_ChipFactory, "pin-code"}
	kConfigKey_SetupDiscriminator    = Key{KConfigNamespace_ChipFactory, "discriminator"}
	kConfigKey_Spake2pIterationCount = Key{KConfigNamespace_ChipFactory, "iteration-count"}
	kConfigKey_Spake2pSalt           = Key{KConfigNamespace_ChipFactory, "salt"}
	kConfigKey_Spake2pVerifier       = Key{KConfigNamespace_ChipFactory, "verifier"}
	kConfigKey_VendorId              = Key{KConfigNamespace_ChipFactory, "vendor-id"}
	kConfigKey_ProductId             = Key{KConfigNamespace_ChipFactory, "product-i=d"}
	kConfigKey_ServiceConfig         = Key{KConfigNamespace_ChipConfig, "service-config"}
	kConfigKey_PairedAccountId       = Key{KConfigNamespace_ChipConfig, "account-id"}
	kConfigKey_ServiceId             = Key{KConfigNamespace_ChipConfig, "service-id"}
	kConfigKey_LastUsedEpochKeyId    = Key{KConfigNamespace_ChipConfig, "last-ek-id"}
	kConfigKey_FailSafeArmed         = Key{KConfigNamespace_ChipConfig, "fail-safe-armed"}
	kConfigKey_WiFiStationSecType    = Key{KConfigNamespace_ChipConfig, "sta-sec-type"}
	kConfigKey_RegulatoryLocation    = Key{KConfigNamespace_ChipConfig, "regulatory-location"}
	kConfigKey_CountryCode           = Key{KConfigNamespace_ChipConfig, "country-code"}
	kConfigKey_LocationCapability    = Key{KConfigNamespace_ChipConfig, "location-capability"}
	kConfigKey_UniqueId              = Key{KConfigNamespace_ChipConfig, "unique-id"}

	kCounterKey_RebootCount           = Key{KConfigNamespace_ChipCounters, "reboot-count"}
	kCounterKey_UpTime                = Key{KConfigNamespace_ChipCounters, "up-time"}
	kCounterKey_TotalOperationalHours = Key{KConfigNamespace_ChipCounters, "total-operational-hours"}
	kCounterKey_BootReason            = Key{KConfigNamespace_ChipCounters, "boot-reason"}
)
