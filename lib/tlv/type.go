package tlv

type Type int8
type fieldSize int8
type elementType int8

const (
	TypeNotSpecified        Type = -1
	TypeUnknownContainer    Type = -2
	TypeSignedInteger       Type = 0x00 //• Signed integers
	TypeUnsignedInteger     Type = 0x04 //• Unsigned integers
	TypeUTF8String          Type = 0x0C //• UTF-8 Strings
	TypeByteString          Type = 0x10 //• Octet Strings
	TypeFloatingPointNumber Type = 0x0A //• Single or double-precision floating point numbers (following IEEE 754-2019) • Booleans
	TypeBoolean             Type = 0x08 //• Booleans
	TypeNull                Type = 0x14 //• Nulls
	TypeStructure           Type = 0x15
	TypeArray               Type = 0x16
	TypeList                Type = 0x17
)

const (
	notSpecified          elementType = -1
	i8                    elementType = 0x00
	i16                   elementType = 0x01
	i32                   elementType = 0x02
	i64                   elementType = 0x03
	u8                    elementType = 0x04
	u16                   elementType = 0x05
	u32                   elementType = 0x06
	u64                   elementType = 0x07
	booleanFalse          elementType = 0x08
	booleanTrue           elementType = 0x09
	floatingPointNumber32 elementType = 0x0A
	floatingPointNumber64 elementType = 0x0B
	utf8String1ByteLength elementType = 0x0C
	utf8String2ByteLength elementType = 0x0D
	utf8String4ByteLength elementType = 0x0E
	utf8String8ByteLength elementType = 0x0F
	byteString1ByteLength elementType = 0x10
	byteString2ByteLength elementType = 0x11
	byteString4ByteLength elementType = 0x12
	byteString8ByteLength elementType = 0x13
	null                  elementType = 0x14
	structure             elementType = 0x15
	array                 elementType = 0x16
	list                  elementType = 0x17
	endOfContainer        elementType = 0x18
)

const (
	fieldSize0Byte fieldSize = -1
	fieldSize1Byte fieldSize = 0
	fieldSize2Byte fieldSize = 1
	fieldSize4Byte fieldSize = 2
	fieldSize8Byte fieldSize = 3
)

const (
	fTypeMask     uint8 = 0x1F
	fTypeSizeMask uint8 = 0x03
)

func (f fieldSize) byteSize() uint8 {
	if f == fieldSize0Byte {
		return 0
	}
	return uint8(f) << 1
}

func elemType[T ~uint8 | ~uint16](val T) elementType {
	return elementType(uint8(val) & fTypeMask)
}

func (t elementType) withFieldSize(size fieldSize) elementType {
	if size == fieldSize0Byte {
		return t
	}
	return elementType(uint8(t) | uint8(size))
}

func (t elementType) isContainer() bool {
	return t <= list && t >= structure
}

func (t Type) isContainer() bool {
	return elementType(t) <= list && elementType(t) >= structure
}

func (t elementType) isValid() bool {
	return t <= endOfContainer
}

func (t elementType) hasValue() bool {
	return t <= u64 || (t >= floatingPointNumber32 && t <= byteString8ByteLength)
}

func (t elementType) hasLength() bool {
	return t >= utf8String1ByteLength && t <= byteString8ByteLength
}

func (t elementType) isString() bool {
	return t >= utf8String1ByteLength && t <= byteString8ByteLength
}

func (t elementType) isUTF8String() bool {
	return t >= utf8String1ByteLength && t <= utf8String8ByteLength
}

func (t elementType) isByteString() bool {
	return t >= byteString1ByteLength && t <= byteString8ByteLength
}

func (t elementType) fieldSize() fieldSize {
	if t.hasValue() {
		return fieldSize(uint8(t) & fTypeSizeMask)
	}
	return fieldSize0Byte
}

func (t Type) String() string {
	switch t {
	case TypeNotSpecified:
		return "notSpecified"
	case TypeUnknownContainer:
		return "Unknown Container"
	case TypeSignedInteger:
		return "Signed Integer"
	case TypeUnsignedInteger:
		return "Unsigned Integer"
	case TypeUTF8String:
		return "UTF8 String"
	case TypeByteString:
		return "Byte String"
	case TypeFloatingPointNumber:
		return "Floating Point Number"
	case TypeBoolean:
		return "Boolean"
	case TypeNull:
		return "Null"
	case TypeStructure:
		return "Structure"
	case TypeArray:
		return "Array"
	case TypeList:
		return "List"
	default:
		return "Unknown"
	}
}

func (t elementType) String() string {
	index := int(t)
	if index < 0 || index > 0x18 {
		return "notSpecified"
	}
	types := []string{
		"i8",
		"i16",
		"i32",
		"i64",
		"u8 ",
		"u16",
		"u32",
		"u64",
		"booleanFalse",
		"booleanTrue",
		"floatingPointNumber32",
		"floatingPointNumber64",
		"utf8String1ByteLength",
		"utf8String2ByteLength",
		"utf8String4ByteLength",
		"utf8String8ByteLength",
		"byteString1ByteLength",
		"byteString2ByteLength",
		"byteString4ByteLength",
		"byteString8ByteLength",
		"null",
		"structure",
		"array",
		"list",
		"endOfContainer",
	}
	return types[index]
}
