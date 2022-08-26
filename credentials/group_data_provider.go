package credentials

import (
	"github.com/galenliu/chip/crypto"
	"github.com/galenliu/chip/lib"
	"github.com/galenliu/chip/pkg/storage"
	"time"
)

const KEpochKeysMax = 3

type KeySet struct {
	keysetId    lib.KeysetId
	NumKeysUsed uint8
	EpochKeys   []EpochKey
	Policy      any
}

// 对称加密密钥长度
const lengthBytes = crypto.SymmetricKeyLengthBytes //Crypto::CHIP_CRYPTO_SYMMETRIC_KEY_LENGTH_BYTES;
type EpochKey struct {
	StartTime time.Time
	Key       []byte
}

func (e *EpochKey) Clear() {
	e.StartTime = time.Time{}
	e.Key = make([]byte, lengthBytes)
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

func (g *GroupDataProviderImpl) GetIpkKeySet(index lib.FabricIndex) (outKeyset *KeySet, err error) {
	outKeyset = &KeySet{}
	fabric := NewFabricData(index)
	err = fabric.Load(g.mStorage)

	//mapping := NewKeyMapData(fabric.fabricIndex, fabric.firstMap)

	keyset := KeySetData{}
	keyset.Find(g.mStorage, fabric, lib.KeysetId(0))

	outKeyset.keysetId = keyset.keysetId
	outKeyset.NumKeysUsed = keyset.keysetCount
	outKeyset.Policy = keyset.policy

	for i, epoch := range outKeyset.EpochKeys {
		if uint8(i) < keyset.keysetCount {
			epoch.StartTime = keyset.operationalKeys[i].StateTime
			epoch.Key = keyset.operationalKeys[i].EncryptionKey
		}
	}

	if err != nil {
		return
	}
	return
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
