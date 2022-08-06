package lib

import (
	"github.com/galenliu/chip/pkg/storage"
)

type AttributePersistenceProvider interface {
	Init(storage storage.PersistentStorage) error
}

type AttributePersistence struct {
	mStorage storage.PersistentStorage
}

func NewAttributePersistence() *AttributePersistence {

	return &AttributePersistence{}
}

func (p AttributePersistence) Init(storage storage.PersistentStorage) (err error) {
	p.mStorage = storage
	return
}
