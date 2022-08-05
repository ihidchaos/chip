package storage

import "github.com/spf13/cast"

type KeyValuePersistentStorage interface {
	SyncGetKeyValue(key string) (any, error)
	SyncSetKeyValue(key string, value any) error
	SyncDeleteKeyValue(key string) error
	SyncDoesKeyExist(key string) bool
	PersistentStorage
}

type KeyValuePersistentStorageImpl struct {
	*PersistentStorageImpl
}

func NewKeyValuePersistentStorage() *KeyValuePersistentStorageImpl {
	impl := &KeyValuePersistentStorageImpl{
		PersistentStorageImpl: NewPersistentStorageImpl(),
	}
	return impl
}

func (k KeyValuePersistentStorageImpl) SyncGetKeyValue(key string) (any, error) {
	return k.PersistentStorageImpl.ReadValueStr(key)
}

func (k KeyValuePersistentStorageImpl) SyncSetKeyValue(key string, value any) error {
	switch value.(type) {
	case string:
		return k.PersistentStorageImpl.WriteValueStr(key, cast.ToString(value))
	case uint16, uint32, uint64, uint8, uint:
		return k.PersistentStorageImpl.WriteValueUint64(key, cast.ToUint64(value))
	case bool:
		return k.PersistentStorageImpl.WriteValueBool(key, cast.ToBool(value))
	default:
		return k.PersistentStorageImpl.WriteValueStr(key, cast.ToString(value))
	}
}

func (k KeyValuePersistentStorageImpl) SyncDeleteKeyValue(key string) error {
	return k.PersistentStorageImpl.ClearValue(key)
}

func (k KeyValuePersistentStorageImpl) SyncDoesKeyExist(key string) bool {
	return k.PersistentStorageImpl.HasValue(key)
}
