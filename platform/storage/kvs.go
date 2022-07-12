package storage

type PersistentStorageDelegate interface {
	SyncGetKeyValue(key string) (any, error)
	SyncSetKeyValue(key string, value any) error
	SyncDeleteKeyValue(key string) error
	SyncDoesKeyExist(key string) bool
}

type KvsPersistentStorageImpl struct {
	StorageImpl
}

func (k KvsPersistentStorageImpl) SyncGetKeyValue(key string) (any, error) {
	//TODO implement me
	panic("implement me")
}

func (k KvsPersistentStorageImpl) SyncSetKeyValue(key string, value any) error {
	//TODO implement me
	panic("implement me")
}

func (k KvsPersistentStorageImpl) SyncDeleteKeyValue(key string) error {
	//TODO implement me
	panic("implement me")
}

func (k KvsPersistentStorageImpl) SyncDoesKeyExist(key string) bool {
	//TODO implement me
	panic("implement me")
}
