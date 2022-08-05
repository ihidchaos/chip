package device

import (
	"github.com/galenliu/chip/pkg/storage"
	"sync"
)

type DeviceInfoProvider interface {
}

type DeviceInfoProviderImpl struct {
	storage storage.PersistentStorage
}

func (i *DeviceInfoProviderImpl) SetStorageDelegate(storage storage.PersistentStorage) {
	i.storage = storage
}

var _deviceInfoProvider *DeviceInfoProviderImpl
var _deviceInfoProviderOnce sync.Once

func GetDeviceInfoProvider() *DeviceInfoProviderImpl {
	_deviceInfoProviderOnce.Do(func() {
		if _deviceInfoProvider == nil {
			_deviceInfoProvider = &DeviceInfoProviderImpl{}
		}
	})
	return _deviceInfoProvider
}
