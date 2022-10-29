package core

import (
	"github.com/galenliu/chip/credentials"
	"github.com/galenliu/chip/pkg/storage"
)

type AclStorage interface {
	Init(storage storage.KvsPersistentStorageDelegate, fabrics *credentials.FabricTable) error
}

type AclStorageImpl struct {
}

func NewAclStorageImpl() *AclStorageImpl {
	return &AclStorageImpl{}
}

func (d AclStorageImpl) Init(storage storage.KvsPersistentStorageDelegate, fabrics *credentials.FabricTable) error {
	return nil
}
