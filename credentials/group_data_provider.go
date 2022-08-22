package credentials

import (
	"github.com/galenliu/chip/crypto"
	"github.com/galenliu/chip/lib"
	"github.com/galenliu/chip/pkg/storage"
	"time"
)

const KEpochKeysMax = 3

type KeySet struct {
	NumKeysUsed int
	EpochKeys   []EpochKey
}

// 对称加密密钥长度
const lengthBytes = crypto.SymmetricKeyLengthBytes //Crypto::CHIP_CRYPTO_SYMMETRIC_KEY_LENGTH_BYTES;
type EpochKey struct {
	StartTime time.Time
	Key       []byte
}

func NewEpochKey() *EpochKey {
	return &EpochKey{
		StartTime: time.Time{},
		Key:       make([]byte, lengthBytes),
	}
}

type GroupInfo struct {
	Id lib.GroupId
}

type GroupDataProvider interface {
	SetStorageDelegate(delegate storage.KvsPersistentStorageDelegate)
	Init() error
	SetListener(listener GroupListener)
	GetIpkKeySet(index lib.FabricIndex) (*KeySet, error)
}

type GroupDataProviderImpl struct {
	mStorage       storage.KvsPersistentStorageDelegate
	mGroupListener GroupListener
}

func NewGroupDataProviderImpl() *GroupDataProviderImpl {
	return &GroupDataProviderImpl{}
}

func (g *GroupDataProviderImpl) GetIpkKeySet(index lib.FabricIndex) (*KeySet, error) {
	fabricData := NewFabricData(index)
	fabricData.Load(g.mStorage)
	return nil, nil
}

func (g *GroupDataProviderImpl) SetListener(listener GroupListener) {
	g.mGroupListener = listener
}

func (g *GroupDataProviderImpl) SetStorageDelegate(delegate storage.KvsPersistentStorageDelegate) {
	g.mStorage = delegate
}

func (g *GroupDataProviderImpl) Init() error {
	return nil
}

var gGroupDataProvider GroupDataProvider

func GetGroupDataProvider() GroupDataProvider {
	return gGroupDataProvider
}

func SetGroupDataProvider(g GroupDataProvider) {
	gGroupDataProvider = g
}

type PersistentData struct {
}
