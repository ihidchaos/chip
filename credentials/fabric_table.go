package credentials

import (
	"github.com/galenliu/chip/crypto"
	storage2 "github.com/galenliu/chip/crypto/persistent_storage"
	"github.com/galenliu/chip/lib"
	"github.com/galenliu/chip/pkg/storage"
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
	Init(*FabricTableInitParams) error
	Delete(index lib.FabricIndex)
	DeleteAllFabrics()
	GetDeletedFabricFromCommitMarker() lib.FabricIndex
	ClearCommitMarker()
	Forget(index lib.FabricIndex)
	AddFabricDelegate(delegate FabricTableDelegate) error
	RemoveFabricDelegate(delegate FabricTableDelegate)
	SetFabricLabel(label string) error
	GetFabricLabel(index lib.FabricIndex) (string, error)
	GetLastKnownGoodChipEpochTime() (time.Time, error)
	SetLastKnownGoodChipEpochTime(time.Time) error
	FabricCount() uint8

	HasPendingFabricUpdate() bool

	FetchRootCert(lib.FabricIndex) ([]byte, error)
	FetchPendingNonFabricAssociatedRootCert() ([]byte, error)
	FetchICACert(index lib.FabricIndex) ([]byte, error)
	FetchNOCCert(index lib.FabricIndex) ([]byte, error)
	FetchRootPubkey(index lib.FabricIndex) ([]byte, error)
	FetchCATs(index lib.FabricIndex) ([]byte, error)
	SignWithOpKeypair(lib.FabricIndex) crypto.P256ECDSASignature
	FindFabricWithIndex(index lib.FabricIndex) *FabricInfo
}

type FabricTable struct {
	mStates        []*FabricInfo
	mPendingFabric *FabricInfo
	mFabricLabel   string
}

func (f *FabricTable) AddFabricDelegate(delegate FabricTableDelegate) error {
	//TODO implement me
	panic("implement me")
}

func (f *FabricTable) HasPendingFabricUpdate() bool {
	//TODO implement me
	panic("implement me")
}

func NewFabricTable() *FabricTable {
	return &FabricTable{}
}

func (f *FabricTable) Init(params *FabricTableInitParams) error {
	//TODO implement me
	panic("implement me")
}

func (f *FabricTable) Delete(index lib.FabricIndex) {
	//TODO implement me
	panic("implement me")
}

func (f *FabricTable) DeleteAllFabrics() {
	//TODO implement me
	panic("implement me")
}

func (f *FabricTable) GetDeletedFabricFromCommitMarker() lib.FabricIndex {
	//TODO implement me
	panic("implement me")
}

func (f *FabricTable) ClearCommitMarker() {
	//TODO implement me
	panic("implement me")
}

func (f *FabricTable) Forget(index lib.FabricIndex) {
	//TODO implement me
	panic("implement me")
}

func (f *FabricTable) RemoveFabricDelegate(delegate FabricTableDelegate) {
	//TODO implement me
	panic("implement me")
}

func (f *FabricTable) SetFabricLabel(label string) error {
	f.mFabricLabel = label
	return nil
}

func (f *FabricTable) GetFabricLabel(index lib.FabricIndex) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (f *FabricTable) GetLastKnownGoodChipEpochTime() (time.Time, error) {
	//TODO implement me
	panic("implement me")
}

func (f FabricTable) SetLastKnownGoodChipEpochTime(t time.Time) error {
	//TODO implement me
	panic("implement me")
}

func (f *FabricTable) FabricCount() uint8 {
	//TODO implement me
	panic("implement me")
}

func (f *FabricTable) FetchRootCert(index lib.FabricIndex) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (f *FabricTable) FetchPendingNonFabricAssociatedRootCert() ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (f *FabricTable) FetchICACert(index lib.FabricIndex) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (f *FabricTable) FetchNOCCert(index lib.FabricIndex) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (f *FabricTable) FetchRootPubkey(index lib.FabricIndex) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (f *FabricTable) FetchCATs(index lib.FabricIndex) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (f *FabricTable) SignWithOpKeypair(index lib.FabricIndex) crypto.P256ECDSASignature {
	//TODO implement me
	panic("implement me")
}

func (f *FabricTable) FindFabricWithIndex(index lib.FabricIndex) *FabricInfo {
	if f.HasPendingFabricUpdate() && f.mPendingFabric.GetFabricIndex() == index {
		return f.mPendingFabric
	}
	for _, f := range f.mStates {
		if !f.IsInitialized() {
			continue
		}
		if f.GetFabricIndex() == index {
			return f
		}
	}
	return nil
}

func NewFabricTableInitParams() *FabricTableInitParams {
	return &FabricTableInitParams{}
}
