package lib

import "github.com/galenliu/chip/storage"

type SessionResumptionStorage interface {
	Init(delegate storage.PersistentStorageDelegate) error
}

type SessionResumptionStorageImpl struct {
	mStorage storage.PersistentStorageDelegate
}

func (s *SessionResumptionStorageImpl) Init(delegate storage.PersistentStorageDelegate) error {
	s.mStorage = delegate
	return nil
}

func NewSimpleSessionResumptionStorage() *SessionResumptionStorageImpl {
	return &SessionResumptionStorageImpl{}
}
