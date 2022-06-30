package server

import (
	"github.com/galenliu/chip/credentials"
	"github.com/galenliu/chip/lib"
)

type AclStorage struct {
}

func (s AclStorage) Init(storage lib.PersistentStorageDelegate, fabrics *credentials.FabricTable) error {
	return nil
}
