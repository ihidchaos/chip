package store

import (
	"fmt"
	"github.com/spf13/cast"
	"strconv"
	"sync"
)

type Storage interface {
	Init(mConfigPath string) error

	WriteValue(key string, val any) error

	ReadValueUint64(key string) (uint64, error)
	WriteValueUint64(string, uint64) error

	ReadValueFloat(key string) (float64, error)
	WriteValueFloat(string, float64) error

	ReadValueBool(key string) (bool, error)
	WriteValueBool(string, bool) error

	ReadValueInt(key string) (int64, error)
	WriteValueInt(string, int64) error

	ReadValueString(key string) (string, error)
	WriteValueString(string, string) error

	HasValue(key string) bool

	Delete(key string) error

	DeleteAll() error
}

type PersistentStorageImpl struct {
	storage      ini
	mConfigFile  string
	mLock        sync.Locker
	mDirty       bool
	mInitialized bool
}

func NewPersistentStorageImpl() *PersistentStorageImpl {
	return &PersistentStorageImpl{
		storage:      newIniStorage(),
		mConfigFile:  "",
		mLock:        &sync.Mutex{},
		mDirty:       false,
		mInitialized: false,
	}
}

func (s *PersistentStorageImpl) ReadValueFloat(key string) (float64, error) {
	return s.storage.readFloatValue(key)
}

func (s *PersistentStorageImpl) WriteValueFloat(key string, v float64) error {
	err := s.storage.addEntry(key, strconv.FormatFloat(v, 'f', 10, 64))
	if err != nil {
		return err
	}
	s.mDirty = true
	return s.commit()
}

func (s *PersistentStorageImpl) ReadValueBool(key string) (bool, error) {
	value, err := s.storage.readUintValue(key)
	if value == 0 {
		return false, err
	}
	return true, err
}

func (s *PersistentStorageImpl) ReadValueInt(key string) (int64, error) {
	s.mLock.Lock()
	defer s.mLock.Unlock()
	value, err := s.storage.readUint64Value(key)
	return int64(value), err

}

func (s *PersistentStorageImpl) WriteValue(k string, v any) error {
	switch v.(type) {
	case int, int8, int16, int32, int64:
		val := cast.ToInt64(v)
		err := s.storage.addEntry(k, strconv.FormatInt(val, 10))
		if err != nil {
			return err
		}
	case uint, uint8, uint16, uint32, uint64:
		val := cast.ToUint64(v)
		err := s.storage.addEntry(k, strconv.FormatUint(val, 10))
		if err != nil {
			return err
		}
	case bool:
		val := cast.ToBool(v)
		err := s.storage.addEntry(k, strconv.FormatBool(val))
		if err != nil {
			return err
		}
	case string:
		val := cast.ToString(v)
		err := s.storage.addEntry(k, val)
		if err != nil {
			return err
		}
	case []byte:
		val, _ := v.([]byte)
		err := s.storage.addEntry(k, string(val))
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("type err")
	}
	s.mDirty = true
	return s.commit()
}

func (s *PersistentStorageImpl) WriteValueInt(k string, v int64) error {
	err := s.storage.addEntry(k, strconv.FormatInt(v, 10))
	if err != nil {
		return err
	}
	s.mDirty = true
	return s.commit()
}

func (s *PersistentStorageImpl) ReadValueString(key string) (string, error) {
	return s.storage.readStringValue(key)
}

func (s *PersistentStorageImpl) WriteValueString(key string, v string) error {
	err := s.storage.addEntry(key, v)
	if err != nil {
		return err
	}
	s.mDirty = true
	return s.commit()
}

func (s *PersistentStorageImpl) Delete(key string) error {
	err := s.storage.removeEntry(key)
	if err != nil {
		return err
	}
	s.mDirty = true
	return s.commit()
}

func (s *PersistentStorageImpl) DeleteAll() error {
	err := s.storage.init()
	if err != nil {
		return err
	}
	s.mDirty = true
	return s.commit()
}

func (s *PersistentStorageImpl) Init(mConfigPath string) error {

	s.mConfigFile = mConfigPath
	if s.storage == nil {
		s.storage = newIniStorage()
	}
	_ = s.storage.init()
	err := s.storage.addConfig(mConfigPath)
	if err != nil {
		return err
	}
	s.mInitialized = true
	return nil
}

func (s *PersistentStorageImpl) ReadValueUint16(key string) (uint16, error) {
	s.mLock.Lock()
	defer s.mLock.Unlock()
	return s.storage.readUInt16Value(key)
}

func (s *PersistentStorageImpl) ReadValueUint64(key string) (uint64, error) {
	return s.storage.readUint64Value(key)
}

func (s *PersistentStorageImpl) ReadValueStr(key string) (string, error) {
	return s.storage.readStringValue(key)
}

func (s *PersistentStorageImpl) ReadValueBin(key string) ([]byte, error) {
	return s.storage.readBinaryValue(key)
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
	err := s.storage.addEntry(key, strconv.FormatUint(uint64(v), 10))
	s.mDirty = true
	return err
}

func (s *PersistentStorageImpl) WriteValueUint32(key string, v uint32) error {
	err := s.storage.addEntry(key, strconv.FormatUint(uint64(v), 10))
	s.mDirty = true
	return err
}

func (s *PersistentStorageImpl) WriteValueUint64(key string, v uint64) error {
	err := s.storage.addEntry(key, strconv.FormatUint(v, 10))
	s.mDirty = true
	return err
}

func (s *PersistentStorageImpl) commit() error {
	if s.mConfigFile != "" && s.mDirty && s.mInitialized {
		err := s.storage.commitConfig(s.mConfigFile)
		if err != nil {
			return err
		}
		s.mDirty = false
	} else {
		return fmt.Errorf("strore state  error")
	}
	return nil
}

func (s *PersistentStorageImpl) HasValue(key string) bool {
	return s.storage.hasValue(key)
}