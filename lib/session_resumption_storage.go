package lib

import (
	"github.com/galenliu/chip/lib/store"
)

type SessionResumptionStorage interface {
	Init(delegate store.PersistentStorageDelegate) error
}

type SessionResumptionStorageImpl struct {
	mStorage store.PersistentStorageDelegate
}

func (s *SessionResumptionStorageImpl) Init(delegate store.PersistentStorageDelegate) error {
	s.mStorage = delegate
	return nil
}

func NewSimpleSessionResumptionStorage() *SessionResumptionStorageImpl {
	return &SessionResumptionStorageImpl{}
}
