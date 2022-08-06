package storage

type KvsPersistentStorageDelegate interface {
	Init(kvsMgr ChipStorage)
	GetBool(key string) (bool, error)
	GetUint8(key string) (uint8, error)
	GetUint16(key string) (uint16, error)
	GetUint32(key string) (uint32, error)
	GetUint64(key string) (uint64, error)
	GetString(key string) (string, error)
	GetBytes(key string) ([]byte, error)

	SetBool(key string, v bool) error
	SetUint8(key string, v uint16) error
	SetUint16(key string, v uint16) error
	SetUint32(key string, v uint32) error
	SetUint64(key string, v uint64) error
	SetString(key string, v string) error
	SetBytes(key string, v []byte) error

	DeleteKeyValue(key string) error
	DoesKeyExist(key string) bool
}

type KvsPersistentStorageDelegateImpl struct {
	mKvsManager ChipStorage
}

func NewKvsPersistentStorageDelegateImpl() *KvsPersistentStorageDelegateImpl {
	return &KvsPersistentStorageDelegateImpl{}
}

func (p *KvsPersistentStorageDelegateImpl) Init(kvsMgr ChipStorage) {
	p.mKvsManager = kvsMgr
}

func (p *KvsPersistentStorageDelegateImpl) GetBool(key string) (bool, error) {
	return p.mKvsManager.ReadValueBool(key)
}

func (p *KvsPersistentStorageDelegateImpl) GetUint8(key string) (v uint8, e error) {
	value, e := p.mKvsManager.ReadValueUint64(key)
	if e != nil {
		return
	}
	return uint8(value), nil
}

func (p *KvsPersistentStorageDelegateImpl) GetUint16(key string) (v uint16, e error) {
	value, e := p.mKvsManager.ReadValueUint64(key)
	if e != nil {
		return
	}
	return uint16(value), nil
}

func (p *KvsPersistentStorageDelegateImpl) GetUint32(key string) (v uint32, e error) {
	value, e := p.mKvsManager.ReadValueUint64(key)
	if e != nil {
		return
	}
	return uint32(value), nil
}

func (p *KvsPersistentStorageDelegateImpl) GetUint64(key string) (v uint64, e error) {
	return p.mKvsManager.ReadValueUint64(key)
}

func (p *KvsPersistentStorageDelegateImpl) GetString(key string) (v string, e error) {
	return p.mKvsManager.ReadValueString(key)
}

func (p *KvsPersistentStorageDelegateImpl) GetBytes(key string) (v []byte, e error) {
	value, e := p.mKvsManager.ReadValueString(key)
	if e != nil {
		return
	}
	return []byte(value), nil
}

func (p *KvsPersistentStorageDelegateImpl) SetBool(key string, v bool) error {
	return p.mKvsManager.WriteValueBool(key, v)
}

func (p *KvsPersistentStorageDelegateImpl) SetUint8(key string, v uint16) error {
	return p.mKvsManager.WriteValueUint64(key, uint64(v))
}

func (p *KvsPersistentStorageDelegateImpl) SetUint16(key string, v uint16) error {
	return p.mKvsManager.WriteValueUint64(key, uint64(v))
}

func (p *KvsPersistentStorageDelegateImpl) SetUint32(key string, v uint32) error {
	return p.mKvsManager.WriteValueUint64(key, uint64(v))
}

func (p *KvsPersistentStorageDelegateImpl) SetUint64(key string, v uint64) error {
	return p.mKvsManager.WriteValueUint64(key, v)
}

func (p *KvsPersistentStorageDelegateImpl) SetString(key string, v string) error {
	return p.mKvsManager.WriteValueString(key, v)
}

func (p *KvsPersistentStorageDelegateImpl) SetBytes(key string, v []byte) error {
	return p.mKvsManager.WriteValueString(key, string(v))
}

func (p *KvsPersistentStorageDelegateImpl) DeleteKeyValue(key string) error {
	return p.mKvsManager.DeleteKeyValue(key)
}

func (p *KvsPersistentStorageDelegateImpl) DoesKeyExist(key string) bool {
	return p.mKvsManager.HasValue(key)
}
