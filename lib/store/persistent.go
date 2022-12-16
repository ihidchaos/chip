package store

import (
	"fmt"
	"github.com/spf13/cast"
)

type PersistentStorageDelegate interface {
	SetStorage(storage KvsStorage)
	//GetBool(key string) (bool, error)
	//GetUint8(key string) (uint8, error)
	//GetUint16(key string) (uint16, error)
	//GetUint32(key string) (uint32, error)
	//GetUint64(key string) (uint64, error)
	//GetString(key string) (string, error)
	//GetBytes(key string) ([]byte, error)
	//
	//SetBool(key string, v bool) error
	//SetUint8(key string, v uint16) error
	//SetUint16(key string, v uint16) error
	//SetUint32(key string, v uint32) error
	//SetUint64(key string, v uint64) error
	//SetString(key string, v string) error
	//SetBytes(key string, v []byte) error

	DeleteKeyValue(key string) error
	DoesKeyExist(key string) bool
	SetKeyValue(key string, value any) error
	GetKeyValue(key string, outValue any) error
}

type DefaultPersistentImpl struct {
	mStorage KvsStorage
}

func (p *DefaultPersistentImpl) SetStorage(storage KvsStorage) {
	p.mStorage = storage
}

//func (p *DefaultPersistentImpl) GetBool(key string) (bool, error) {
//	return p.mStorage.ReadValueBool(key)
//}
//
//func (p *DefaultPersistentImpl) GetUint8(key string) (v uint8, e error) {
//	value, e := p.mStorage.ReadValueUint(key)
//	return uint8(value), e
//}
//
//func (p *DefaultPersistentImpl) GetUint16(key string) (v uint16, e error) {
//	value, e := p.mStorage.ReadValueUint(key)
//	return uint16(value), e
//}
//
//func (p *DefaultPersistentImpl) GetUint32(key string) (v uint32, e error) {
//	value, e := p.mStorage.ReadValueUint(key)
//	return uint32(value), e
//}
//
//func (p *DefaultPersistentImpl) GetUint64(key string) (v uint64, e error) {
//	return p.mStorage.ReadValueUint(key)
//}
//
//func (p *DefaultPersistentImpl) GetString(key string) (v string, e error) {
//	return p.mStorage.ReadValueString(key)
//}
//
//func (p *DefaultPersistentImpl) GetBytes(key string) (v []byte, e error) {
//	value, e := p.mStorage.ReadValueString(key)
//	return []byte(value), e
//}
//
//func (p *DefaultPersistentImpl) SetBool(key string, v bool) error {
//	return p.mStorage.WriteValueBool(key, v)
//}
//
//func (p *DefaultPersistentImpl) SetUint8(key string, v uint16) error {
//	return p.mStorage.WriteValueUint(key, uint64(v))
//}
//
//func (p *DefaultPersistentImpl) SetUint16(key string, v uint16) error {
//	return p.mStorage.WriteValueUint(key, uint64(v))
//}
//
//func (p *DefaultPersistentImpl) SetUint32(key string, v uint32) error {
//	return p.mStorage.WriteValueUint(key, uint64(v))
//}
//
//func (p *DefaultPersistentImpl) SetUint64(key string, v uint64) error {
//	return p.mStorage.WriteValueUint(key, v)
//}
//
//func (p *DefaultPersistentImpl) SetString(key string, v string) error {
//	return p.mStorage.WriteValueString(key, v)
//}
//
//func (p *DefaultPersistentImpl) SetBytes(key string, v []byte) error {
//	return p.mStorage.WriteValueString(key, string(v))
//}

func (p *DefaultPersistentImpl) DeleteKeyValue(key string) error {
	return p.mStorage.Delete(key)
}

func (p *DefaultPersistentImpl) DoesKeyExist(key string) bool {
	return p.mStorage.HasValue(key)
}

func (p *DefaultPersistentImpl) SetKeyValue(key string, value any) error {
	switch value.(type) {
	case uint, uint8, uint16, uint32, uint64:
		v := cast.ToUint64(value)
		return p.mStorage.WriteValueUint(key, v)
	case int, int8, int16, int32, int64:
		v := cast.ToInt64(value)
		return p.mStorage.WriteValueInt(key, v)
	case bool:
		var v int64 = -1
		if cast.ToBool(value) {
			v = 1
		}
		return p.mStorage.WriteValueInt(key, v)
	case string:
		v := cast.ToString(value)
		return p.mStorage.WriteValueString(key, v)
	default:
		return fmt.Errorf("value type invaild")
	}
}

func (p *DefaultPersistentImpl) GetKeyValue(key string, outValue any) error {
	switch outValue.(type) {
	case *int:
		if v, err := p.mStorage.ReadValueInt(key); err != nil {
			return err
		} else {
			out := int(v)
			outValue = &out
		}
	case *int8:
		if v, err := p.mStorage.ReadValueInt(key); err != nil {
			return err
		} else {
			out := int8(v)
			outValue = &out
		}
	case *int16:
		if v, err := p.mStorage.ReadValueInt(key); err != nil {
			return err
		} else {
			out := int16(v)
			outValue = &out
		}
	case *int32:
		if v, err := p.mStorage.ReadValueInt(key); err != nil {
			return err
		} else {
			out := int32(v)
			outValue = &out
		}
	case *int64:
		if v, err := p.mStorage.ReadValueInt(key); err != nil {
			return err
		} else {
			outValue = &v
		}
	case *uint:
		if v, err := p.mStorage.ReadValueUint(key); err != nil {
			return err
		} else {
			out := uint(v)
			outValue = &out
		}
	case *uint8:
		if v, err := p.mStorage.ReadValueUint(key); err != nil {
			return err
		} else {
			out := uint8(v)
			outValue = &out
		}

	case *uint16:
		if v, err := p.mStorage.ReadValueUint(key); err != nil {
			return err
		} else {
			out := uint16(v)
			outValue = &out
		}
	case *uint32:
		if v, err := p.mStorage.ReadValueUint(key); err != nil {
			return err
		} else {
			out := uint32(v)
			outValue = &out
		}

	case *uint64:
		if v, err := p.mStorage.ReadValueUint(key); err != nil {
			return err
		} else {
			outValue = &v
		}

	case *bool:
		if v, err := p.mStorage.ReadValueInt(key); err != nil {
			return err
		} else {
			out := false
			if v > 0 {
				out = true
			}
			outValue = &out
		}
	case *string:
		if v, err := p.mStorage.ReadValueString(key); err != nil {
			return err
		} else {
			out := v
			outValue = &out
		}

	default:
		return fmt.Errorf("vlaue type invaild")
	}
	return nil
}
