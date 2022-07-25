package transport

import (
	"github.com/galenliu/chip/credentials"
	"github.com/galenliu/chip/storage"
)

type SessionManager interface {
	Init(transports TransportManager, storage storage.PersistentStorageDelegate, table *credentials.FabricTable) error
}

type SessionManagerImpl struct {
}
