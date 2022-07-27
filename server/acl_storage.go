package server

import (
	"github.com/galenliu/chip/credentials"
	"github.com/galenliu/chip/storage"
)

type AclStorage interface {
	Init(storage storage.PersistentStorageDelegate, fabrics *credentials.FabricTable) error
}

type AclStorageImpl struct {
}

func NewAclStorageImpl() *AclStorageImpl {
	return &AclStorageImpl{}
}

func (d AclStorageImpl) Init(storage storage.PersistentStorageDelegate, fabrics *credentials.FabricTable) error {
	return nil
}
