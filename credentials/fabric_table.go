package credentials

import (
	"github.com/galenliu/chip/crypto"
	storage2 "github.com/galenliu/chip/crypto/persistent_storage"
	"github.com/galenliu/chip/lib"
	"github.com/galenliu/chip/storage"
	"time"
)

type FabricTableInitParams struct {
	Storage             storage.StorageDelegate
	OperationalKeystore storage2.PersistentStorageOperationalKeystore
	OpCertStore         PersistentStorageOpCertStore
}

type FabricTableDelegate interface {
	FabricWillBeRemoved(table FabricTable, index lib.FabricIndex)
	OnFabricRemoved(table FabricTable, index lib.FabricIndex)
	OnFabricCommitted(table FabricTable, index lib.FabricIndex)
	OnFabricUpdated(table FabricTable, index lib.FabricIndex)
}

type FabricTableContainer interface {
	Init(FabricTableInitParams) error
	Delete(index lib.FabricIndex)
	DeleteAllFabrics() error
	GetDeletedFabricFromCommitMarker() lib.FabricIndex
	ClearCommitMarker()
	Forget(index lib.FabricIndex)
	AddFabricDelegate(delegate FabricTableDelegate) error
	RemoveFabricDelegate(delegate FabricTableDelegate)
	SetFabricLabel(index lib.FabricIndex, label string) error
	GetFabricLabel(index lib.FabricIndex) (string, error)
	GetLastKnownGoodChipEpochTime() (time.Time, error)
	SetLastKnownGoodChipEpochTime(time.Time) error
	FabricCount() uint8

	FetchRootCert(lib.FabricIndex) ([]byte, error)
	FetchPendingNonFabricAssociatedRootCert() ([]byte, error)
	FetchICACert(index lib.FabricIndex) ([]byte, error)
	FetchNOCCert(index lib.FabricIndex) ([]byte, error)
	FetchRootPubkey(index lib.FabricIndex) ([]byte, error)
	FetchCATs(index lib.FabricIndex) ([]byte, error)

	SignWithOpKeypair(lib.FabricIndex) crypto.P256ECDSASignature
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
