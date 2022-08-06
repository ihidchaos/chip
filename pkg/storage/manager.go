package storage

import "sync"

type KeyValueStoreManagerImpl struct {
	*ChipStorageImpl
}

var _instance *KeyValueStoreManagerImpl
var once sync.Once

func KeyValueStoreMgr() *KeyValueStoreManagerImpl {
	once.Do(func() {
		_instance = &KeyValueStoreManagerImpl{
			ChipStorageImpl: NewPersistentStorageImpl(),
		}
	})
	return _instance
}

func (m *KeyValueStoreManagerImpl) Init(file string) error {
	return m.ChipStorageImpl.Init(file)
}
