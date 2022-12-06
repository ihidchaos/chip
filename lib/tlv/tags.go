package tlv

import "math"

type tagControl uint8
type Tag uint64

const (
	Anonymous             tagControl = 0x00
	ContextSpecific       tagControl = 0x20
	CommonProfile2Bytes   tagControl = 0x40
	CommonProfile4Bytes   tagControl = 0x60
	ImplicitProfile2Bytes tagControl = 0x80
	ImplicitProfile4Bytes tagControl = 0xA0
	FullyQualified6Bytes  tagControl = 0xC0
	FullyQualified8Bytes  tagControl = 0xE0
)

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
	kCommonProfileId  = 0
)

func tagCtl[T ~uint8 | ~uint16](val T) tagControl {
	const fTagControlMask = 0xE0
	return tagControl(uint8(val) & fTagControlMask)
}

func (tc tagControl) bytesSize() uint8 {
	const fTagControlShift = 5
	return []uint8{0, 1, 2, 4, 2, 4, 6, 8}[uint(tc>>fTagControlShift)]
}

func (tc tagControl) withElemType(et elementType) uint8 {
	return uint8(tc) | uint8(et)
}

const unknownImplicitTag = fSpecialTagMarker | 0x00000000FFFFFFFE

type commonProfilesU32 uint32

const (
	profileIdNotSpecified commonProfilesU32 = 0xFFFFFFFF
	profileCommonId       commonProfilesU32 = 0
)

/**
 * Used to indicate the absence of a profile id in a variable or member.
 * This is essentially the same as kCHIPProfile_NotSpecified defined in CHIPProfiles.h
 */

func AnonymousTag() Tag {
	return Tag(0xFFFFFFFF00000000 | 0x00000000FFFFFFFF)
}

func ContextTag(tagNum uint8) Tag {
	return Tag(0xFFFFFFFF00000000 | uint64(tagNum))
}

func commonTag[T ~uint16 | ~uint32](val T) Tag {
	return profileTag(kCommonProfileId, val)
}

func profileTag[T ~uint16 | ~uint32](profileId uint32, tagNum T) Tag {
	return Tag((uint64(profileId))<<32 | uint64(tagNum))
}

func profileSpecificTag[T uint16 | uint32](vendorId uint16, profileNum uint16, tagNum T) Tag {
	return profileTag(uint32(vendorId)<<16|uint32(profileNum), tagNum)
}

func (t Tag) vendorId() uint16 {
	value := (t & fVendorIdMask) >> fVendorIdShift
	return uint16(value)
}

func (t Tag) profileNumber() uint16 {
	return uint16((uint64(t) & fProfileNumMask) >> fProfileIdShift)
}

func (t Tag) number() uint32 {
	return uint32(t & fTagNumMask)
}

func (t Tag) isSpecial() bool {
	return (uint64(t) & fProfileIdMask) == fSpecialTagMarker
}

func (t Tag) isContext() bool {
	return t.isSpecial() && t.number() <= kContextTagMaxNum
}
