package device

import (
	"github.com/galenliu/chip/lib/store"
	"sync/atomic"
)

type InfoProvider interface {
}

type InfoProviderImpl struct {
	storage store.KvsPersistentStorageBase
}

func (i *InfoProviderImpl) SetStorageDelegate(storage store.KvsPersistentStorageBase) {
	i.storage = storage
}

var defaultInfoProvider atomic.Value

func init() {
	info := &InfoProviderImpl{}
	defaultInfoProvider.Store(info)
}

func DefaultInfoProvider() *InfoProviderImpl {
	_deviceInfoProvider := defaultInfoProvider.Load().(*InfoProviderImpl)
	return _deviceInfoProvider
}

func SetDefaultInfoProvider(provider *InfoProviderImpl) {
	defaultInfoProvider.Store(provider)
}
