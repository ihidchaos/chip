package storage

import "sync"

type KeyValueStoreManager struct {
	*KeyValuePersistentStorageImpl
}

var _instance *KeyValueStoreManager
var once sync.Once

func KeyValueStoreMgr() *KeyValueStoreManager {
	once.Do(func() {
		_instance = &KeyValueStoreManager{
			KeyValuePersistentStorageImpl: NewKeyValuePersistentStorage(),
		}
	})
	return _instance
}
