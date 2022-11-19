package store

type KvsPersistentStorageBase interface {
	Init(kvsMgr Storage)
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

type KvsPersistentStorageImpl struct {
	mStorage Storage
}

func NewKvsPersistentStorageImpl() *KvsPersistentStorageImpl {
	return &KvsPersistentStorageImpl{}
}

func (p *KvsPersistentStorageImpl) Init(kvsMgr Storage) {
	p.mStorage = kvsMgr
}

func (p *KvsPersistentStorageImpl) GetBool(key string) (bool, error) {
	return p.mStorage.ReadValueBool(key)
}

func (p *KvsPersistentStorageImpl) GetUint8(key string) (v uint8, e error) {
	value, e := p.mStorage.ReadValueUint64(key)
	return uint8(value), e
}

func (p *KvsPersistentStorageImpl) GetUint16(key string) (v uint16, e error) {
	value, e := p.mStorage.ReadValueUint64(key)
	return uint16(value), e
}

func (p *KvsPersistentStorageImpl) GetUint32(key string) (v uint32, e error) {
	value, e := p.mStorage.ReadValueUint64(key)
	return uint32(value), e
}

func (p *KvsPersistentStorageImpl) GetUint64(key string) (v uint64, e error) {
	return p.mStorage.ReadValueUint64(key)
}

func (p *KvsPersistentStorageImpl) GetString(key string) (v string, e error) {
	return p.mStorage.ReadValueString(key)
}

func (p *KvsPersistentStorageImpl) GetBytes(key string) (v []byte, e error) {
	value, e := p.mStorage.ReadValueString(key)
	return []byte(value), e
}

func (p *KvsPersistentStorageImpl) SetBool(key string, v bool) error {
	return p.mStorage.WriteValueBool(key, v)
}

func (p *KvsPersistentStorageImpl) SetUint8(key string, v uint16) error {
	return p.mStorage.WriteValueUint64(key, uint64(v))
}

func (p *KvsPersistentStorageImpl) SetUint16(key string, v uint16) error {
	return p.mStorage.WriteValueUint64(key, uint64(v))
}

func (p *KvsPersistentStorageImpl) SetUint32(key string, v uint32) error {
	return p.mStorage.WriteValueUint64(key, uint64(v))
}

func (p *KvsPersistentStorageImpl) SetUint64(key string, v uint64) error {
	return p.mStorage.WriteValueUint64(key, v)
}

func (p *KvsPersistentStorageImpl) SetString(key string, v string) error {
	return p.mStorage.WriteValueString(key, v)
}

func (p *KvsPersistentStorageImpl) SetBytes(key string, v []byte) error {
	return p.mStorage.WriteValueString(key, string(v))
}

func (p *KvsPersistentStorageImpl) DeleteKeyValue(key string) error {
	return p.mStorage.Delete(key)
}

func (p *KvsPersistentStorageImpl) DoesKeyExist(key string) bool {
	return p.mStorage.HasValue(key)
}
