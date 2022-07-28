package credentials

import (
	"github.com/galenliu/chip/storage"
)

type GroupDataProvider interface {
	SetStorageDelegate(delegate storage.StorageDelegate)
	Init() error
	SetListener(listener GroupDataProviderListener)
}

type GroupDataProviderImpl struct {
	mStorage storage.StorageDelegate
}

func NewGroupDataProviderImpl() *GroupDataProviderImpl {
	return &GroupDataProviderImpl{}
}

func (g GroupDataProviderImpl) SetListener(listener GroupDataProviderListener) {

}

func (g *GroupDataProviderImpl) SetStorageDelegate(delegate storage.StorageDelegate) {
	g.mStorage = delegate
}

func (g GroupDataProviderImpl) Init() error {
	return nil
}
