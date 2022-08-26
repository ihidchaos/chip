package tlv

const (
	fTLVTypeMask      uint8 = 0x1F
	fTLVTypeSizeMask  uint8 = 0x03
	fProfileIdMask          = 0xFFFFFFFF00000000
	fProfileNumMask         = 0x0000FFFF00000000
	fVendorIdMask           = 0xFFFF000000000000
	fProfileIdShift         = 32
	fVendorIdShift          = 48
	fProfileNumShift        = 32
	fTagNumMask             = 0x00000000FFFFFFFF
	fSpecialTagMarker       = 0xFFFFFFFF00000000
	fContextTagMaxNum uint8 = 0xFF
)

const UnknownImplicitTag = fSpecialTagMarker | 0x00000000FFFFFFFE

type TLVCommonProfiles uint64

/**
 * Used to indicate the absence of a profile id in a variable or member.
 * This is essentially the same as kCHIPProfile_NotSpecified defined in CHIPProfiles.h
 */

type Tag uint64

func (t Tag) Equal(tag Tag) bool {
	return t == tag
}

func (t Tag) IsSpecialTag() bool {
	return t&0xFFFFFFFF00000000 == 0xFFFFFFFF00000000
}

func AnonymousTag() Tag {
	return Tag(0xFFFFFFFF00000000 | 0x0000000FFFFFFFFF)
}

func ContextTag(tagNum uint8) Tag {
	return Tag(0xFFFFFFFF00000000 | uint64(tagNum))
}

func CommonTag4Byte(val uint32) Tag {
	return ProfileTag(0x0, val)
}

func CommonTag2Byte(val uint16) Tag {
	return ProfileTag(0x0, uint32(val))
}

func ProfileTag(profileId, tagNum uint32) Tag {
	return Tag((uint64(profileId))<<32 | uint64(tagNum))
}

func ProfileTag4Byte(vendorId uint16, profileNum uint16, tagNum uint32) Tag {
	tag := uint64(vendorId)<<48 | uint64(profileNum)<<32 | uint64(tagNum)
	return Tag(tag)
}

func GetTLVFieldSize(elementType ElementType) FieldSize {
	if elementType.HasValue() {
		return FieldSize(uint8(elementType) & 0x03)
	}
	return kTLVFieldSize0Byte
}
