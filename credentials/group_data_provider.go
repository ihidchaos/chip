package credentials

import (
	"github.com/galenliu/chip/crypto"
	"github.com/galenliu/chip/lib"
	"github.com/galenliu/chip/pkg/store"
	"time"
)

type GroupListener interface {
	OnGroupAdded(fabricIndex lib.FabricIndex, newGroup *GroupInfo)
	OnGroupRemoved(fabricIndex lib.FabricIndex, newGroup *GroupInfo)
}

const KEpochKeysMax = 3

type SecurityPolicy uint8

const (
	TrustFirst SecurityPolicy = 1
)

type KeySet struct {
	keysetId    lib.KeysetId
	NumKeysUsed uint8
	EpochKeys   []EpochKey
	Policy      any
}

type EpochKey struct {
	StartTime time.Time
	Key       [crypto.SymmetricKeyLengthBytes]byte // 对称加密密钥长度
}

type GroupInfo struct {
	Id lib.GroupId
}

type GroupSession struct {
	GroupId        lib.GroupId
	FabricIndex    lib.FabricIndex
	SecurityPolicy SecurityPolicy
	Key            *crypto.SymmetricKeyContext
}

type GroupKey struct {
	groupId  lib.GroupId
	keysetId lib.KeysetId
}

type GroupEndpoint struct {
	GroupId    lib.GroupId
	EndPointId lib.EndpointId
}

func (e *EpochKey) Clear() {
	e.StartTime = time.Time{}
	e.Key = [crypto.SymmetricKeyLengthBytes]byte{}
}

func NewEpochKey() *EpochKey {
	return &EpochKey{
		StartTime: time.Time{},
		Key:       [crypto.SymmetricKeyLengthBytes]byte{},
	}
}

type GroupDataProvider interface {
	SetStorageDelegate(delegate store.KvsPersistentStorageBase)
	Init() error
	SetListener(listener GroupListener)
	GetIpkKeySet(index lib.FabricIndex) (*KeySet, error)
	GroupSessions(sessionId uint16) []*GroupSession
}

type GroupDataProviderImpl struct {
	mStorage       store.KvsPersistentStorageBase
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
			epoch.Key[i] = keyset.operationalKeys[i].EncryptionKey[i]
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

func (g *GroupDataProviderImpl) SetStorageDelegate(delegate store.KvsPersistentStorageBase) {
	g.mStorage = delegate
}

func (g *GroupDataProviderImpl) GroupSessions(sessionId uint16) []*GroupSession {
	return nil
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
