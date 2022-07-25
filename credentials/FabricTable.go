package credentials

import (
	"github.com/galenliu/chip/crypto"
	"github.com/galenliu/chip/storage"
)

type FabricTable struct {
	mState []FabricInfo
}

func (f FabricTable) FabricCount() int {
	//TODO implement me
	panic("implement me")
}

func (f FabricTable) Init(params *InitParams) (err error) {
	return
}

func (f FabricTable) GetFabricInfos() []FabricInfo {
	return f.mState
}

func (f FabricTable) DeleteAllFabrics() {

}

func (f FabricTable) AddFabricDelegate(delegate ServerFabricDelegate) {

}

type InitParams struct {
	Storage             storage.PersistentStorageDelegate
	OperationalKeystore crypto.PersistentStorageOperationalKeystore
	OpCertStore         PersistentStorageOpCertStore
}

func NewFabricTableInitParams() *InitParams {
	return &InitParams{}
}
