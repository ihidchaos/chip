package lib

import (
	"github.com/galenliu/chip/lib/store"
)

type AttributePersistenceProvider interface {
	Init(storage store.PersistentStorageDelegate) error
}

type AttributePersistence struct {
	mStorage store.PersistentStorageDelegate
}

func NewAttributePersistence() *AttributePersistence {

	return &AttributePersistence{}
}

func (p AttributePersistence) Init(storage store.PersistentStorageDelegate) (err error) {
	p.mStorage = storage
	return
}
