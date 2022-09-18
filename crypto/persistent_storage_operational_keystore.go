package crypto

import (
	"crypto"
	"github.com/galenliu/chip/lib"
)

type PersistentStorageOperationalKeystore interface {
	OperationalKeystore
}

type PersistentStorageOperationalKeystoreImpl struct {
	*OperationalKeystoreImpl
}

func NewPersistentStorageOperationalKeystoreImpl() *PersistentStorageOperationalKeystoreImpl {
	return &PersistentStorageOperationalKeystoreImpl{
		NewOperationalKeystoreImpl(),
	}
}

func (p *PersistentStorageOperationalKeystoreImpl) HasPendingOpKeypair() bool {
	//TODO implement me
	panic("implement me")
}

func (p *PersistentStorageOperationalKeystoreImpl) HasOpKeypairForFabric(fabricIndex lib.FabricIndex) bool {
	//TODO implement me
	panic("implement me")
}

func (p *PersistentStorageOperationalKeystoreImpl) NewOpKeypairForFabric(fabricIndex lib.FabricIndex) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (p *PersistentStorageOperationalKeystoreImpl) ActivateOpKeypairForFabric(fabricIndex lib.FabricIndex, key crypto.PublicKey) error {
	//TODO implement me
	panic("implement me")
}

func (p *PersistentStorageOperationalKeystoreImpl) CommitOpKeypairForFabric(fabricIndex lib.FabricIndex) error {
	//TODO implement me
	panic("implement me")
}

func (p *PersistentStorageOperationalKeystoreImpl) RemoveOpKeypairForFabric(fabricIndex lib.FabricIndex) error {
	//TODO implement me
	panic("implement me")
}

func (p *PersistentStorageOperationalKeystoreImpl) RevertPendingKeypair() {
	//TODO implement me
	panic("implement me")
}

func (p *PersistentStorageOperationalKeystoreImpl) SignWithOpKeypair(fabricIndex lib.FabricIndex, message []byte) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (p *PersistentStorageOperationalKeystoreImpl) AllocateEphemeralKeypairForCASE() *P256Keypair {
	//TODO implement me
	panic("implement me")
}

func (p *PersistentStorageOperationalKeystoreImpl) ReleaseEphemeralKeypair(key *P256Keypair) {
	//TODO implement me
	panic("implement me")
}
