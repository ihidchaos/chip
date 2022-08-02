package lib

type NodeId uint64

const (
	UndefinedNodeId NodeId = 0

	// The range of possible      values has some pieces carved out for special
	// uses.
	kMinGroupNodeId NodeId = 0xFFFFFFFFFFFF0000
	// The max group id is complicated, depending on how we want to count the
	// various special group ids.  Let  s not define it for now, until we have use
	// cases.
	kMaskGroupId NodeId = 0x000000000000FFFF

	kMinTemporaryLocalId NodeId = 0xFFFFFFFE00000000

	// We use the largest available temporary local id to represent
	// kPlaceholder    , so the range is narrowed compared to the spec.
	kMaxTemporaryLocalId NodeId = 0xFFFFFFFEFFFFFFFE
	kPlaceholder         NodeId = 0xFFFFFFFEFFFFFFFF

	kMinCASEAuthTag  NodeId = 0xFFFFFFFD00000000
	kMaxCASEAuthTag  NodeId = 0xFFFFFFFDFFFFFFFF
	kMaskCASEAuthTag NodeId = 0x00000000FFFFFFFF

	kMinPAKEKeyId        NodeId = 0xFFFFFFFB00000000
	kMaxPAKEKeyId        NodeId = 0xFFFFFFFBFFFFFFFF
	kMaskPAKEKeyId       NodeId = 0x000000000000FFFF
	kMaskUnusedPAKEKeyId NodeId = 0x00000000FFFF0000

	kMaxOperational NodeId = 0xFFFFFFEFFFFFFFFF
)

func IsOperationalNodeId(aNodeId NodeId) bool {
	return aNodeId != UndefinedNodeId && aNodeId < kMaxOperational
}

func IsGroupId(aNodeId NodeId) bool {
	return aNodeId >= kMinGroupNodeId
}

func IsCASEAuthTag(aNodeId NodeId) bool {
	return aNodeId >= kMinCASEAuthTag && aNodeId <= kMaxCASEAuthTag
}

func NodeIdFromGroupId(aGroupId GroupId) NodeId {
	return NodeId(uint64(kMinGroupNodeId) | uint64(aGroupId))
}

func GroupIdFromNodeId(aNodeId NodeId) GroupId {
	return GroupId(uint64(aNodeId) & uint64(kMaskGroupId))
}

func NodeIdFromPAKEKeyId(aPAKEKeyId PasscodeId) NodeId {
	return NodeId(uint64(kMinPAKEKeyId) | uint64(aPAKEKeyId))
}

func PAKEKeyIdFromNodeId(aNodeId NodeId) PasscodeId {
	return PasscodeId(aNodeId & kMaskPAKEKeyId)
}
