package store

import (
	"sync/atomic"
)

type KeyValueManager struct {
	KvsPersistentStorageBase
}

var defaultKeyValueManager atomic.Value

func init() {
	storeManager := NewPersistentStorageImpl()
	defaultKeyValueManager.Store(storeManager)
}

func DefaultKeyValueMgr() *KeyValueManager {
	storeManager := defaultKeyValueManager.Load().(*KeyValueManager)
	return storeManager
}

func SetDefaultKeyValueManager(impl *KeyValueManager) {
	defaultKeyValueManager.Store(impl)
}

func (m *KeyValueManager) Init(file string) error {
	storage := NewPersistentStorageImpl()
	if err := storage.Init(file); err != nil {
		return err
	}
	kvsStorage := NewKvsPersistentStorageImpl()
	kvsStorage.Init(storage)
	m.KvsPersistentStorageBase = kvsStorage
	return nil
}
