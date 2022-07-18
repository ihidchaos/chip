package credentials

import (
	"github.com/galenliu/chip/storage"
)

type GroupDataProvider interface {
	SetStorageDelegate(delegate storage.PersistentStorageDelegate)
	Init() error
}

type GroupDataProviderImpl struct {
}

func (g GroupDataProviderImpl) SetStorageDelegate(delegate storage.PersistentStorageDelegate) {
	//TODO implement me
	panic("implement me")
}

func (g GroupDataProviderImpl) Init() error {
	//TODO implement me
	panic("implement me")
}
