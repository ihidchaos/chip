package storage

import (
	"fmt"
	"strconv"
	"sync"
)

type PersistentStorage interface {
	Init(file string) error
	ReadBoolValue(key string) (bool, error)
	ReadValueUint16(key string) (uint16, error)
	ReadValueUint32(key string) (uint32, error)
	ReadValueUint64(key string) (uint64, error)
	ReadValueStr(key string) (string, error)
	ReadValueBin(key string) ([]byte, error)
	WriteValueBool(key string, v bool) error
	WriteValueUint16(key string, v uint16) error
	WriteValueUint32(key string, v uint32) error
	WriteValueUint64(key string, v uint64) error
	WriteValueStr(key string, v string) error
	WriteValueBin(key string, v []byte) error
	ClearValue(key string) error
	ClearAll() error
	Commit() error
	HasValue(key string) bool
}

type PersistentStorageImpl struct {
	storage      Storage
	mConfigFile  string
	mLock        sync.Locker
	mDirty       bool
	mInitialized bool
}

func NewPersistentStorageImpl() *PersistentStorageImpl {
	return &PersistentStorageImpl{
		storage:      NewIniStorage(),
		mConfigFile:  "",
		mLock:        &sync.Mutex{},
		mDirty:       false,
		mInitialized: false,
	}
}

func (s *PersistentStorageImpl) Init(mConfigPath string) error {
	if s.mConfigFile == "" && s.mInitialized {
		return fmt.Errorf("initialized")
	}
	s.mConfigFile = mConfigPath
	s.storage = NewIniStorage()
	err := s.storage.AddConfig(mConfigPath)
	if err != nil {
		return err
	}
	s.mInitialized = true
	return nil
}

func (s *PersistentStorageImpl) ReadBoolValue(key string) (bool, error) {
	value, err := s.storage.GetUIntValue(key)
	if value == 0 {
		return true, err
	}
	return false, err
}

func (s *PersistentStorageImpl) ReadValueUint16(key string) (uint16, error) {
	s.mLock.Lock()
	defer s.mLock.Unlock()
	return s.storage.GetUInt16Value(key)
}

func (s *PersistentStorageImpl) ReadValueUint32(key string) (uint32, error) {
	s.mLock.Lock()
	defer s.mLock.Unlock()
	value, err := s.storage.GetUInt64Value(key)
	return uint32(value), err
}

func (s *PersistentStorageImpl) ReadValueUint64(key string) (uint64, error) {
	return s.storage.GetUInt64Value(key)
}

func (s *PersistentStorageImpl) ReadValueStr(key string) (string, error) {

	return s.storage.GetStringValue(key)
}

func (s *PersistentStorageImpl) ReadValueBin(key string) ([]byte, error) {
	return s.storage.GetBinaryBlobValue(key)
}

func (s *PersistentStorageImpl) ReadValue(key string) ([]byte, error) {
	return s.storage.GetBinaryBlobValue(key)
}

func (s *PersistentStorageImpl) WriteValueBool(key string, v bool) error {
	if v {
		return s.WriteValueUint16(key, 1)
	}
	err := s.WriteValueUint16(key, 0)
	s.mDirty = true
	return err
}

func (s *PersistentStorageImpl) WriteValueUint16(key string, v uint16) error {
	err := s.storage.AddEntry(key, strconv.FormatUint(uint64(v), 10))
	s.mDirty = true
	return err
}

func (s *PersistentStorageImpl) WriteValueUint32(key string, v uint32) error {
	err := s.storage.AddEntry(key, strconv.FormatUint(uint64(v), 10))
	s.mDirty = true
	return err
}

func (s *PersistentStorageImpl) WriteValueUint64(key string, v uint64) error {
	err := s.storage.AddEntry(key, strconv.FormatUint(uint64(v), 10))
	s.mDirty = true
	return err
}

func (s *PersistentStorageImpl) WriteValueStr(key string, v string) error {
	err := s.storage.AddEntry(key, v)
	if err != nil {
		return err
	}
	s.mDirty = true
	return s.Commit()

}

func (s *PersistentStorageImpl) WriteValueBin(key string, v []byte) error {
	err := s.storage.AddEntry(key, string(v))
	s.mDirty = true
	return err
}

func (s *PersistentStorageImpl) ClearValue(key string) error {

	err := s.storage.RemoveEntry(key)
	s.mDirty = true
	return err
}

func (s *PersistentStorageImpl) ClearAll() error {

	err := s.storage.RemoveAll()
	s.mDirty = true
	return err
}

func (s *PersistentStorageImpl) Commit() error {
	if s.mConfigFile != "" && s.mDirty {
		return s.storage.CommitConfig(s.mConfigFile)
	}
	s.mDirty = false
	return nil
}

func (s *PersistentStorageImpl) HasValue(key string) bool {
	return s.storage.HasValue(key)
}
