package credentials

import (
	"github.com/galenliu/chip/storage"
)

type FabricTable struct {
	mState []FabricInfo
}

func (f FabricTable) FabricCount() int {
	//TODO implement me
	panic("implement me")
}

func (f FabricTable) Init(storage storage.PersistentStorageDelegate) (err error) {
	return
}

func (f FabricTable) GetFabricInfos() []FabricInfo {
	return f.mState
}

func (f FabricTable) DeleteAllFabrics() {

}
