package device

import (
	"github.com/galenliu/chip/storage"
	"sync"
)

type DeviceInfoProvider interface {
}

type DeviceInfoProviderImpl struct {
	storage storage.StorageDelegate
}

func (i *DeviceInfoProviderImpl) SetStorageDelegate(storage storage.StorageDelegate) {
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
