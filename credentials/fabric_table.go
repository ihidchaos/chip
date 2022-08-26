package credentials

import (
	"github.com/galenliu/chip/crypto"
	"github.com/galenliu/chip/lib"
	"github.com/galenliu/chip/pkg/storage"
	"time"
)

type FabricTableInitParams struct {
	Storage             storage.KvsPersistentStorageDelegate
	OperationalKeystore crypto.OperationalKeystore
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
	Delete(index lib.FabricIndex) error
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
	SignWithOpKeypair(lib.FabricIndex) *crypto.P256ECDSASignature
	FindFabricWithIndex(index lib.FabricIndex) *FabricInfo
}

type FabricTable struct {
	mStates                   []*FabricInfo
	mPendingFabric            *FabricInfo
	mFabricLabel              string
	mStorage                  storage.KvsPersistentStorageDelegate
	operationalKeystore       crypto.OperationalKeystore
	mOpCertStore              PersistentStorageOpCertStore
	mFabricCount              uint8
	mNextAvailableFabricIndex lib.FabricIndex
	mDelegate                 FabricTableDelegate
}

func NewFabricTableInitParams() *FabricTableInitParams {
	return &FabricTableInitParams{}
}

func (f *FabricTable) AddFabricDelegate(delegate FabricTableDelegate) error {
	f.mDelegate = delegate
	return nil
}

func (f *FabricTable) HasPendingFabricUpdate() bool {
	//TODO implement me
	panic("implement me")
}

func NewFabricTable() *FabricTable {
	return &FabricTable{}
}

func (f *FabricTable) Init(params *FabricTableInitParams) error {
	f.mStorage = params.Storage
	f.mOpCertStore = params.OpCertStore
	f.mFabricCount = 0
	for _, f := range f.mStates {
		f.Reset()
	}
	f.mNextAvailableFabricIndex = lib.KMinValidFabricIndex
	return nil
}

func (f *FabricTable) Delete(index lib.FabricIndex) error {
	if f.mStorage == nil || !index.IsValidFabricIndex() {
		return lib.ChipErrorInvalidArgument
	}
	return nil
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
	fabricInfo := f.GetMutableFabricByIndex(index)
	if fabricInfo == nil {
		return
	}
	f.RevertPendingFabricData()
	fabricInfo.Reset()

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
	fabricInfo := f.FindFabricWithIndex(index)
	if fabricInfo == nil {
		return "", lib.ChipErrorInvalidFabricIndex
	}
	return fabricInfo.GetFabricLabel(), nil
}

func (f *FabricTable) GetLastKnownGoodChipEpochTime() (time.Time, error) {
	//TODO implement me
	panic("implement me")
}

func (f *FabricTable) SetLastKnownGoodChipEpochTime(t time.Time) error {
	//TODO implement me
	panic("implement me")
}

func (f *FabricTable) FabricCount() uint8 {
	return uint8(len(f.mStates))
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
	if f.mOpCertStore == nil {
		return nil, lib.ChipErrorIncorrectState
	}
	icaCert, err := f.mOpCertStore.GetCertificate(index, CertChainElement_Icac)
	if err != nil {
		if f.mOpCertStore.HasCertificateForFabric(index, CertChainElement_Noc) {
			return icaCert, nil
		}
	}

	return icaCert, err
}

func (f *FabricTable) FetchNOCCert(index lib.FabricIndex) ([]byte, error) {
	if f.mStorage == nil {
		return nil, lib.ChipErrorIncorrectState
	}
	return f.mOpCertStore.GetCertificate(index, CertChainElement_Noc)
}

func (f *FabricTable) FetchRootPubkey(index lib.FabricIndex) (*crypto.P256PublicKey, error) {
	fabricInfo := f.FindFabricWithIndex(index)
	if fabricInfo == nil {
		return nil, lib.ChipErrorInvalidFabricIndex
	}
	return fabricInfo.FetchRootPubkey()

}

func (f *FabricTable) FetchCATs(index lib.FabricIndex) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (f *FabricTable) SignWithOpKeypair(index lib.FabricIndex) *crypto.P256ECDSASignature {
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

func (f *FabricTable) GetFabrics() []*FabricInfo {
	return f.mStates
}

func (f *FabricTable) GetMutableFabricByIndex(index lib.FabricIndex) *FabricInfo {
	if f.HasPendingFabricUpdate() && f.mPendingFabric.GetFabricIndex() == index {
		return f.mPendingFabric
	}
	for _, fabricInfo := range f.mStates {
		if !fabricInfo.IsInitialized() {
			continue
		}
		if fabricInfo.GetFabricIndex() == index {
			return fabricInfo
		}
	}
	return nil
}

func (f *FabricTable) RevertPendingFabricData() {
	//TODO implement me
	panic("implement me")
}

func (f *FabricTable) AllocateEphemeralKeypairForCASE() *crypto.P256Keypair {
	if f.operationalKeystore != nil {
		return f.operationalKeystore.AllocateEphemeralKeypairForCASE()
	}
	return crypto.GenericP256Keypair()
}
