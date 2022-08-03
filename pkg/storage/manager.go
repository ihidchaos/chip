package storage

import "sync"

type KeyValueStoreManager struct {
	*KvsPersistentStorageImpl
}

var _instance *KeyValueStoreManager
var once sync.Once

func KeyValueStoreMgr() *KeyValueStoreManager {
	once.Do(func() {
		_instance = &KeyValueStoreManager{
			KvsPersistentStorageImpl: NewKvsPersistentStorage(),
		}
	})
	return _instance
}
