package credentials

import (
	"github.com/galenliu/chip/lib/store"
)

const (
	CertChainElement_Rcac uint8 = 0
	CertChainElement_Icac uint8 = 1
	CertChainElement_Noc  uint8 = 2
)

type OperationalCertificateStore interface {
	Init(persistentStorage store.PersistentStorageDelegate) error
}

type OperationalCertificateStoreImpl struct {
	mStorage store.PersistentStorageDelegate
}

func NewOperationalCertificateStoreImpl() *OperationalCertificateStoreImpl {
	return &OperationalCertificateStoreImpl{}
}

func (o OperationalCertificateStoreImpl) Init(persistentStorage store.PersistentStorageDelegate) error {
	o.mStorage = persistentStorage
	return nil
}
