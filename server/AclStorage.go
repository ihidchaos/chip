package server

import (
	"github.com/galenliu/chip/credentials"
	"github.com/galenliu/chip/storage"
)

type AclStorage interface {
	Init(storage storage.PersistentStorageDelegate, fabrics *credentials.FabricTable) error
}

type DefaultAclStorage struct {
}

func (d DefaultAclStorage) Init(storage storage.PersistentStorageDelegate, fabrics *credentials.FabricTable) error {
	//TODO implement me
	panic("implement me")
}
