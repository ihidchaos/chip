package lib

type VendorId uint16
type ProductId uint16

const ProductIdAnonymous ProductId = 0x0000

const (
	UndefinedFabricIndex FabricIndex = 0x0
	MinValidFabricIndex  FabricIndex = 0x1
	kMaxValidFabricIndex FabricIndex = 0xFE
)

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

func NodeIdFromPAKEKeyId(aPAKEKeyId PasscodeId) NodeId {
	return NodeId(uint64(minPAKEKeyId) | uint64(aPAKEKeyId))
}

func PAKEKeyIdFromNodeId(aNodeId NodeId) PasscodeId {
	return PasscodeId(aNodeId & maskPAKEKeyId)
}

func NodeIdFromGroupId(aGroupId GroupId) NodeId {
	return NodeId(uint64(minGroupNodeId) | uint64(aGroupId))
}

func (index FabricIndex) IsValidFabricIndex() bool {
	return index >= MinValidFabricIndex && index <= kMaxValidFabricIndex
}
