package lib

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
