package storage

type keyValueStoreManager interface {
	ReadValueUint(key string) (uint64, error)
	WriteValueUint(string, uint64)

	ReadValueFloat(key string) (float64, error)
	WriteValueUint(string, uint64)
}

type KvsPersistentStorageDelegate interface {
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
	mKvsManager keyValueStoreManager
}

func (p *KvsPersistentStorageDelegateImpl) GetBool(key string) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (p *KvsPersistentStorageDelegateImpl) GetUint8(key string) (uint8, error) {
	//TODO implement me
	panic("implement me")
}

func (p *KvsPersistentStorageDelegateImpl) GetUint16(key string) (uint16, error) {
	//TODO implement me
	panic("implement me")
}

func (p *KvsPersistentStorageDelegateImpl) GetUint32(key string) (uint32, error) {
	//TODO implement me
	panic("implement me")
}

func (p *KvsPersistentStorageDelegateImpl) GetUint64(key string) (uint64, error) {
	//TODO implement me
	panic("implement me")
}

func (p *KvsPersistentStorageDelegateImpl) GetString(key string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (p *KvsPersistentStorageDelegateImpl) GetBytes(key string) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (p *KvsPersistentStorageDelegateImpl) SetBool(key string, v bool) error {
	//TODO implement me
	panic("implement me")
}

func (p *KvsPersistentStorageDelegateImpl) SetUint8(key string, v uint16) error {
	//TODO implement me
	panic("implement me")
}

func (p *KvsPersistentStorageDelegateImpl) SetUint16(key string, v uint16) error {
	//TODO implement me
	panic("implement me")
}

func (p *KvsPersistentStorageDelegateImpl) SetUint32(key string, v uint32) error {
	//TODO implement me
	panic("implement me")
}

func (p *KvsPersistentStorageDelegateImpl) SetUint64(key string, v uint64) error {
	//TODO implement me
	panic("implement me")
}

func (p *KvsPersistentStorageDelegateImpl) SetString(key string, v string) error {
	//TODO implement me
	panic("implement me")
}

func (p *KvsPersistentStorageDelegateImpl) SetBytes(key string, v []byte) error {
	//TODO implement me
	panic("implement me")
}

func (p *KvsPersistentStorageDelegateImpl) DeleteKeyValue(key string) error {
	//TODO implement me
	panic("implement me")
}

func (p *KvsPersistentStorageDelegateImpl) DoesKeyExist(key string) bool {
	//TODO implement me
	panic("implement me")
}

func (p *KvsPersistentStorageDelegateImpl) Init(kvsMgr keyValueStoreManager) {
	p.mKvsManager = kvsMgr
}
