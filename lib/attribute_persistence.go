package lib

import (
	"github.com/galenliu/chip/pkg/storage"
)

type AttributePersistenceProvider interface {
	Init(storage storage.KvsPersistentStorageDelegate) error
}

type AttributePersistence struct {
	mStorage storage.KvsPersistentStorageDelegate
}

func NewAttributePersistence() *AttributePersistence {

	return &AttributePersistence{}
}

func (p AttributePersistence) Init(storage storage.KvsPersistentStorageDelegate) (err error) {
	p.mStorage = storage
	return
}
