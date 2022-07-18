package storage

import "sync"

type KeyValueStoreManager struct {
	KvsPersistentStorageImpl
}

var _instance *KeyValueStoreManager
var once sync.Once

func KeyValueStoreMgr() *KeyValueStoreManager {
	once.Do(func() {
		_instance = &KeyValueStoreManager{
			KvsPersistentStorageImpl: KvsPersistentStorageImpl{
				PersistentStorageImpl: PersistentStorageImpl{
					storage:      &iniStorageImpl{},
					mConfigPath:  "",
					mLock:        new(sync.Mutex),
					mDirty:       false,
					mInitialized: false,
				},
			},
		}
	})
	return _instance
}
