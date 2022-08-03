package lib

import (
	"github.com/galenliu/chip/pkg/storage"
)

type AttributePersistenceProvider interface {
	Init(storage storage.StorageDelegate) error
}

type AttributePersistence struct {
	mStorage storage.StorageDelegate
}

func NewAttributePersistence() *AttributePersistence {
	return &AttributePersistence{}
}

func (p AttributePersistence) Init(storage storage.StorageDelegate) (err error) {
	p.mStorage = storage
	return
}
