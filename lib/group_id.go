package lib

import "fmt"

type GroupId uint16

const (
	kUndefinedGroupId    GroupId = 0
	kMinUniversalGroupId GroupId = 0x8000
	kMaxUniversalGroupId GroupId = 0xFFFF

	kMinFabricGroupId GroupId = 0x0001
	kMaxFabricGroupId GroupId = 0x7FFF

	kAllNodes     GroupId = 0xFFFF
	kAllNonSleepy GroupId = 0xFFFE
	kAllProxies   GroupId = 0xFFFD

	kMinUniversalGroupIdReserved = 0x8000
	kMaxUniversalGroupIdReserved = 0xFFFC
)

func UndefinedGroupId() GroupId {
	return kUndefinedGroupId
}

func (i GroupId) String() string {
	value := uint16(i)
	return fmt.Sprintf("%04X", value)
}

func (i GroupId) IsFabric() bool {
	return i >= kMinFabricGroupId && i <= kMaxFabricGroupId
}

func (i GroupId) IsUniversal() bool {
	return i >= kMinUniversalGroupId
}

func (i GroupId) IsValid() bool {
	return i != kUndefinedGroupId
}

func (i GroupId) IsOperationalGroupId() bool {
	return i != kUndefinedGroupId && (i < kMinUniversalGroupId || i > kMaxUniversalGroupId)
}

func (i GroupId) NodeId() NodeId {
	return NodeId(uint64(kMinGroupNodeId) | uint64(i))
}
