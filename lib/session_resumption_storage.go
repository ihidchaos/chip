package lib

import (
	"github.com/galenliu/chip/pkg/storage"
)

type SessionResumptionStorage interface {
	Init(delegate storage.KeyValuePersistentStorage) error
}

type SessionResumptionStorageImpl struct {
	mStorage storage.KeyValuePersistentStorage
}

func (s *SessionResumptionStorageImpl) Init(delegate storage.KeyValuePersistentStorage) error {
	s.mStorage = delegate
	return nil
}

func NewSimpleSessionResumptionStorage() *SessionResumptionStorageImpl {
	return &SessionResumptionStorageImpl{}
}
