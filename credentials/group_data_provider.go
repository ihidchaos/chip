package credentials

import (
	"github.com/galenliu/chip/lib"
	"github.com/galenliu/chip/storage"
)

type GroupInfo struct {
	Id lib.GroupId
}

type GroupDataProvider interface {
	SetStorageDelegate(delegate storage.StorageDelegate)
	Init() error
	SetListener(listener GroupListener)
}

type GroupDataProviderImpl struct {
	mStorage       storage.StorageDelegate
	mGroupListener GroupListener
}

func NewGroupDataProviderImpl() *GroupDataProviderImpl {
	return &GroupDataProviderImpl{}
}

func (g *GroupDataProviderImpl) SetListener(listener GroupListener) {
	g.mGroupListener = listener
}

func (g *GroupDataProviderImpl) SetStorageDelegate(delegate storage.StorageDelegate) {
	g.mStorage = delegate
}

func (g *GroupDataProviderImpl) Init() error {
	return nil
}

var gGroupDataProvider GroupDataProvider

func GetGroupDataProvider() GroupDataProvider {
	return gGroupDataProvider
}

func SetGroupDataProvider(g GroupDataProvider) {
	gGroupDataProvider = g
}
