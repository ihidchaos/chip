package core

import (
	"github.com/galenliu/chip/credentials"
	"github.com/galenliu/chip/lib/store"
)

type AclStorage interface {
	Init(storage store.PersistentStorageDelegate, fabrics *credentials.FabricTable) error
}

type AclStorageImpl struct {
}

func NewAclStorageImpl() *AclStorageImpl {
	return &AclStorageImpl{}
}

func (d AclStorageImpl) Init(storage store.PersistentStorageDelegate, fabrics *credentials.FabricTable) error {
	return nil
}
