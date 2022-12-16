package credentials

import (
	"github.com/galenliu/chip/crypto"
	"github.com/galenliu/chip/lib"
	"github.com/galenliu/chip/lib/store"
	"sync/atomic"
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
	Key            crypto.SymmetricKeyContextBase
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

type GroupDataProviderBase interface {
	SetStorageDelegate(delegate store.PersistentStorageDelegate)
	Init() error
	SetListener(listener GroupListener)
	GetIpkKeySet(index lib.FabricIndex) (*KeySet, error)
	GroupSessions(sessionId uint16) []*GroupSession
}

var gGroupDataProvider atomic.Value

func init() {
	gGroupDataProvider.Store(&GroupDataProvider{})
}

func GetGroupDataProvider() *GroupDataProvider {
	return gGroupDataProvider.Load().(*GroupDataProvider)
}

func SetGroupDataProvider(g *GroupDataProvider) {
	gGroupDataProvider.Store(g)
}

type GroupDataProvider struct {
	mStorage       store.PersistentStorageDelegate
	mGroupListener GroupListener
}

func NewGroupDataProvider() *GroupDataProvider {
	return &GroupDataProvider{}
}

func (g *GroupDataProvider) GetIpkKeySet(index lib.FabricIndex) (outKeyset *KeySet, err error) {
	outKeyset = &KeySet{}
	fabric := &FabricData{
		FabricIndex: index,
	}
	if err = fabric.Load(g.mStorage); err != nil {
		return nil, err
	}

	//mapping := NewKeyMapData(fabric.FabricIndex, fabric.FirstMap)

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

	return
}

func (g *GroupDataProvider) SetListener(listener GroupListener) {
	g.mGroupListener = listener
}

func (g *GroupDataProvider) SetStorageDelegate(delegate store.PersistentStorageDelegate) {
	g.mStorage = delegate
}

func (g *GroupDataProvider) GroupSessions(sessionId uint16) []*GroupSession {
	return nil
}

func (g *GroupDataProvider) Init() error {
	return nil
}

func (g *GroupDataProvider) KeyContext(fabricIndex lib.FabricIndex, groupId lib.GroupId) crypto.SymmetricKeyContextBase {
	fabric := &FabricData{
		FabricIndex: fabricIndex,
	}
	if err := fabric.Load(g.mStorage); err != nil {
		return nil
	}
	mapping := KeyMapData{
		fabricIndex: 0,
		groupId:     0,
		keysetId:    0,
		LinkedData: LinkedData{
			id: fabric.FirstMap,
		},
	}
	var i uint16 = 0
	for ; i < fabric.MapCount; i++ {
		mapping.id = mapping.next
		if err := mapping.Load(g.mStorage); err != nil {
			return nil
		}
		if mapping.keysetId > 0 && mapping.groupId == groupId {

		}
	}

	return nil
}

type PersistentData struct {
}
