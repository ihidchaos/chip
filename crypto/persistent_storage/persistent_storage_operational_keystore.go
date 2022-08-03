package persistent_storage

import (
	"crypto"
	"github.com/galenliu/chip/device"
	"github.com/galenliu/chip/pkg/storage"
)

type PersistentStorageOperationalKeystore interface {
	HasPendingOpKeypair() bool
	HasOpKeypairForFabric(fabricIndex device.FabricIndex) bool
	NewOpKeypairForFabric(fabricIndex device.FabricIndex) ([]byte, error)
	ActivateOpKeypairForFabric(fabricIndex device.FabricIndex, key crypto.PublicKey) error
	CommitOpKeypairForFabric(fabricIndex device.FabricIndex) error
	RemoveOpKeypairForFabric(fabricIndex device.FabricIndex) error
	RevertPendingKeypair()
	SignWithOpKeypair(fabricIndex device.FabricIndex, message []byte) ([]byte, error)
	AllocateEphemeralKeypairForCASE() crypto.PrivateKey
	ReleaseEphemeralKeypair(key crypto.PrivateKey)
	Init(delegate storage.StorageDelegate)
}

type PersistentStorageOperationalKeystoreImpl struct {
	mPersistentStorage storage.StorageDelegate
}

func NewPersistentStorageOperationalKeystoreImpl() *PersistentStorageOperationalKeystoreImpl {
	return &PersistentStorageOperationalKeystoreImpl{}
}

func (p PersistentStorageOperationalKeystoreImpl) HasPendingOpKeypair() bool {
	//TODO implement me
	panic("implement me")
}

func (p PersistentStorageOperationalKeystoreImpl) HasOpKeypairForFabric(fabricIndex device.FabricIndex) bool {
	//TODO implement me
	panic("implement me")
}

func (p PersistentStorageOperationalKeystoreImpl) NewOpKeypairForFabric(fabricIndex device.FabricIndex) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (p PersistentStorageOperationalKeystoreImpl) ActivateOpKeypairForFabric(fabricIndex device.FabricIndex, key crypto.PublicKey) error {
	//TODO implement me
	panic("implement me")
}

func (p PersistentStorageOperationalKeystoreImpl) CommitOpKeypairForFabric(fabricIndex device.FabricIndex) error {
	//TODO implement me
	panic("implement me")
}

func (p PersistentStorageOperationalKeystoreImpl) RemoveOpKeypairForFabric(fabricIndex device.FabricIndex) error {
	//TODO implement me
	panic("implement me")
}

func (p PersistentStorageOperationalKeystoreImpl) RevertPendingKeypair() {
	//TODO implement me
	panic("implement me")
}

func (p PersistentStorageOperationalKeystoreImpl) SignWithOpKeypair(fabricIndex device.FabricIndex, message []byte) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (p PersistentStorageOperationalKeystoreImpl) AllocateEphemeralKeypairForCASE() crypto.PrivateKey {
	//TODO implement me
	panic("implement me")
}

func (p PersistentStorageOperationalKeystoreImpl) ReleaseEphemeralKeypair(key crypto.PrivateKey) {
	//TODO implement me
	panic("implement me")
}

func (p PersistentStorageOperationalKeystoreImpl) Init(delegate storage.StorageDelegate) {
	p.mPersistentStorage = delegate
}
