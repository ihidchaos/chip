package crypto

import (
	"bytes"
	"github.com/galenliu/chip/lib"
)

type TLVElementType int

const (

	// IMPORTANT: All values here except NotSpecified must have no bits in
	// common with values of TagControl.
	NotSpecified          TLVElementType = -1
	Int8                  TLVElementType = 0x00
	Int16                 TLVElementType = 0x01
	Int32                 TLVElementType = 0x02
	Int64                 TLVElementType = 0x03
	UInt8                 TLVElementType = 0x04
	UInt16                TLVElementType = 0x05
	UInt32                TLVElementType = 0x06
	UInt64                TLVElementType = 0x07
	BooleanFalse          TLVElementType = 0x08
	BooleanTrue           TLVElementType = 0x09
	FloatingPointNumber32 TLVElementType = 0x0A
	FloatingPointNumber64 TLVElementType = 0x0B
	Utf8string1ByteLength TLVElementType = 0x0C
	UTF8String2ByteLength TLVElementType = 0x0D
	UTF8String4ByteLength TLVElementType = 0x0E
	UTF8String8ByteLength TLVElementType = 0x0F
	ByteString1ByteLength TLVElementType = 0x10
	ByteString2ByteLength TLVElementType = 0x11
	ByteString4ByteLength TLVElementType = 0x12
	ByteString8ByteLength TLVElementType = 0x13

	// IMPORTANT: Values starting at Null must match the corresponding values of
	// TLVType.
	Null           TLVElementType = 0x14
	Structure      TLVElementType = 0x15
	Array          TLVElementType = 0x16
	List           TLVElementType = 0x17
	EndOfContainer TLVElementType = 0x18
)

func GetValueLength(buf *bytes.Buffer) (uint64, error) {
	var f uint8 = 0b10000000
	data, err := buf.ReadByte()
	if err != nil {
		return 0, err
	}
	if lib.HasFlags(data, f) {
		size := data & (^f)
		data := make([]byte, size)
		_, err := buf.Read(data)
		if err != nil {
			return 0, err
		}
	}
	return 0, err
}
