package crypto

import (
	"crypto"
	"github.com/galenliu/chip/core"
	"github.com/galenliu/chip/storage"
)

type PersistentStorageOperationalKeystore interface {
	HasPendingOpKeypair() bool
	HasOpKeypairForFabric(fabricIndex core.FabricIndex) bool
	NewOpKeypairForFabric(fabricIndex core.FabricIndex) ([]byte, error)
	ActivateOpKeypairForFabric(fabricIndex core.FabricIndex, key crypto.PublicKey) error
	CommitOpKeypairForFabric(fabricIndex core.FabricIndex) error
	RemoveOpKeypairForFabric(fabricIndex core.FabricIndex) error
	RevertPendingKeypair()
	SignWithOpKeypair(fabricIndex core.FabricIndex, message []byte) ([]byte, error)
	AllocateEphemeralKeypairForCASE() crypto.PrivateKey
	ReleaseEphemeralKeypair(key crypto.PrivateKey)
	Init(delegate storage.PersistentStorageDelegate)
}

type PersistentStorageOperationalKeystoreImpl struct {
	mPersistentStorage storage.PersistentStorageDelegate
}

func NewPersistentStorageOperationalKeystoreImpl() *PersistentStorageOperationalKeystoreImpl {
	return &PersistentStorageOperationalKeystoreImpl{}
}

func (p PersistentStorageOperationalKeystoreImpl) HasPendingOpKeypair() bool {
	//TODO implement me
	panic("implement me")
}

func (p PersistentStorageOperationalKeystoreImpl) HasOpKeypairForFabric(fabricIndex core.FabricIndex) bool {
	//TODO implement me
	panic("implement me")
}

func (p PersistentStorageOperationalKeystoreImpl) NewOpKeypairForFabric(fabricIndex core.FabricIndex) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (p PersistentStorageOperationalKeystoreImpl) ActivateOpKeypairForFabric(fabricIndex core.FabricIndex, key crypto.PublicKey) error {
	//TODO implement me
	panic("implement me")
}

func (p PersistentStorageOperationalKeystoreImpl) CommitOpKeypairForFabric(fabricIndex core.FabricIndex) error {
	//TODO implement me
	panic("implement me")
}

func (p PersistentStorageOperationalKeystoreImpl) RemoveOpKeypairForFabric(fabricIndex core.FabricIndex) error {
	//TODO implement me
	panic("implement me")
}

func (p PersistentStorageOperationalKeystoreImpl) RevertPendingKeypair() {
	//TODO implement me
	panic("implement me")
}

func (p PersistentStorageOperationalKeystoreImpl) SignWithOpKeypair(fabricIndex core.FabricIndex, message []byte) ([]byte, error) {
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

func (p PersistentStorageOperationalKeystoreImpl) Init(delegate storage.PersistentStorageDelegate) {
	p.mPersistentStorage = delegate
}
