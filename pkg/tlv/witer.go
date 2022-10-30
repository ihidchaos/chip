package tlv

import (
	"bytes"
	"fmt"
	"github.com/galenliu/chip/lib"
	"github.com/galenliu/chip/pkg/tlv/buffer"
	"io"
	"math"
)

type Element struct {
	TagControl  TagControl
	ElementType ElementType
	Tag         Tag
}

type Writer interface {
	WriteByte(byte2 byte) error
	io.Writer
}

type WriterImpl struct {
	*bytes.Buffer
	mContainerType    TLVType
	mContainerOpen    bool
	mValueOrLength    uint64
	mElementType      ElementType
	ImplicitProfileId uint32
}

func (w *WriterImpl) WriteByte(b byte) error {
	data := []byte{b}
	_, err := w.Buffer.Write(data)
	return err
}

func NewWriter() *WriterImpl {
	w := &WriterImpl{}
	w.Buffer = bytes.NewBuffer(make([]byte, 0))
	return w
}

func (w *WriterImpl) StartContainer(tag Tag, tlvType TLVType) error {
	if !IsContainerType(ElementType(tlvType)) {
		return fmt.Errorf("container type not supported")
	}
	w.mContainerType = tlvType
	w.mContainerOpen = false
	return w.WriteElement(ElementType(tlvType), tag, 0)

}

func (w *WriterImpl) WriteElement(elemType ElementType, tag Tag, lengthOrValue uint64) error {

	tagNum := uint32(tag)
	if tag.IsSpecialTag() {
		if tagNum <= 0xFF {
			if w.mContainerType != Type_Structure && w.mContainerType != Type_List {
				return lib.MatterErrorInvalidTlvTag
			}
			w.Buffer.WriteByte(controlByte(ContextSpecific, elemType))
			w.Buffer.WriteByte(byte(tagNum))
		} else {
			if elemType != EndOfContainer &&
				w.mContainerType != Type_NotSpecified &&
				w.mContainerType != Type_Array &&
				w.mContainerType != Type_List {
				return lib.MatterErrorInvalidTlvTag
			}
			w.Buffer.WriteByte(controlByte(Anonymous, elemType))
		}
	} else {
		profileId := uint32((uint64(tag) & 0xFFFFFFFF00000000) >> 32)
		if profileId == 0 {
			if tagNum < 65536 {
				w.Buffer.WriteByte(controlByte(CommonProfile2Bytes, elemType))
				err := buffer.LittleEndianWrite16(w, uint16(tagNum))
				if err != nil {
					return err
				}
			} else {
				w.Buffer.WriteByte(controlByte(CommonProfile4Bytes, elemType))
				err := buffer.LittleEndianWrite32(w, tagNum)
				if err != nil {
					return err
				}
			}
		} else if profileId == w.ImplicitProfileId {
			if tagNum < 65536 {
				w.Buffer.WriteByte(controlByte(ImplicitProfile2Bytes, elemType))
				err := buffer.LittleEndianWrite16(w, uint16(tagNum))
				if err != nil {
					return err
				}
			} else {
				w.Buffer.WriteByte(controlByte(ImplicitProfile4Bytes, elemType))
				err := buffer.LittleEndianWrite32(w, tagNum)
				if err != nil {
					return err
				}
			}
		} else {
			vendorId := uint16(profileId >> 16)
			profileNum := uint16(profileId)
			if tagNum < 65536 {
				w.Buffer.WriteByte(controlByte(FullyQualified6Bytes, elemType))
				err := buffer.LittleEndianWrite16(w, vendorId)
				if err != nil {
					return err
				}
				err = buffer.LittleEndianWrite16(w, profileNum)
				if err != nil {
					return err
				}
				err = buffer.LittleEndianWrite16(w, uint16(tagNum))
				if err != nil {
					return err
				}
			} else {
				w.Buffer.WriteByte(controlByte(FullyQualified8Bytes, elemType))
				err := buffer.LittleEndianWrite16(w, vendorId)
				if err != nil {
					return err
				}
				err = buffer.LittleEndianWrite16(w, profileNum)
				if err != nil {
					return err
				}
				err = buffer.LittleEndianWrite32(w, tagNum)
				if err != nil {
					return err
				}
			}

		}
	}

	switch GetTLVFieldSize(elemType) {
	case kTLVFieldSize1Byte:
		err := w.WriteByte(byte(lengthOrValue))
		if err != nil {
			return err
		}
	case kTLVFieldSize2Byte:
		err := buffer.LittleEndianWrite16(w, uint16(lengthOrValue))
		if err != nil {
			return err
		}
	case kTLVFieldSize4Byte:

		err := buffer.LittleEndianWrite32(w, uint32(lengthOrValue))
		if err != nil {
			return err
		}
	case kTLVFieldSize8Byte:
		err := buffer.LittleEndianWrite64(w, lengthOrValue)
		if err != nil {
			return err
		}
	}
	return nil
}

func (w *WriterImpl) WriteData(data []byte) error {
	if w.mElementType < 0x0C && w.mElementType > 0x13 {
		return fmt.Errorf("element type must be 0x0C-0x13, got %d", w.mElementType)
	}
	i, err := w.Buffer.Write(data)
	w.mValueOrLength = w.mValueOrLength - uint64(i)
	return err
}

func (w *WriterImpl) PutBytes(tag Tag, data []byte) error {
	return w.WriteElementWithData(Type_ByteString, tag, data)
}

func (w *WriterImpl) WriteElementWithData(tlvType TLVType, tag Tag, data []byte) error {
	dataLen := len(data)
	var lenFieldSize FieldSize
	if dataLen <= 0xff {
		lenFieldSize = kTLVFieldSize1Byte
	} else if dataLen <= 0xffff {
		lenFieldSize = kTLVFieldSize2Byte
	} else if dataLen <= 0xffffffff {
		lenFieldSize = kTLVFieldSize4Byte
	}
	err := w.WriteElement(ElementType(byte(tlvType)|byte(lenFieldSize)), tag, uint64(dataLen))
	if err != nil {
		return err
	}
	return w.WriteData(data)
}

func (w *WriterImpl) EndContainer(t TLVType) error {
	w.mContainerType = t
	return w.WriteElement(EndOfContainer, AnonymousTag(), 0)
}

func (w *WriterImpl) Put(tag Tag, val uint64) error {
	var elementType ElementType
	if val <= math.MaxUint8 {
		elementType = UInt8
	} else if val <= math.MaxUint16 {
		elementType = UInt16
	} else if val <= math.MaxUint32 {
		elementType = UInt32
	} else {
		elementType = UInt64
	}
	return w.WriteElement(elementType, tag, val)
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
