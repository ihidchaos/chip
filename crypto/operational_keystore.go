package crypto

import (
	"crypto"
	"github.com/galenliu/chip/lib"
	"github.com/galenliu/chip/pkg/store"
)

type OperationalKeystore interface {
	Init(persistentStorage store.KvsPersistentStorageBase) error
	HasPendingOpKeypair() bool
	HasOpKeypairForFabric(fabricIndex lib.FabricIndex) bool
	NewOpKeypairForFabric(fabricIndex lib.FabricIndex) ([]byte, error)
	ActivateOpKeypairForFabric(fabricIndex lib.FabricIndex, key crypto.PublicKey) error
	CommitOpKeypairForFabric(fabricIndex lib.FabricIndex) error
	RemoveOpKeypairForFabric(fabricIndex lib.FabricIndex) error
	RevertPendingKeypair()
	SignWithOpKeypair(fabricIndex lib.FabricIndex, message []byte) ([]byte, error)
	AllocateEphemeralKeypairForCASE() *P256Keypair
	ReleaseEphemeralKeypair(key *P256Keypair)
}

func NewOperationalKeystoreImpl() *OperationalKeystoreImpl {
	return &OperationalKeystoreImpl{}
}

type OperationalKeystoreImpl struct {
	mPersistentStorage store.KvsPersistentStorageBase
}

func (p *OperationalKeystoreImpl) Init(persistentStorage store.KvsPersistentStorageBase) error {
	p.mPersistentStorage = persistentStorage
	return nil
}

func (p *OperationalKeystoreImpl) HasPendingOpKeypair() bool {
	//TODO implement me
	panic("implement me")
}

func (p *OperationalKeystoreImpl) HasOpKeypairForFabric(fabricIndex lib.FabricIndex) bool {
	//TODO implement me
	panic("implement me")
}

func (p *OperationalKeystoreImpl) NewOpKeypairForFabric(fabricIndex lib.FabricIndex) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (p *OperationalKeystoreImpl) ActivateOpKeypairForFabric(fabricIndex lib.FabricIndex, key crypto.PublicKey) error {
	//TODO implement me
	panic("implement me")
}

func (p *OperationalKeystoreImpl) CommitOpKeypairForFabric(fabricIndex lib.FabricIndex) error {
	//TODO implement me
	panic("implement me")
}

func (p *OperationalKeystoreImpl) RemoveOpKeypairForFabric(fabricIndex lib.FabricIndex) error {
	//TODO implement me
	panic("implement me")
}

func (p *OperationalKeystoreImpl) RevertPendingKeypair() {
	//TODO implement me
	panic("implement me")
}

func (p *OperationalKeystoreImpl) SignWithOpKeypair(fabricIndex lib.FabricIndex, message []byte) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (p *OperationalKeystoreImpl) AllocateEphemeralKeypairForCASE() *P256Keypair {
	return GenericP256Keypair()
}

func (p *OperationalKeystoreImpl) ReleaseEphemeralKeypair(key *P256Keypair) {
	//TODO implement me
	panic("implement me")
}
