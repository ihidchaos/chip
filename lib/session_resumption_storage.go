package lib

import (
	"github.com/galenliu/chip/pkg/store"
)

type SessionResumptionStorage interface {
	Init(delegate store.KvsPersistentStorageBase) error
}

type SessionResumptionStorageImpl struct {
	mStorage store.KvsPersistentStorageBase
}

func (s *SessionResumptionStorageImpl) Init(delegate store.KvsPersistentStorageBase) error {
	s.mStorage = delegate
	return nil
}

func NewSimpleSessionResumptionStorage() *SessionResumptionStorageImpl {
	return &SessionResumptionStorageImpl{}
}
