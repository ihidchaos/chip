package credentials

import (
	"github.com/galenliu/chip/crypto"
	storage2 "github.com/galenliu/chip/crypto/persistent_storage"
	"github.com/galenliu/chip/storage"
	"time"
)

type FabricTableInitParams struct {
	Storage             storage.StorageDelegate
	OperationalKeystore storage2.PersistentStorageOperationalKeystore
	OpCertStore         PersistentStorageOpCertStore
}

type FabricTableDelegate interface {
}

type FabricTableProvider interface {
	Init(FabricTableInitParams) error
	Delete(index FabricIndex)
	DeleteAllFabrics() error
	GetDeletedFabricFromCommitMarker() FabricIndex
	ClearCommitMarker()
	Forget(index FabricIndex)
	AddFabricDelegate(delegate FabricTableDelegate) error
	RemoveFabricDelegate(delegate FabricTableDelegate)
	SetFabricLabel(index FabricIndex, label string) error
	GetFabricLabel(index FabricIndex) (string, error)
	GetLastKnownGoodChipEpochTime() (time.Time, error)
	SetLastKnownGoodChipEpochTime(time.Time) error
	FabricCount() uint8

	FetchRootCert(FabricIndex) ([]byte, error)
	FetchPendingNonFabricAssociatedRootCert() ([]byte, error)
	FetchICACert(index FabricIndex) ([]byte, error)
	FetchNOCCert(index FabricIndex) ([]byte, error)
	FetchRootPubkey(index FabricIndex) ([]byte, error)
	FetchCATs(index FabricIndex) ([]byte, error)

	SignWithOpKeypair(FabricIndex) crypto.P256ECDSASignature
}

type FabricTable struct {
	mState []FabricInfo
}

func NewFabricTable() *FabricTable {
	return &FabricTable{}
}

func (f FabricTable) FabricCount() int {
	return len(f.mState)
}

func (f FabricTable) Init(params *FabricTableInitParams) (err error) {
	return
}

func (f FabricTable) GetFabricInfos() []FabricInfo {
	return f.mState
}

func (f FabricTable) DeleteAllFabrics() {

}

func (f FabricTable) AddFabricDelegate(delegate ServerFabricDelegate) {

}

func NewFabricTableInitParams() *FabricTableInitParams {
	return &FabricTableInitParams{}
}
