package lib

import "github.com/galenliu/chip/platform/storage"

type AttributePersistenceProvider interface {
	Init(storage storage.PersistentStorageDelegate) error
}

type AttributePersistence struct {
	mStorage storage.PersistentStorageDelegate
}

func (p AttributePersistence) Init(storage storage.PersistentStorageDelegate) (err error) {
	p.mStorage = storage
	return
}
