package store

import (
	"sync/atomic"
)

var defaultKeyValueManager atomic.Value

func init() {
	storeManager := &DefaultPersistentImpl{}
	defaultKeyValueManager.Store(storeManager)
}

func DefaultPersistentStorage() PersistentStorageDelegate {
	storeManager := defaultKeyValueManager.Load().(PersistentStorageDelegate)
	return storeManager
}

func SetDefaultPersistentStorage(impl PersistentStorageDelegate) {
	defaultKeyValueManager.Store(impl)
}
