package store

import (
	"fmt"
	"github.com/spf13/cast"
	"strconv"
	"sync"
)

type KvsStorage interface {
	WriteValue(key string, val any) error

	ReadValueUint(key string) (uint64, error)
	WriteValueUint(string, uint64) error

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

type KvsPersistentStorage struct {
	kvs          *iniFile
	mConfigFile  string
	mLock        sync.Locker
	mDirty       bool
	mInitialized bool
}

func (s *KvsPersistentStorage) Init(mConfigPath string) error {

	s.mConfigFile = mConfigPath
	if s.kvs == nil {
		s.kvs = newIniFile()
	}
	_ = s.kvs.init()
	err := s.kvs.addConfig(mConfigPath)
	if err != nil {
		return err
	}
	s.mInitialized = true
	return nil
}

func NewInitStorage(fileName string) *KvsPersistentStorage {
	s := &KvsPersistentStorage{
		kvs:          newIniFile(),
		mConfigFile:  "",
		mLock:        &sync.Mutex{},
		mDirty:       false,
		mInitialized: false,
	}
	if err := s.Init(fileName); err != nil {
		panic(err.Error())
	}
	return s
}

func (s *KvsPersistentStorage) ReadValueFloat(key string) (float64, error) {
	return s.kvs.readFloatValue(key)
}

func (s *KvsPersistentStorage) WriteValueFloat(key string, v float64) error {
	err := s.kvs.addEntry(key, strconv.FormatFloat(v, 'f', 10, 64))
	if err != nil {
		return err
	}
	s.mDirty = true
	return s.commit()
}

func (s *KvsPersistentStorage) ReadValueBool(key string) (bool, error) {
	value, err := s.kvs.readUintValue(key)
	if value == 0 {
		return false, err
	}
	return true, err
}

func (s *KvsPersistentStorage) ReadValueInt(key string) (int64, error) {
	s.mLock.Lock()
	defer s.mLock.Unlock()
	value, err := s.kvs.readUint64Value(key)
	return int64(value), err

}

func (s *KvsPersistentStorage) WriteValue(k string, v any) error {
	switch v.(type) {
	case int, int8, int16, int32, int64:
		val := cast.ToInt64(v)
		err := s.kvs.addEntry(k, strconv.FormatInt(val, 10))
		if err != nil {
			return err
		}
	case uint, uint8, uint16, uint32, uint64:
		val := cast.ToUint64(v)
		err := s.kvs.addEntry(k, strconv.FormatUint(val, 10))
		if err != nil {
			return err
		}
	case bool:
		val := cast.ToBool(v)
		err := s.kvs.addEntry(k, strconv.FormatBool(val))
		if err != nil {
			return err
		}
	case string:
		val := cast.ToString(v)
		err := s.kvs.addEntry(k, val)
		if err != nil {
			return err
		}
	case []byte:
		val, _ := v.([]byte)
		err := s.kvs.addEntry(k, string(val))
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("type err")
	}
	s.mDirty = true
	return s.commit()
}

func (s *KvsPersistentStorage) WriteValueInt(k string, v int64) error {
	err := s.kvs.addEntry(k, strconv.FormatInt(v, 10))
	if err != nil {
		return err
	}
	s.mDirty = true
	return s.commit()
}

func (s *KvsPersistentStorage) ReadValueString(key string) (string, error) {
	return s.kvs.readStringValue(key)
}

func (s *KvsPersistentStorage) WriteValueString(key string, v string) error {
	err := s.kvs.addEntry(key, v)
	if err != nil {
		return err
	}
	s.mDirty = true
	return s.commit()
}

func (s *KvsPersistentStorage) Delete(key string) error {
	err := s.kvs.removeEntry(key)
	if err != nil {
		return err
	}
	s.mDirty = true
	return s.commit()
}

func (s *KvsPersistentStorage) DeleteAll() error {
	err := s.kvs.init()
	if err != nil {
		return err
	}
	s.mDirty = true
	return s.commit()
}

func (s *KvsPersistentStorage) ReadValueUint16(key string) (uint16, error) {
	s.mLock.Lock()
	defer s.mLock.Unlock()
	return s.kvs.readUInt16Value(key)
}

func (s *KvsPersistentStorage) ReadValueUint(key string) (uint64, error) {
	return s.kvs.readUint64Value(key)
}

func (s *KvsPersistentStorage) ReadValueStr(key string) (string, error) {
	return s.kvs.readStringValue(key)
}

func (s *KvsPersistentStorage) ReadValueBin(key string) ([]byte, error) {
	return s.kvs.readBinaryValue(key)
}

func (s *KvsPersistentStorage) WriteValueBool(key string, v bool) error {
	if v {
		return s.WriteValueUint16(key, 1)
	}
	err := s.WriteValueUint16(key, 0)
	s.mDirty = true
	return err
}

func (s *KvsPersistentStorage) WriteValueUint16(key string, v uint16) error {
	err := s.kvs.addEntry(key, strconv.FormatUint(uint64(v), 10))
	s.mDirty = true
	return err
}

func (s *KvsPersistentStorage) WriteValueUint32(key string, v uint32) error {
	err := s.kvs.addEntry(key, strconv.FormatUint(uint64(v), 10))
	s.mDirty = true
	return err
}

func (s *KvsPersistentStorage) WriteValueUint(key string, v uint64) error {
	err := s.kvs.addEntry(key, strconv.FormatUint(v, 10))
	s.mDirty = true
	return err
}

func (s *KvsPersistentStorage) commit() error {
	if s.mConfigFile != "" && s.mDirty && s.mInitialized {
		err := s.kvs.commitConfig(s.mConfigFile)
		if err != nil {
			return err
		}
		s.mDirty = false
	} else {
		return fmt.Errorf("strore state  error")
	}
	return nil
}

func (s *KvsPersistentStorage) HasValue(key string) bool {
	return s.kvs.hasValue(key)
}
