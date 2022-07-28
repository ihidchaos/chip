package transport

import (
	"github.com/galenliu/chip/credentials"
	"github.com/galenliu/chip/storage"
)

type SessionManager interface {
	Init(transports TransportManager, storage storage.StorageDelegate, table *credentials.FabricTable) error
}

type SessionManagerImpl struct {
}

func (s SessionManagerImpl) Init(transports TransportManager, storage storage.StorageDelegate, table *credentials.FabricTable) error {
	return nil
}

func NewSessionManagerImpl() *SessionManagerImpl {
	return &SessionManagerImpl{}
}
