package tlv

type TagFields uint16
type FieldSize int8
type TagControl uint8
type TLVType int8
type ElementType int8

const (
	TypeNull                 TLVType = 0x14
	Type_Structure           TLVType = 0x15
	Type_Array               TLVType = 0x16
	Type_List                TLVType = 0x17
	Type_NotSpecified        TLVType = -1
	TypeUnknownContainer     TLVType = -2
	TypeSignedInteger        TLVType = 0x00
	Type_UnsignedInteger     TLVType = 0x04
	TypeBoolean              TLVType = 0x08
	Type_FloatingPointNumber TLVType = 0x0A
	Type_UTF8String          TLVType = 0x0C
	Type_ByteString          TLVType = 0x10
)

// IMPORTANT: All values here must have no bits in common with specified
// values of TLVElementType.
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

func controlByte(control TagControl, elementType ElementType) byte {
	return byte(control) & byte(elementType)
}

func (t TagControl) Uint8() uint8 {
	return uint8(t)
}

const (
	kTLVFieldSize0Byte FieldSize = -1
	kTLVFieldSize1Byte FieldSize = 0
	kTLVFieldSize2Byte FieldSize = 1
	kTLVFieldSize4Byte FieldSize = 2
	kTLVFieldSize8Byte FieldSize = 3
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

func IsContainerType(elementType ElementType) bool {
	return elementType <= List && elementType <= Structure
}

func (t ElementType) HasValue() bool {
	return t <= UInt64 || (t >= FloatingPointNumber32 && t <= ByteString8ByteLength)
}

func (t ElementType) Uint8() uint8 {
	return uint8(t)
}

func WriterControl(buffer Writer, control TagControl, elementType ElementType) error {
	data := control.Uint8() & elementType.Uint8()
	err := buffer.WriteByte(data)
	return err
}

//func WriterUin16(buffer Writer, val uint16) error {
//	data := make([]byte,2)
//	binary.LittleEndian.PutUint16(data,val)
//	_,err:= buffer.Write(data)
//	return err
//}
//
//func WriterUin16(buffer Writer, val uint16) error {
//	data := make([]byte,2)
//	binary.LittleEndian.PutUint16(data,val)
//	_,err:= buffer.Write(data)
//	return err
//
//
//}
