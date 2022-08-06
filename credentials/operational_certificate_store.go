package credentials

import "github.com/galenliu/chip/pkg/storage"

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
