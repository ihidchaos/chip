package lib

type GroupId uint16

// 0xFF00-0xFFFC Reserved for future use
const (
	KUndefinedGroupId    GroupId = 0x0000
	KMinUniversalGroupId GroupId = 0x8000
	KMaxUniversalGroupId GroupId = 0xFFFF

	KMinFabricGroupId GroupId = 0x0001
	KMaxFabricGroupId GroupId = 0x7FFF

	KAllNodes                    GroupId = 0xFFFF
	KAllNonSleepy                GroupId = 0xFFFE
	KAllProxies                  GroupId = 0xFFFD
	KMinUniversalGroupIdReserved GroupId = 0x8000
	KMaxUniversalGroupIdReserved GroupId = 0xFFFC
)

func (i GroupId) IsOperationalGroupId() bool {
	return i != KUndefinedGroupId && (i < KMinUniversalGroupIdReserved || i > KMaxUniversalGroupIdReserved)
}

func (i GroupId) IsFabricGroupId() bool {
	return i >= KMinFabricGroupId && i <= KMaxFabricGroupId
}

func (i GroupId) IsUniversalGroupId() bool {
	return i >= KMinUniversalGroupId
}

func (i GroupId) IsValidGroupId() bool {
	return i != KUndefinedGroupId
}

func (i GroupId) HasValue() bool {
	return i != KUndefinedGroupId
}
