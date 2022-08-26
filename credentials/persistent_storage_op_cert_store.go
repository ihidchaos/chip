package credentials

import (
	"github.com/galenliu/chip/lib"
)

type PersistentStorageOpCertStore interface {
	OperationalCertificateStore
	HasPendingRootCert() bool
	HasPendingNocChain() bool
	HasCertificateForFabric(fabricIndex lib.FabricIndex, element uint8) bool

	AddNewTrustedRootCertForFabric(fabricIndex lib.FabricIndex, rcac []byte) error

	AddNewOpCertsForFabric(fabricIndex lib.FabricIndex, noc []byte, icac []byte) error

	UpdateOpCertsForFabric(fabricIndex lib.FabricIndex, noc []byte, icac []byte) error

	CommitOpCertsForFabric(fabricIndex lib.FabricIndex) error

	RemoveOpCertsForFabric(fabricIndex lib.FabricIndex) error

	RevertPendingOpCerts()

	RevertPendingOpCertsExceptRoot()

	GetCertificate(fabricIndex lib.FabricIndex, element uint8) ([]byte, error)
}

type PersistentStorageOpCertStoreImpl struct {
	*OperationalCertificateStoreImpl
}

func NewPersistentStorageOpCertStoreImpl() *PersistentStorageOpCertStoreImpl {
	return &PersistentStorageOpCertStoreImpl{
		NewOperationalCertificateStoreImpl(),
	}
}

func (s PersistentStorageOpCertStoreImpl) HasPendingRootCert() bool {
	return true
}

func (s PersistentStorageOpCertStoreImpl) HasPendingNocChain() bool {
	return true
}

func (s PersistentStorageOpCertStoreImpl) HasCertificateForFabric(fabricIndex lib.FabricIndex, element uint8) bool {
	return true
}

func (s PersistentStorageOpCertStoreImpl) AddNewTrustedRootCertForFabric(fabricIndex lib.FabricIndex, rcac []byte) error {
	return nil
}

func (s PersistentStorageOpCertStoreImpl) AddNewOpCertsForFabric(fabricIndex lib.FabricIndex, noc []byte, icac []byte) error {
	return nil
}

func (s PersistentStorageOpCertStoreImpl) UpdateOpCertsForFabric(fabricIndex lib.FabricIndex, noc []byte, icac []byte) error {
	//TODO implement me
	panic("implement me")
}

func (s PersistentStorageOpCertStoreImpl) CommitOpCertsForFabric(fabricIndex lib.FabricIndex) error {
	return nil
}

func (s PersistentStorageOpCertStoreImpl) RemoveOpCertsForFabric(fabricIndex lib.FabricIndex) error {
	return nil
}

func (s PersistentStorageOpCertStoreImpl) RevertPendingOpCerts() {

}

func (s PersistentStorageOpCertStoreImpl) RevertPendingOpCertsExceptRoot() {

}

func (s PersistentStorageOpCertStoreImpl) GetCertificate(fabricIndex lib.FabricIndex, element uint8) ([]byte, error) {
	return nil, nil
}
