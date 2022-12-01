package tlv

type TagFields uint16
type FieldSize int8
type TLVType int8
type ElementType int8

const (
	TypeNotSpecified     TLVType = -1
	TypeUnknownContainer TLVType = -2

	TypeSignedInteger       TLVType = 0x00 //• Signed integers
	TypeUnsignedInteger     TLVType = 0x04 //• Unsigned integers
	TypeUTF8String          TLVType = 0x0C //• UTF-8 Strings
	TypeByteString          TLVType = 0x10 //• Octet Strings
	TypeFloatingPointNumber TLVType = 0x0A //• Single or double-precision floating point numbers (following IEEE 754-2019) • Booleans
	TypeBoolean             TLVType = 0x08 //• Booleans
	TypeNull                TLVType = 0x14 //• Nulls

	TypeStructure TLVType = 0x15
	TypeArray     TLVType = 0x16
	TypeList      TLVType = 0x17
)

func (t TLVType) WithFieldSize(size FieldSize) ElementType {
	if size == FieldSize0Byte {
		return ElementType(t)
	}
	return ElementType(uint8(t) | uint8(size))
}

const (
	FieldSize0Byte FieldSize = -1
	FieldSize1Byte FieldSize = 0
	FieldSize2Byte FieldSize = 1
	FieldSize4Byte FieldSize = 2
	FieldSize8Byte FieldSize = 3
)

func (f FieldSize) ByteSize() uint8 {
	if f == FieldSize0Byte {
		return 0
	}
	return uint8(f) << 1
}

const (
	fTypeMask     uint8 = 0x1F
	fTypeSizeMask uint8 = 0x03
)

const (
	NotSpecified          ElementType = -1
	Int8                  ElementType = 0x00
	Int16                 ElementType = 0x01
	Int32                 ElementType = 0x02
	Int64                 ElementType = 0x03
	UInt8                 ElementType = 0x04
	UInt16                ElementType = 0x05
	UInt32                ElementType = 0x06
	UInt64                ElementType = 0x07
	BooleanFalse          ElementType = 0x08
	BooleanTrue           ElementType = 0x09
	FloatingPointNumber32 ElementType = 0x0A
	FloatingPointNumber64 ElementType = 0x0B
	UTF8String1ByteLength ElementType = 0x0C
	UTF8String2ByteLength ElementType = 0x0D
	UTF8String4ByteLength ElementType = 0x0E
	UTF8String8ByteLength ElementType = 0x0F
	ByteString1ByteLength ElementType = 0x10
	ByteString2ByteLength ElementType = 0x11
	ByteString4ByteLength ElementType = 0x12
	ByteString8ByteLength ElementType = 0x13

	Null           ElementType = 0x14
	Structure      ElementType = 0x15
	Array          ElementType = 0x16
	List           ElementType = 0x17
	EndOfContainer ElementType = 0x18
)

func NewElementType[T ~uint8 | ~uint16](val T) ElementType {
	return ElementType(uint8(val) & fTypeMask)
}

func (t ElementType) IsContainer() bool {
	return t <= List && t >= Structure
}

func (t ElementType) WithTagControl(tag TagControl) uint8 {
	return uint8(t) | uint8(tag)
}

func (t ElementType) IsValid() bool {
	return t <= EndOfContainer
}

func (t ElementType) HasValue() bool {
	return t <= UInt64 || (t >= FloatingPointNumber32 && t <= ByteString8ByteLength)
}

func (t ElementType) HasLength() bool {
	return t >= UTF8String1ByteLength && t <= ByteString8ByteLength
}

func (t ElementType) IsString() bool {
	return t >= UTF8String1ByteLength && t <= ByteString8ByteLength
}

func (t ElementType) IsUTF8String() bool {
	return t >= UTF8String1ByteLength && t <= UTF8String8ByteLength
}

func (t ElementType) IsByteString() bool {
	return t >= ByteString1ByteLength && t <= ByteString8ByteLength
}

func (t ElementType) FieldSize() FieldSize {
	if t.HasValue() {
		return FieldSize(uint8(t) & fTypeSizeMask)
	}
	return FieldSize0Byte
}

//func WriterControl(buffer WriterBase, control TagControl, elementType ElementType) error {
//	data := control.Uint8() & uint8(elementType)
//	err := buffer.WriteByte(data)
//	return err
//}
