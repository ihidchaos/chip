package credentials

import (
	"github.com/galenliu/chip/storage"
)

type GroupDataProvider interface {
	SetStorageDelegate(delegate storage.PersistentStorageDelegate)
	Init() error
	SetListener(listener GroupDataProviderListener)
}

type GroupDataProviderImpl struct {
	mStorage storage.PersistentStorageDelegate
}

func NewGroupDataProviderImpl() *GroupDataProviderImpl {
	return &GroupDataProviderImpl{}
}

func (g GroupDataProviderImpl) SetListener(listener GroupDataProviderListener) {

}

func (g *GroupDataProviderImpl) SetStorageDelegate(delegate storage.PersistentStorageDelegate) {
	g.mStorage = delegate
}

func (g GroupDataProviderImpl) Init() error {
	return nil
}
