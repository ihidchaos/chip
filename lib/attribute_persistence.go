package lib

import (
	"github.com/galenliu/chip/pkg/store"
)

type AttributePersistenceProvider interface {
	Init(storage store.KvsPersistentStorageBase) error
}

type AttributePersistence struct {
	mStorage store.KvsPersistentStorageBase
}

func NewAttributePersistence() *AttributePersistence {

	return &AttributePersistence{}
}

func (p AttributePersistence) Init(storage store.KvsPersistentStorageBase) (err error) {
	p.mStorage = storage
	return
}
