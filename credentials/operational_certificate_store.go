package credentials

import "github.com/galenliu/chip/pkg/storage"

type OperationalCertificateStore interface {
	Init(persistentStorage storage.PersistentStorage) error
}

type OperationalCertificateStoreImpl struct {
	mStorage storage.PersistentStorage
}

func (o OperationalCertificateStoreImpl) Init(persistentStorage storage.PersistentStorage) error {
	o.mStorage = persistentStorage
	return nil
}
