package credentials

import "github.com/galenliu/chip/pkg/storage"

const (
	CertChainElement_Rcac uint8 = 0
	CertChainElement_Icac uint8 = 1
	CertChainElement_Noc  uint8 = 2
)

type OperationalCertificateStore interface {
	Init(persistentStorage storage.KvsPersistentStorageDelegate) error
}

type OperationalCertificateStoreImpl struct {
	mStorage storage.KvsPersistentStorageDelegate
}

func NewOperationalCertificateStoreImpl() *OperationalCertificateStoreImpl {
	return &OperationalCertificateStoreImpl{}
}

func (o OperationalCertificateStoreImpl) Init(persistentStorage storage.KvsPersistentStorageDelegate) error {
	o.mStorage = persistentStorage
	return nil
}
