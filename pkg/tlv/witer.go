package tlv

import (
	"bytes"
	"github.com/galenliu/chip/lib"
	"github.com/galenliu/chip/lib/bitflags"
	"github.com/galenliu/chip/platform/system/buffer"
	"io"
	"math"
)

type Element struct {
	TagControl  TagControl
	ElementType ElementType
	Tag         Tag
}

type WriterBase interface {
	WriteByte(byte2 byte) error
	io.Writer
}

type Writer struct {
	*bytes.Buffer
	mContainerType    TLVType
	mContainerOpen    bool
	mElementType      ElementType
	ImplicitProfileId uint32
}

func NewWriter() *Writer {
	w := &Writer{}
	w.mContainerType = TypeNotSpecified
	w.mContainerOpen = false
	w.Buffer = bytes.NewBuffer(make([]byte, 0))
	return w
}

func (w *Writer) WriteByte(b byte) error {
	data := []byte{b}
	_, err := w.Buffer.Write(data)
	return err
}

func (w *Writer) StartContainer(tag Tag, tlvType TLVType) (outerContainerType TLVType, err error) {
	outerContainerType = w.mContainerType
	if !ElementType(tlvType).IsContainer() {
		return outerContainerType, lib.MATTER_ERROR_WRONG_TLV_TYPE
	}
	err = w.WriteElementHead(ElementType(tlvType), tag, 0)
	if err != nil {
		return outerContainerType, err
	}
	w.mContainerType = tlvType
	w.mContainerOpen = false
	return outerContainerType, nil
}

func (w *Writer) EndContainer(outerContainerType TLVType) error {
	if !ElementType(w.mContainerType).IsContainer() {
		return lib.MATTER_ERROR_WRONG_TLV_TYPE
	}
	w.mContainerType = outerContainerType
	return w.WriteElementHead(EndOfContainer, AnonymousTag(), 0)
}

func (w *Writer) WriteElementWithData(tlvType TLVType, tag Tag, data []byte) error {

	dataLen := len(data)
	var lenFieldSize FieldSize
	if dataLen <= math.MaxUint8 {
		lenFieldSize = FieldSize1Byte
	} else if dataLen <= math.MaxUint16 {
		lenFieldSize = FieldSize2Byte
	} else if dataLen <= math.MaxUint32 {
		lenFieldSize = FieldSize4Byte
	}
	err := w.WriteElementHead(tlvType.WithFieldSize(lenFieldSize), tag, uint64(dataLen))
	if err != nil {
		return err
	}
	return w.write(data)
}

func (w *Writer) WriteElementHead(elementType ElementType, tag Tag, lenOrVal uint64) (err error) {
	buf := bytes.NewBuffer(make([]byte, 0))
	if w.mContainerOpen {
		return lib.MATTER_ERROR_TLV_CONTAINER_OPEN
	}
	if tag.IsContextSpecial() {
		if tag.TagNumber() <= kContextTagMaxNum {
			if w.mContainerType != TypeStructure && w.mContainerType != TypeList {
				return lib.MATTER_ERROR_INVALID_TLV_TAG
			}
			_, err = buf.Write([]byte{elementType.WithTagControl(ContextSpecific)})
			if err != nil {
				return
			}
			_, err = buf.Write([]byte{uint8(tag.TagNumber())})
		} else {
			if elementType != EndOfContainer &&
				w.mContainerType != TypeNotSpecified &&
				w.mContainerType != TypeList &&
				w.mContainerType != TypeArray {
				return lib.MATTER_ERROR_INVALID_TLV_TAG
			}
			buf.Write([]byte{elementType.WithTagControl(Anonymous)})
		}
	} else {
		if w.mContainerType != TypeNotSpecified && w.mContainerType != TypeStructure && w.mContainerType != TypeList {
			return lib.MATTER_ERROR_INVALID_TLV_TAG
		}
		switch tag.ProfileNumber() {
		case uint32(kProfileCommonId):
			if tag.TagNumber() <= math.MaxUint16 {
				err = buffer.Write8(buf, uint8(CommonProfile2Bytes)|uint8(elementType))
				if err != nil {
					return
				}
				err = buffer.LittleEndianWrite16(buf, uint16(tag.TagNumber()))
				if err != nil {
					return
				}
			} else {
				err := buffer.Write8(buf, uint8(CommonProfile4Bytes)|uint8(elementType))
				if err != nil {
					return err
				}
				err = buffer.LittleEndianWrite32(buf, tag.TagNumber())
				if err != nil {
					return err
				}
			}
		case w.ImplicitProfileId:
			if tag.TagNumber() <= math.MaxUint16 {
				err := buffer.Write8(buf, uint8(ImplicitProfile2Bytes)|uint8(elementType))
				if err != nil {
					return err
				}
				err = buffer.LittleEndianWrite16(buf, uint16(tag.TagNumber()))
				if err != nil {
					return err
				}
			} else {
				err := buffer.Write8(buf, uint8(ImplicitProfile4Bytes)|uint8(elementType))
				if err != nil {
					return err
				}
				err = buffer.LittleEndianWrite32(buf, tag.TagNumber())
				if err != nil {
					return err
				}
			}
		default:
			vendorId, profileNum := bitflags.U32To16(tag.ProfileNumber())
			if tag.TagNumber() <= math.MaxUint16 {
				err := buffer.Write8(buf, elementType.WithTagControl(FullyQualified6Bytes))
				if err != nil {
					return err
				}
				err = buffer.LittleEndianWrite16(buf, vendorId)
				if err != nil {
					return err
				}
				err = buffer.LittleEndianWrite16(buf, profileNum)
				if err != nil {
					return err
				}
				err = buffer.LittleEndianWrite16(buf, uint16(tag.TagNumber()))
			} else {
				err := buffer.Write8(buf, elementType.WithTagControl(FullyQualified8Bytes))
				if err != nil {
					return err
				}
				err = buffer.LittleEndianWrite16(buf, vendorId)
				if err != nil {
					return err
				}
				err = buffer.LittleEndianWrite16(buf, profileNum)
				if err != nil {
					return err
				}
				err = buffer.LittleEndianWrite32(buf, tag.TagNumber())
				if err != nil {
					return err
				}
			}

		}
	}

	switch elementType.FieldSize() {
	case FieldSize0Byte:
		break
	case FieldSize1Byte:
		err := buffer.Write8(buf, uint8(lenOrVal))
		if err != nil {
			return err
		}
	case FieldSize2Byte:
		err := buffer.LittleEndianWrite16(buf, uint16(lenOrVal))
		if err != nil {
			if err != nil {
				return err
			}
		}
	case FieldSize4Byte:
		err := buffer.LittleEndianWrite32(buf, uint32(lenOrVal))
		if err != nil {
			if err != nil {
				return err
			}
		}
	case FieldSize8Byte:
		err := buffer.LittleEndianWrite64(buf, lenOrVal)
		if err != nil {
			if err != nil {
				return err
			}
		}
	}
	return w.write(buf.Bytes())
}

func (w *Writer) write(data []byte) error {
	_, err := w.Buffer.Write(data)
	return err
}

func (w *Writer) PutBytes(tag Tag, data []byte) error {
	return w.WriteElementWithData(TypeByteString, tag, data)
}

func (w *Writer) PutUint(tag Tag, val uint64) error {
	if val <= math.MaxUint8 {
		return w.PutU8(tag, uint8(val))
	} else if val <= math.MaxUint16 {
		return w.PutU16(tag, uint16(val))
	} else if val <= math.MaxUint32 {
		return w.PutU32(tag, uint32(val))
	} else {
		return w.PutU64(tag, val)
	}
}

func (w *Writer) PutInt(tag Tag, val int64) error {
	if math.MinInt8 <= val && val <= math.MaxInt8 {
		return w.PutI8(tag, int8(val))
	} else if math.MinInt16 <= val && val <= math.MaxInt16 {
		return w.PutI16(tag, int16(val))
	} else if math.MinInt32 <= val && val <= math.MaxInt32 {
		return w.PutI32(tag, int32(val))
	} else {
		return w.PutI64(tag, val)
	}
}

func (w *Writer) PutU8(tag Tag, val uint8) error {
	return w.WriteElementHead(UInt8, tag, uint64(val))
}

func (w *Writer) PutU16(tag Tag, val uint16) error {
	return w.WriteElementHead(UInt16, tag, uint64(val))
}

func (w *Writer) PutU32(tag Tag, val uint32) error {
	return w.WriteElementHead(UInt32, tag, uint64(val))
}

func (w *Writer) PutU64(tag Tag, val uint64) error {
	return w.WriteElementHead(UInt64, tag, val)
}

func (w *Writer) PutI8(tag Tag, val int8) error {
	return w.WriteElementHead(Int8, tag, uint64(val))
}

func (w *Writer) PutI16(tag Tag, val int16) error {
	return w.WriteElementHead(Int16, tag, uint64(val))
}

func (w *Writer) PutI32(tag Tag, val int32) error {
	return w.WriteElementHead(Int32, tag, uint64(val))
}

func (w *Writer) PutI64(tag Tag, val int64) error {
	return w.WriteElementHead(Int64, tag, uint64(val))
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
