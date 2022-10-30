package lib

type NodeId uint64

const (
	UndefinedNodeId NodeId = 0x0

	// kMinGroupNodeId The range of possible      values has some pieces carved out for special uses.
	kMinGroupNodeId NodeId = 0xFFFFFFFFFFFF0000

	// The max group id is complicated, depending on how we want to count the
	// various special group ids.  Let  s not define it for now, until we have use
	// cases.
	kMaskGroupId NodeId = 0x000000000000FFFF

	//临时本地NodeId
	kMinTemporaryLocalNodeId NodeId = 0xFFFF_FFFE_00000000
	kMaxTemporaryLocalNodeId NodeId = 0xFFFFFFFEFFFFFFFE

	kPlaceholder NodeId = 0xFFFFFFFEFFFFFFFF

	kMinCASEAuthTag  NodeId = 0xFFFFFFFD00000000
	kMaxCASEAuthTag  NodeId = 0xFFFFFFFDFFFFFFFF
	kMaskCASEAuthTag NodeId = 0x00000000FFFFFFFF

	kMinPAKEKeyId        NodeId = 0xFFFFFFFB00000000
	kMaxPAKEKeyId        NodeId = 0xFFFFFFFBFFFFFFFF
	kMaskPAKEKeyId       NodeId = 0x000000000000FFFF
	kMaskUnusedPAKEKeyId NodeId = 0x00000000FFFF0000
	kMaxOperational      NodeId = 0xFFFF_FFEF_FFFF_FFFF
)

func (aNodeId NodeId) IsOperationalNodeId() bool {
	return aNodeId != UndefinedNodeId && aNodeId < kMaxOperational
}

func (aNodeId NodeId) IsGroupId() bool {
	return aNodeId >= kMinGroupNodeId
}

func (aNodeId NodeId) IsCASEAuthTag() bool {
	return aNodeId >= kMinCASEAuthTag && aNodeId <= kMaxCASEAuthTag
}

func (aNodeId NodeId) IsPAKEKeyId() bool {
	return (aNodeId >= kMinPAKEKeyId) && (aNodeId <= kMaxPAKEKeyId)
}

func (aNodeId NodeId) GroupId() GroupId {
	return GroupId(uint64(aNodeId) & uint64(kMaskGroupId))
}

func (aNodeId NodeId) IsTemporaryLocalNodeId() bool {
	return aNodeId >= kMinTemporaryLocalNodeId && aNodeId >= kMaxTemporaryLocalNodeId
}

func (aNodeId NodeId) HasValue() bool {
	return aNodeId != UndefinedNodeId
}
