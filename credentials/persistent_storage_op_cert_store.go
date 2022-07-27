package credentials

import (
	"github.com/galenliu/chip/core"
	"github.com/galenliu/chip/storage"
)

type PersistentStorageOpCertStore interface {
	Init(delegate storage.PersistentStorageDelegate)

	HasPendingRootCert() bool
	HasPendingNocChain() bool
	HasCertificateForFabric(fabricIndex core.FabricIndex, element uint8) bool

	AddNewTrustedRootCertForFabric(fabricIndex core.FabricIndex, rcac []byte) error

	AddNewOpCertsForFabric(fabricIndex core.FabricIndex, noc []byte, icac []byte) error

	UpdateOpCertsForFabric(fabricIndex core.FabricIndex, noc []byte, icac []byte) error

	CommitOpCertsForFabric(fabricIndex core.FabricIndex) error

	RemoveOpCertsForFabric(fabricIndex core.FabricIndex) error

	RevertPendingOpCerts()

	RevertPendingOpCertsExceptRoot()

	GetCertificate(fabricIndex core.FabricIndex, element uint8) []byte
}

type PersistentStorageOpCertStoreImpl struct {
	mPersistentStorage storage.PersistentStorageDelegate
}

func NewPersistentStorageOpCertStoreImpl() *PersistentStorageOpCertStoreImpl {
	return &PersistentStorageOpCertStoreImpl{}
}

func (s PersistentStorageOpCertStoreImpl) HasPendingRootCert() bool {
	return true
}

func (s PersistentStorageOpCertStoreImpl) HasPendingNocChain() bool {
	return true
}

func (s PersistentStorageOpCertStoreImpl) HasCertificateForFabric(fabricIndex core.FabricIndex, element uint8) bool {
	return true
}

func (s PersistentStorageOpCertStoreImpl) AddNewTrustedRootCertForFabric(fabricIndex core.FabricIndex, rcac []byte) error {
	return nil
}

func (s PersistentStorageOpCertStoreImpl) AddNewOpCertsForFabric(fabricIndex core.FabricIndex, noc []byte, icac []byte) error {
	return nil
}

func (s PersistentStorageOpCertStoreImpl) UpdateOpCertsForFabric(fabricIndex core.FabricIndex, noc []byte, icac []byte) error {
	//TODO implement me
	panic("implement me")
}

func (s PersistentStorageOpCertStoreImpl) CommitOpCertsForFabric(fabricIndex core.FabricIndex) error {
	return nil
}

func (s PersistentStorageOpCertStoreImpl) RemoveOpCertsForFabric(fabricIndex core.FabricIndex) error {
	return nil
}

func (s PersistentStorageOpCertStoreImpl) RevertPendingOpCerts() {

}

func (s PersistentStorageOpCertStoreImpl) RevertPendingOpCertsExceptRoot() {

}

func (s PersistentStorageOpCertStoreImpl) GetCertificate(fabricIndex core.FabricIndex, element uint8) []byte {
	return nil
}

func (s PersistentStorageOpCertStoreImpl) Init(delegate storage.PersistentStorageDelegate) {
	s.mPersistentStorage = delegate
}
