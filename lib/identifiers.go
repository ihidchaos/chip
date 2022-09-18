package lib

type VendorId uint16
type ProductId uint16
type GroupId uint16
type NodeId uint64
type FabricId uint64
type FabricIndex uint8
type CompressedFabricId uint64
type KeysetId uint16

const InvalidKeysetId KeysetId = 0xFFFF
const FabricIndexUndefined FabricIndex = 0
const ProductIdAnonymous ProductId = 0x0000

const (
	VendorIdMatterStandard VendorId = 0x0000
	VendorIdUnspecified    VendorId = 0x0000
	VendorIdApple          VendorId = 0x1349
	VendorIdGoogle         VendorId = 0x6006
	VendorTest1            VendorId = 0xFFF1
	VendorTest2            VendorId = 0xFFF2
	VendorTest3            VendorId = 0xFFF3
	VendorTest4            VendorId = 0xFFF4
	VendorIdNotSpecified   VendorId = 0xFFFF
)

const (
	UndefinedFabricIndex FabricIndex = 0x0
	MinValidFabricIndex  FabricIndex = 0x1
	kMaxValidFabricIndex FabricIndex = 0xFE
)

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

func (i GroupId) IsOperationalGroupId() bool {
	return i != UndefinedGroupId && (i < kMinUniversalGroupId || i > kMaxUniversalGroupId)
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

func (aNodeId NodeId) HasValue() bool {
	return aNodeId != UndefinedNodeId
}

//Range Type
//0xFFFF_FFFF_FFFF_xxxx Group Node ID
//0xFFFF_FFFF_0000_0000 to 0xFFFF_FFFF_FFFE_FFFF
//Reserved for future use
//0xFFFF_FFFE_xxxx_xxxx Temporary Local Node ID
//0xFFFF_FFFD_xxxx_xxxx CASE Authenticated Tag
//0xFFFF_FFFC_xxxx_xxxx to 0xFFFF_FFFC_FFFF_FFFF
//Reserved for future use
//0xFFFF_FFFB_xxxx_xxxx PAKE key identifiers
//0xFFFF_FFF0_0000_0000 to 0xFFFF_FFFA_FFFF_FFFF
//Reserved for future use
//0x0000_0000_0000_0001 to 0xFFFF_FFEF_FFFF_FFFF
//Operational Node ID
//0x0000_0000_0000_0000 Unspecified Node ID

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

func NodeIdFromPAKEKeyId(aPAKEKeyId PasscodeId) NodeId {
	return NodeId(uint64(kMinPAKEKeyId) | uint64(aPAKEKeyId))
}

func PAKEKeyIdFromNodeId(aNodeId NodeId) PasscodeId {
	return PasscodeId(aNodeId & kMaskPAKEKeyId)
}

func NodeIdFromGroupId(aGroupId GroupId) NodeId {
	return NodeId(uint64(kMinGroupNodeId) | uint64(aGroupId))
}

func (index FabricIndex) IsValidFabricIndex() bool {
	return index >= MinValidFabricIndex && index <= kMaxValidFabricIndex
}
