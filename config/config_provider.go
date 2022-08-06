package config

import (
	"github.com/galenliu/chip/lib"
	"github.com/galenliu/chip/pkg/storage"
	log "github.com/sirupsen/logrus"
	"sync"
)

const KConfigNamespaceChipFactory = "chip-factory"
const KConfigNamespaceChipConfig = "chip-config"
const KConfigNamespaceChipCounters = "chip-counters"

type StorageDelegate interface {
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

type Provider interface {
	StorageDelegate
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
	mChipFactoryStorage  storage.ChipStorage
	mChipConfigStorage   storage.ChipStorage
	mChipCountersStorage storage.ChipStorage
}

func NewConfigProviderImpl() *ConfigProviderImpl {
	return GetConfigProviderInstance()
}

func (conf *ConfigProviderImpl) ReadConfigValueBool(k Key) (bool, error) {
	store := conf.GetStorageForNamespace(k)
	return store.ReadValueBool(k.Name)
}

func (conf *ConfigProviderImpl) ReadConfigValueUint16(k Key) (uint16, error) {
	store := conf.GetStorageForNamespace(k)
	v, err := store.ReadValueUint64(k.Name)
	return uint16(v), err
}

func (conf *ConfigProviderImpl) ReadConfigValueUint32(k Key) (uint32, error) {
	store := conf.GetStorageForNamespace(k)
	v, err := store.ReadValueUint64(k.Name)
	return uint32(v), err
}

func (conf *ConfigProviderImpl) ReadConfigValueUint64(k Key) (uint64, error) {
	store := conf.GetStorageForNamespace(k)
	return store.ReadValueUint64(k.Name)
}

func (conf *ConfigProviderImpl) ReadConfigValueStr(k Key) (string, error) {
	store := conf.GetStorageForNamespace(k)
	v, err := store.ReadValueString(k.Name)
	return v, err
}

func (conf *ConfigProviderImpl) ReadConfigValueBin(k Key) ([]byte, error) {
	store := conf.GetStorageForNamespace(k)
	v, err := store.ReadValueString(k.Name)
	return []byte(v), err
}

func (conf *ConfigProviderImpl) WriteConfigValueBool(k Key, val bool) error {
	store := conf.GetStorageForNamespace(k)
	return store.WriteValueBool(k.Name, val)
}

func (conf *ConfigProviderImpl) WriteConfigValueUint16(k Key, val uint16) error {
	store := conf.GetStorageForNamespace(k)
	return store.WriteValueUint64(k.Name, uint64(val))
}

func (conf *ConfigProviderImpl) WriteConfigValueUint32(k Key, val uint32) error {
	store := conf.GetStorageForNamespace(k)
	return store.WriteValueUint64(k.Name, uint64(val))
}

func (conf *ConfigProviderImpl) WriteConfigValueUint64(k Key, val uint64) error {
	store := conf.GetStorageForNamespace(k)
	return store.WriteValueUint64(k.Name, val)
}

func (conf *ConfigProviderImpl) WriteConfigValueStr(k Key, val string) error {
	store := conf.GetStorageForNamespace(k)
	return store.WriteValueString(k.Name, val)
}

func (conf *ConfigProviderImpl) WriteConfigValueBin(k Key, val []byte) error {
	store := conf.GetStorageForNamespace(k)
	return store.WriteValueString(k.Name, string(val))
}

func (conf *ConfigProviderImpl) ClearConfigValue(k Key) error {
	store := conf.GetStorageForNamespace(k)
	return store.DeleteKeyValue(k.Name)
}

func (conf *ConfigProviderImpl) ConfigValueExists(k Key) bool {
	store := conf.GetStorageForNamespace(k)
	return store.HasValue(k.Name)
}

func (conf *ConfigProviderImpl) FactoryResetConfig() error {
	if conf.mChipFactoryStorage == nil {
		log.Printf("storage get failed")
		return lib.ChipDeviceErrorConfigNotFound
	}
	err := conf.mChipFactoryStorage.DeleteAll()
	if err != nil {
		log.Printf("storage ClearAll failed: %s", err.Error())
		return err
	}
	return nil
}

func (conf *ConfigProviderImpl) FactoryResetCounters() error {
	if conf.mChipCountersStorage == nil {
		log.Printf("storage get failed")
		return lib.ChipDeviceErrorConfigNotFound
	}
	err := conf.mChipCountersStorage.DeleteAll()
	if err != nil {
		log.Printf("storage ClearAll failed: %s", err.Error())
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

func (conf *ConfigProviderImpl) GetStorageForNamespace(k Key) storage.ChipStorage {
	if k.Namespace == KConfigNamespaceChipFactory {
		if conf.mChipFactoryStorage == nil {
			conf.mChipFactoryStorage = storage.NewPersistentStorageImpl()
		}
		return conf.mChipFactoryStorage
	}
	if k.Namespace == KConfigNamespaceChipConfig {
		if conf.mChipConfigStorage == nil {
			conf.mChipConfigStorage = storage.NewPersistentStorageImpl()
		}
		return conf.mChipConfigStorage
	}
	if k.Namespace == KConfigNamespaceChipCounters {
		if conf.mChipCountersStorage == nil {
			conf.mChipCountersStorage = storage.NewPersistentStorageImpl()
		}
		return conf.mChipCountersStorage
	}
	return nil
}

func (conf *ConfigProviderImpl) EnsureNamespace(k string) error {
	if k == KConfigNamespaceChipFactory {
		if conf.mChipFactoryStorage == nil {
			conf.mChipFactoryStorage = storage.NewPersistentStorageImpl()
		}
		log.Infof("init factory storage: %s", ChipDefaultFactoryPath)
		return conf.mChipFactoryStorage.Init(ChipDefaultFactoryPath)
	}
	if k == KConfigNamespaceChipConfig {
		if conf.mChipConfigStorage == nil {
			conf.mChipConfigStorage = storage.NewPersistentStorageImpl()
		}
		log.Infof("init config storage: %s", ChipDefaultConfigPath)
		return conf.mChipConfigStorage.Init(ChipDefaultConfigPath)
	}
	if k == KConfigNamespaceChipCounters {
		if conf.mChipCountersStorage == nil {
			conf.mChipCountersStorage = storage.NewPersistentStorageImpl()
		}
		log.Infof("init counters storage: %s", ChipDefaultDataPath)
		return conf.mChipCountersStorage.Init(ChipDefaultDataPath)
	}
	return nil
}

var (
	KConfigKey_SerialNum             = Key{KConfigNamespaceChipFactory, "serial-num"}
	kConfigKey_MfrDeviceId           = Key{KConfigNamespaceChipFactory, "device-id"}
	kConfigKey_MfrDeviceCert         = Key{KConfigNamespaceChipFactory, "device-cert"}
	kConfigKey_MfrDeviceICACerts     = Key{KConfigNamespaceChipFactory, "device-ca-certs"}
	kConfigKey_MfrDevicePrivateKey   = Key{KConfigNamespaceChipFactory, "device-key"}
	KConfigKey_HardwareVersion       = Key{KConfigNamespaceChipFactory, "hardware-ver"}
	KConfigKey_ManufacturingDate     = Key{KConfigNamespaceChipFactory, "mfg-date"}
	kConfigKey_SetupPinCode          = Key{KConfigNamespaceChipFactory, "pin-code"}
	kConfigKey_SetupDiscriminator    = Key{KConfigNamespaceChipFactory, "discriminator"}
	kConfigKey_Spake2pIterationCount = Key{KConfigNamespaceChipFactory, "iteration-count"}
	kConfigKey_Spake2pSalt           = Key{KConfigNamespaceChipFactory, "salt"}
	kConfigKey_Spake2pVerifier       = Key{KConfigNamespaceChipFactory, "verifier"}
	KConfigKey_VendorId              = Key{KConfigNamespaceChipFactory, "vendor-id"}
	KConfigKey_ProductId             = Key{KConfigNamespaceChipFactory, "product-id"}
	kConfigKey_ServiceConfig         = Key{KConfigNamespaceChipConfig, "service-config"}
	kConfigKey_PairedAccountId       = Key{KConfigNamespaceChipConfig, "account-id"}
	kConfigKey_ServiceId             = Key{KConfigNamespaceChipConfig, "service-id"}
	kConfigKey_LastUsedEpochKeyId    = Key{KConfigNamespaceChipConfig, "last-ek-id"}
	kConfigKey_FailSafeArmed         = Key{KConfigNamespaceChipConfig, "fail-safe-armed"}
	kConfigKey_WiFiStationSecType    = Key{KConfigNamespaceChipConfig, "sta-sec-type"}
	KConfigKey_RegulatoryLocation    = Key{KConfigNamespaceChipConfig, "regulatory-location"}
	kConfigKey_CountryCode           = Key{KConfigNamespaceChipConfig, "country-code"}
	KConfigKey_LocationCapability    = Key{KConfigNamespaceChipConfig, "location-capability"}
	kConfigKey_UniqueId              = Key{KConfigNamespaceChipConfig, "unique-id"}

	KCounterKey_RebootCount           = Key{KConfigNamespaceChipCounters, "reboot-count"}
	kCounterKey_UpTime                = Key{KConfigNamespaceChipCounters, "up-time"}
	KCounterKey_TotalOperationalHours = Key{KConfigNamespaceChipCounters, "total-operational-hours"}
	KCounterKey_BootReason            = Key{KConfigNamespaceChipCounters, "boot-reason"}
)
