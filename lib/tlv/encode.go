package tlv

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"golang.org/x/exp/constraints"
	"io"
	"math"
)

type Encoder struct {
	w                      io.Writer
	err                    error
	containerType          Type
	containerOpen          bool
	closeContainerReserved bool
	elementType            elementType
	implicitProfileId      commonProfilesU32
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{
		w:                      w,
		containerType:          TypeNotSpecified,
		containerOpen:          false,
		closeContainerReserved: true,
		elementType:            notSpecified,
		implicitProfileId:      profileIdNotSpecified,
	}
}

func (enc *Encoder) StartContainer(tag Tag, tlvType Type) (outerContainerType Type, err error) {
	outerContainerType = enc.containerType
	if !tlvType.isContainer() {
		return outerContainerType, enc.elementTypeError(tlvType)
	}
	err = enc.writeElementHead(elementType(tlvType), tag, 0)
	if err != nil {
		return outerContainerType, err
	}
	enc.containerType = tlvType
	enc.containerOpen = false
	return outerContainerType, nil
}

func (enc *Encoder) EndContainer(outerContainerType Type) error {
	if !enc.containerType.isContainer() {
		return enc.elementTypeError(enc.containerType)
	}
	enc.containerType = outerContainerType
	return enc.writeElementHead(endOfContainer, AnonymousTag(), 0)
}

func (enc *Encoder) WriteElementWithData(tlvType Type, tag Tag, data []byte) error {
	dataLen := len(data)
	var lenFieldSize = fieldSize8Byte
	if dataLen <= math.MaxUint8 {
		lenFieldSize = fieldSize1Byte
	} else if dataLen <= math.MaxUint16 {
		lenFieldSize = fieldSize2Byte
	} else if dataLen <= math.MaxUint32 {
		lenFieldSize = fieldSize4Byte
	}
	err := enc.writeElementHead(elementType(tlvType).withFieldSize(lenFieldSize), tag, uint64(dataLen))
	if err != nil {
		return err
	}
	return enc.writeData(data)
}

func (enc *Encoder) writeElementHead(elemType elementType, tag Tag, lenOrVal uint64) (err error) {
	buf := bytes.NewBuffer(make([]byte, 0))
	if enc.containerOpen {
		return enc.incorrectStateError("container open")
	}
	if tag.isSpecial() {
		if tag.isContext() {
			if enc.containerType != TypeStructure &&
				enc.containerType != TypeList {
				return enc.tagError(tag)
			}
			buf.WriteByte(ContextSpecific.withElemType(elemType))
			buf.WriteByte(uint8(tag.Number()))
		} else {
			if elemType != endOfContainer &&
				enc.containerType != TypeNotSpecified &&
				enc.containerType != TypeList &&
				enc.containerType != TypeArray {
				return enc.tagError(tag)
			}
			buf.WriteByte(Anonymous.withElemType(elemType))
		}
	} else {
		if enc.containerType != TypeNotSpecified &&
			enc.containerType != TypeStructure &&
			enc.containerType != TypeList {
			return enc.tagError(tag)
		}
		switch tag.Number() {
		case uint32(profileCommonId):
			if tag.Number() <= math.MaxUint16 {
				buf.WriteByte(CommonProfile2Bytes.withElemType(elemType))
				var data = make([]byte, 2)
				binary.LittleEndian.PutUint16(data, uint16(tag.Number()))
				buf.Write(data)
			} else {
				buf.WriteByte(CommonProfile4Bytes.withElemType(elemType))
				var data = make([]byte, 4)
				binary.LittleEndian.PutUint32(data, tag.Number())
				buf.Write(data)
			}
		case uint32(enc.implicitProfileId):
			if tag.Number() <= math.MaxUint16 {
				buf.WriteByte(ImplicitProfile2Bytes.withElemType(elemType))
				data := make([]byte, 2)
				binary.LittleEndian.PutUint16(data, uint16(tag.Number()))
				buf.Write(data)
			} else {
				buf.WriteByte(ImplicitProfile4Bytes.withElemType(elemType))
				data := make([]byte, 4)
				binary.LittleEndian.PutUint32(data, tag.Number())
				buf.Write(data)
			}
		default:
			if tag.Number() <= math.MaxUint16 {
				buf.WriteByte(FullyQualified6Bytes.withElemType(elemType))
				v := make([]byte, 2)
				binary.LittleEndian.PutUint16(v, tag.vendorId())
				p := make([]byte, 2)
				binary.LittleEndian.PutUint16(p, tag.profileNumber())
				n := make([]byte, 2)
				binary.LittleEndian.PutUint16(p, uint16(tag.Number()))
				buf.Write(v)
				buf.Write(p)
				buf.Write(n)
			} else {
				buf.WriteByte(FullyQualified8Bytes.withElemType(elemType))
				v := make([]byte, 2)
				binary.LittleEndian.PutUint16(v, tag.vendorId())
				p := make([]byte, 2)
				binary.LittleEndian.PutUint16(p, tag.profileNumber())
				n := make([]byte, 4)
				binary.LittleEndian.PutUint32(p, tag.Number())
				buf.Write(v)
				buf.Write(p)
				buf.Write(n)
			}
		}
	}
	switch elemType.fieldSize() {
	case fieldSize0Byte:
		break
	case fieldSize1Byte:
		buf.WriteByte(uint8(lenOrVal))
	case fieldSize2Byte:
		data := make([]byte, 2)
		binary.LittleEndian.PutUint16(data, uint16(lenOrVal))
		buf.Write(data)
	case fieldSize4Byte:
		data := make([]byte, 4)
		binary.LittleEndian.PutUint32(data, uint32(lenOrVal))
		buf.Write(data)
	case fieldSize8Byte:
		data := make([]byte, 8)
		binary.LittleEndian.PutUint64(data, lenOrVal)
		buf.Write(data)
	}
	return enc.writeData(buf.Bytes())
}

func (enc *Encoder) writeData(data []byte) error {
	n, err := enc.w.Write(data)
	if n != len(data) {
		err = fmt.Errorf("writeData err:%w", err)
	}
	return err
}

func (enc *Encoder) PutBytes(tag Tag, data []byte) error {
	return enc.WriteElementWithData(TypeByteString, tag, data)
}

func (enc *Encoder) PutUTFString(tag Tag, data []byte) error {
	return enc.WriteElementWithData(TypeUTF8String, tag, data)
}

func (enc *Encoder) putUint(tag Tag, val uint64) error {
	if val <= math.MaxUint8 {
		return enc.PutU8(tag, uint8(val))
	} else if val <= math.MaxUint16 {
		return enc.PutU16(tag, uint16(val))
	} else if val <= math.MaxUint32 {
		return enc.PutU32(tag, uint32(val))
	} else {
		return enc.PutU64(tag, val)
	}
}

func (enc *Encoder) putInt(tag Tag, val int64) error {
	if math.MinInt8 <= val && val <= math.MaxInt8 {
		return enc.PutI8(tag, int8(val))
	} else if math.MinInt16 <= val && val <= math.MaxInt16 {
		return enc.PutI16(tag, int16(val))
	} else if math.MinInt32 <= val && val <= math.MaxInt32 {
		return enc.PutI32(tag, int32(val))
	} else {
		return enc.PutI64(tag, val)
	}
}

func (enc *Encoder) PutU8(tag Tag, val uint8) error {
	return enc.writeElementHead(u8, tag, uint64(val))
}

func (enc *Encoder) PutU16(tag Tag, val uint16) error {
	return enc.writeElementHead(u16, tag, uint64(val))
}

func (enc *Encoder) PutU32(tag Tag, val uint32) error {
	return enc.writeElementHead(u32, tag, uint64(val))
}

func (enc *Encoder) PutU64(tag Tag, val uint64) error {
	return enc.writeElementHead(u64, tag, val)
}

func (enc *Encoder) PutI8(tag Tag, val int8) error {
	return enc.writeElementHead(i8, tag, uint64(val))
}

func (enc *Encoder) PutI16(tag Tag, val int16) error {
	return enc.writeElementHead(i16, tag, uint64(val))
}

func (enc *Encoder) PutI32(tag Tag, val int32) error {
	return enc.writeElementHead(i32, tag, uint64(val))
}

func (enc *Encoder) PutI64(tag Tag, val int64) error {
	return enc.writeElementHead(i64, tag, uint64(val))
}

func (enc *Encoder) PutF32(tag Tag, f32 float32) error {
	return enc.writeElementHead(floatingPointNumber32, tag, uint64(f32))
}

func (enc *Encoder) PutF64(tag Tag, f64 float64) error {
	return enc.writeElementHead(floatingPointNumber64, tag, uint64(f64))
}

func (enc *Encoder) PutBoolean(tag Tag, b bool) error {
	if b {
		return enc.writeElementHead(booleanTrue, tag, 0)
	}
	return enc.writeElementHead(booleanFalse, tag, 0)
}

func (enc *Encoder) Status() error {
	return enc.err
}

func (enc *Encoder) incorrectStateError(str string) error {
	return fmt.Errorf("incorrect state err:%s", str)
}

func (enc *Encoder) elementTypeError(val any) error {
	return fmt.Errorf("element type err:%s", val)
}

func (enc *Encoder) tagError(val any) error {
	return fmt.Errorf("tag err:%s", val)
}

func PutUint[T constraints.Unsigned](enc *Encoder, tag Tag, val T) error {
	return enc.putUint(tag, uint64(val))
}

func PutInt[T constraints.Signed](enc *Encoder, tag Tag, val T) error {
	return enc.putInt(tag, int64(val))
}

func EstimateStructOverhead(firstFieldSize int, args ...int) int {
	size := firstFieldSize + 4
	if len(args) > 0 {
		for _, arg := range args {
			size = arg + 4
		}
	}
	return size
}
