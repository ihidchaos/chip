package device

import (
	"github.com/galenliu/chip/pkg/storage"
	"sync"
)

type InfoProvider interface {
}

type InfoProviderImpl struct {
	storage storage.KvsPersistentStorageDelegate
}

func (i *InfoProviderImpl) SetStorageDelegate(storage storage.KvsPersistentStorageDelegate) {
	i.storage = storage
}

var _deviceInfoProvider *InfoProviderImpl
var _deviceInfoProviderOnce sync.Once

func GetDeviceInfoProvider() *InfoProviderImpl {
	_deviceInfoProviderOnce.Do(func() {
		if _deviceInfoProvider == nil {
			_deviceInfoProvider = &InfoProviderImpl{}
		}
	})
	return _deviceInfoProvider
}
