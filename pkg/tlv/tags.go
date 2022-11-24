package tlv

import "math"

type TagControl uint8
type Tag uint64

const (
	Anonymous             TagControl = 0x00
	ContextSpecific       TagControl = 0x20
	CommonProfile2Bytes   TagControl = 0x40
	CommonProfile4Bytes   TagControl = 0x60
	ImplicitProfile2Bytes TagControl = 0x80
	ImplicitProfile4Bytes TagControl = 0xA0
	FullyQualified6Bytes  TagControl = 0xC0
	FullyQualified8Bytes  TagControl = 0xE0
)

func (tc TagControl) WithElementType(et ElementType) uint8 {
	return uint8(tc) | uint8(et)
}

const (
	fTLVTypeMask      = 0x1F
	fTLVTypeSizeMask  = 0x03
	fProfileIdMask    = 0xFFFFFFFF00000000
	fProfileNumMask   = 0x0000FFFF00000000
	fVendorIdMask     = 0xFFFF000000000000
	fProfileIdShift   = 32
	fVendorIdShift    = 48
	fProfileNumShift  = 32
	fTagNumMask       = 0x00000000FFFFFFFF
	fSpecialTagMarker = 0xFFFFFFFF00000000
	kContextTagMaxNum = math.MaxUint8
)

const UnknownImplicitTag = fSpecialTagMarker | 0x00000000FFFFFFFE

type CommonProfiles uint32

const (
	kProfileNotSpecified CommonProfiles = 0xFFFFFFFF
	kProfileCommon       CommonProfiles = 0
)

/**
 * Used to indicate the absence of a profile id in a variable or member.
 * This is essentially the same as kCHIPProfile_NotSpecified defined in CHIPProfiles.h
 */

func AnonymousTag() Tag {
	return Tag(0xFFFFFFFF00000000 | 0x00000000FFFFFFFF)
}

func ContextSpecificTag(tagNum uint8) Tag {
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

func ProfileSpecificTag(vendorId uint16, profileNum uint16, tagNum uint32) Tag {
	tag := uint64(vendorId)<<48 | uint64(profileNum)<<32 | uint64(tagNum)
	return Tag(tag)
}

func (t Tag) Equal(tag Tag) bool {
	return t == tag
}

func (t Tag) VendorId() uint16 {
	value := (t & fVendorIdMask) >> fVendorIdShift
	return uint16(value)
}

func (t Tag) ProfileNumber() uint32 {
	return uint32((uint64(t) & fProfileNumMask) >> fProfileIdShift)
}

func (t Tag) TagNumber() uint32 {
	return uint32(t & fTagNumMask)
}

func (t Tag) IsContextSpecial() bool {
	return (uint64(t) & fProfileIdMask) == fSpecialTagMarker
}
