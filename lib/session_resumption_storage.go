package lib

import (
	"github.com/galenliu/chip/pkg/storage"
)

type SessionResumptionStorage interface {
	Init(delegate storage.KvsPersistentStorageDelegate) error
}

type SessionResumptionStorageImpl struct {
	mStorage storage.KvsPersistentStorageDelegate
}

func (s *SessionResumptionStorageImpl) Init(delegate storage.KvsPersistentStorageDelegate) error {
	s.mStorage = delegate
	return nil
}

func NewSimpleSessionResumptionStorage() *SessionResumptionStorageImpl {
	return &SessionResumptionStorageImpl{}
}
