package operational_storage

import (
	"crypto"
	"github.com/galenliu/chip/lib"
	"github.com/galenliu/chip/pkg/storage"
)

type OperationalKeystore interface {
	Init(persistentStorage storage.KeyValuePersistentStorage) error
	HasPendingOpKeypair() bool
	HasOpKeypairForFabric(fabricIndex lib.FabricIndex) bool
	NewOpKeypairForFabric(fabricIndex lib.FabricIndex) ([]byte, error)
	ActivateOpKeypairForFabric(fabricIndex lib.FabricIndex, key crypto.PublicKey) error
	CommitOpKeypairForFabric(fabricIndex lib.FabricIndex) error
	RemoveOpKeypairForFabric(fabricIndex lib.FabricIndex) error
	RevertPendingKeypair()
	SignWithOpKeypair(fabricIndex lib.FabricIndex, message []byte) ([]byte, error)
	AllocateEphemeralKeypairForCASE() crypto.PrivateKey
	ReleaseEphemeralKeypair(key crypto.PrivateKey)
}

func NewOperationalKeystoreImpl() *OperationalKeystoreImpl {
	return &OperationalKeystoreImpl{}
}

type OperationalKeystoreImpl struct {
	mPersistentStorage storage.PersistentStorage
}

func (p *OperationalKeystoreImpl) Init(persistentStorage storage.KeyValuePersistentStorage) error {
	p.mPersistentStorage = persistentStorage
	return nil
}

func (p OperationalKeystoreImpl) HasPendingOpKeypair() bool {
	//TODO implement me
	panic("implement me")
}

func (p OperationalKeystoreImpl) HasOpKeypairForFabric(fabricIndex lib.FabricIndex) bool {
	//TODO implement me
	panic("implement me")
}

func (p OperationalKeystoreImpl) NewOpKeypairForFabric(fabricIndex lib.FabricIndex) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (p OperationalKeystoreImpl) ActivateOpKeypairForFabric(fabricIndex lib.FabricIndex, key crypto.PublicKey) error {
	//TODO implement me
	panic("implement me")
}

func (p OperationalKeystoreImpl) CommitOpKeypairForFabric(fabricIndex lib.FabricIndex) error {
	//TODO implement me
	panic("implement me")
}

func (p OperationalKeystoreImpl) RemoveOpKeypairForFabric(fabricIndex lib.FabricIndex) error {
	//TODO implement me
	panic("implement me")
}

func (p OperationalKeystoreImpl) RevertPendingKeypair() {
	//TODO implement me
	panic("implement me")
}

func (p OperationalKeystoreImpl) SignWithOpKeypair(fabricIndex lib.FabricIndex, message []byte) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (p OperationalKeystoreImpl) AllocateEphemeralKeypairForCASE() crypto.PrivateKey {
	//TODO implement me
	panic("implement me")
}

func (p OperationalKeystoreImpl) ReleaseEphemeralKeypair(key crypto.PrivateKey) {
	//TODO implement me
	panic("implement me")
}
