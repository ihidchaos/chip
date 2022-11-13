package lib

import "fmt"

type GroupId uint16

// 0xFF00-0xFFFC Reserved for future use
const (
	UndefinedGroupId    GroupId = 0x0000
	AllNodesGroupId     GroupId = 0xFFFF
	AllNonSleepyGroupId GroupId = 0xFFFE
	AllProxiesGroupId   GroupId = 0xFFFD

	kMinUniversalGroupId GroupId = 0xFF00
	kMaxUniversalGroupId GroupId = 0xFFFF

	kMinFabricGroupId GroupId = 0x0001
	kMaxFabricGroupId GroupId = 0x7FFF
)

func (i GroupId) String() string {
	value := uint16(i)
	return fmt.Sprintf("%04X", value)
}

func (i GroupId) Bytes() [4]byte {
	return [4]byte{
		byte(0x00),
		byte(0x00),
		byte((i & 0xFF00) >> 8),
		byte((i & 0x00FF) >> 0),
	}
}

func (i GroupId) IsFabricGroupId() bool {
	return i >= kMinFabricGroupId && i <= kMaxFabricGroupId
}

func (i GroupId) IsUniversalGroupId() bool {
	return i >= kMinUniversalGroupId
}

func (i GroupId) IsValidGroupId() bool {
	return i != UndefinedGroupId
}

func (i GroupId) HasValue() bool {
	return i != UndefinedGroupId
}

func (i GroupId) IsOperationalGroupId() bool {
	return i != UndefinedGroupId && (i < kMinUniversalGroupId || i > kMaxUniversalGroupId)
}
