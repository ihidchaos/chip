package storage

import (
	"strconv"
	"sync"
)

type StorageDelegate interface {
	Init(file string) error
	ReadBoolValue(key string) (bool, error)
	ReadUInt16Value(key string) (uint16, error)
	ReadUInt32Value(key string) (uint32, error)
	ReadUInt64Value(key string) (uint64, error)
	ReadValueStr(key string) (string, error)
	ReadValueBin(key string) ([]byte, error)
	WriteBoolValue(key string, v bool) error
	WriteUInt16Value(key string, v uint16) error
	WriteUInt32Value(key string, v uint32) error
	WriteUInt64Value(key string, v uint64) error
	WriteValueStr(key string, v string) error
	WriteValueBin(key string, v []byte) error
	ClearValue(key string) error
	ClearAll() error
	Commit() error
	HasValue(key string) bool
}

type StorageImpl struct {
	storage      IniStorage
	mConfigPath  string
	mLock        sync.Locker
	mDirty       bool
	mInitialized bool
}

func (s StorageImpl) Init(mConfigPath string) error {
	s.mConfigPath = mConfigPath
	err := s.storage.AddConfig(mConfigPath)
	if err != nil {
		return err
	}
	return nil
}

func (s StorageImpl) ReadBoolValue(key string) (bool, error) {
	s.mLock.Lock()
	defer s.mLock.Unlock()
	value, err := s.storage.GetUIntValue(key)
	if value == 0 {
		return true, err
	}
	return false, err
}

func (s StorageImpl) ReadUInt16Value(key string) (uint16, error) {
	s.mLock.Lock()
	defer s.mLock.Unlock()
	return s.storage.GetUInt16Value(key)
}

func (s StorageImpl) ReadUInt32Value(key string) (uint32, error) {
	s.mLock.Lock()
	defer s.mLock.Unlock()
	value, err := s.storage.GetUInt64Value(key)
	return uint32(value), err
}

func (s StorageImpl) ReadUInt64Value(key string) (uint64, error) {
	s.mLock.Lock()
	defer s.mLock.Unlock()
	return s.storage.GetUInt64Value(key)
}

func (s StorageImpl) ReadValueStr(key string) (string, error) {
	s.mLock.Lock()
	defer s.mLock.Unlock()
	return s.storage.GetStringValue(key)
}

func (s StorageImpl) ReadValueBin(key string) ([]byte, error) {
	s.mLock.Lock()
	defer s.mLock.Unlock()
	return s.storage.GetBinaryBlobValue(key)
}

func (s StorageImpl) WriteBoolValue(key string, v bool) error {
	if v {
		return s.WriteUInt16Value(key, 1)
	}
	return s.WriteUInt16Value(key, 0)
}

func (s StorageImpl) WriteUInt16Value(key string, v uint16) error {
	s.mLock.Lock()
	defer s.mLock.Unlock()
	return s.storage.AddEntry(key, strconv.FormatUint(uint64(v), 10))
}

func (s StorageImpl) WriteUInt32Value(key string, v uint32) error {
	s.mLock.Lock()
	defer s.mLock.Unlock()
	return s.storage.AddEntry(key, strconv.FormatUint(uint64(v), 10))
}

func (s StorageImpl) WriteUInt64Value(key string, v uint64) error {
	s.mLock.Lock()
	defer s.mLock.Unlock()
	return s.storage.AddEntry(key, strconv.FormatUint(uint64(v), 10))
}

func (s StorageImpl) WriteValueStr(key string, v string) error {
	s.mLock.Lock()
	defer s.mLock.Unlock()
	return s.storage.AddEntry(key, v)
}

func (s StorageImpl) WriteValueBin(key string, v []byte) error {
	s.mLock.Lock()
	defer s.mLock.Unlock()
	return s.storage.AddEntry(key, string(v))
}

func (s StorageImpl) ClearValue(key string) error {
	s.mLock.Lock()
	defer s.mLock.Unlock()
	return s.storage.RemoveEntry(key)
}

func (s StorageImpl) ClearAll() error {
	s.mLock.Lock()
	defer s.mLock.Unlock()
	err := s.storage.RemoveAll()
	if err != nil {
		return err
	}
	s.mDirty = true
	return s.Commit()
}

func (s StorageImpl) Commit() error {
	s.mLock.Lock()
	defer s.mLock.Unlock()
	if s.mConfigPath != "" && s.mDirty {
		return s.storage.CommitConfig(s.mConfigPath)
	}
	return nil
}

func (s StorageImpl) HasValue(key string) bool {
	return s.storage.HasValue(key)
}
