package credentials

import (
	"github.com/galenliu/chip/device"
	"github.com/galenliu/chip/pkg/storage"
)

type PersistentStorageOpCertStore interface {
	Init(delegate storage.StorageDelegate)

	HasPendingRootCert() bool
	HasPendingNocChain() bool
	HasCertificateForFabric(fabricIndex device.FabricIndex, element uint8) bool

	AddNewTrustedRootCertForFabric(fabricIndex device.FabricIndex, rcac []byte) error

	AddNewOpCertsForFabric(fabricIndex device.FabricIndex, noc []byte, icac []byte) error

	UpdateOpCertsForFabric(fabricIndex device.FabricIndex, noc []byte, icac []byte) error

	CommitOpCertsForFabric(fabricIndex device.FabricIndex) error

	RemoveOpCertsForFabric(fabricIndex device.FabricIndex) error

	RevertPendingOpCerts()

	RevertPendingOpCertsExceptRoot()

	GetCertificate(fabricIndex device.FabricIndex, element uint8) []byte
}

type PersistentStorageOpCertStoreImpl struct {
	mPersistentStorage storage.StorageDelegate
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

func (s PersistentStorageOpCertStoreImpl) HasCertificateForFabric(fabricIndex device.FabricIndex, element uint8) bool {
	return true
}

func (s PersistentStorageOpCertStoreImpl) AddNewTrustedRootCertForFabric(fabricIndex device.FabricIndex, rcac []byte) error {
	return nil
}

func (s PersistentStorageOpCertStoreImpl) AddNewOpCertsForFabric(fabricIndex device.FabricIndex, noc []byte, icac []byte) error {
	return nil
}

func (s PersistentStorageOpCertStoreImpl) UpdateOpCertsForFabric(fabricIndex device.FabricIndex, noc []byte, icac []byte) error {
	//TODO implement me
	panic("implement me")
}

func (s PersistentStorageOpCertStoreImpl) CommitOpCertsForFabric(fabricIndex device.FabricIndex) error {
	return nil
}

func (s PersistentStorageOpCertStoreImpl) RemoveOpCertsForFabric(fabricIndex device.FabricIndex) error {
	return nil
}

func (s PersistentStorageOpCertStoreImpl) RevertPendingOpCerts() {

}

func (s PersistentStorageOpCertStoreImpl) RevertPendingOpCertsExceptRoot() {

}

func (s PersistentStorageOpCertStoreImpl) GetCertificate(fabricIndex device.FabricIndex, element uint8) []byte {
	return nil
}

func (s PersistentStorageOpCertStoreImpl) Init(delegate storage.StorageDelegate) {
	s.mPersistentStorage = delegate
}
