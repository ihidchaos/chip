package tlv

import (
	"github.com/galenliu/chip/lib"
	"github.com/galenliu/chip/platform/system/buffer"
	"io"
)

const (
	kControlByteNotSpecified = 0xFFFF
)

//var sTagSizes = []uint8{0, 1, 2, 4, 2, 4, 6, 8}

/**********************************************
* TAG
  | 111  xxxxx|  3bit  Tag Control
  | xxx  11111|  5bit  Tag TLVType
**********************************************/

type ReaderBase interface {
	Next()
	NextE(tag Tag, tlvTyp ...TLVType) error
	GetBytesView() ([]byte, error)
	EnterContainer() (TLVType, error)
	GetUint8() (uint8, error)
	GetUint16() (uint16, error)
	GetUint32() (uint32, error)
	GetUint64() (uint64, error)
	Tag() Tag
	ExitContainer(tlvType TLVType) error
}

type Reader struct {
	mBuffer           io.Reader
	mControlTag       TagControl
	mElementType      ElementType
	mContainerType    TLVType
	mElemTag          Tag
	mElemLenOrVal     uint64
	mControlByte      uint16
	mContainerOpen    bool
	ImplicitProfileId CommonProfiles
}

func NewReader(reader io.Reader) *Reader {
	r := &Reader{
		mBuffer:           reader,
		mControlTag:       0,
		mElementType:      NotSpecified,
		mContainerType:    TypeNotSpecified,
		mContainerOpen:    false,
		mControlByte:      0,
		mElemTag:          0,
		mElemLenOrVal:     0,
		ImplicitProfileId: kProfileNotSpecified,
	}
	r.ClearElementState()
	return r
}

// NextTT   读取一个指定tag的TLV
func (r *Reader) NextTT(tlvType TLVType, tag Tag) error {
	err := r.NextT(tag)
	if err != nil {
		return err
	}
	if r.Type() != tlvType {
		return lib.MATTER_ERROR_WRONG_TLV_TYPE
	}
	return nil
}

func (r *Reader) NextT(tag Tag) error {
	err := r.Next()
	if err != nil {
		return err
	}
	if r.mElemTag != tag {
		return lib.MATTER_ERROR_UNEXPECTED_TLV_ELEMENT
	}
	return nil
}

func (r *Reader) Next() error {

	err := r.Skip()
	if err != nil {
		return err
	}
	err = r.ReadElement()
	if err != nil {
		return err
	}
	elemType := r.ElementType()
	if elemType == EndOfContainer {
		return lib.MATTER_END_OF_TLV
	}
	return err
}

func (r *Reader) Skip() error {
	elemType := r.ElementType()
	if elemType == EndOfContainer {
		return lib.MATTER_END_OF_TLV
	}
	if elemType.IsContainer() {
		tlvType, err := r.EnterContainer()
		if err != nil {
			return err
		}
		err = r.ExitContainer(tlvType)
		if err != nil {
			return err
		}
	} else {
		err := r.SkipData()
		if err != nil {
			return err
		}
		r.ClearElementState()
	}
	return nil
}

func (r *Reader) Tag() Tag {
	return r.mElemTag
}

func (r *Reader) ReadElement() error {

	byt, err := buffer.Read8(r.mBuffer)
	if err != nil {
		return err
	}
	r.mControlTag = TagControl(byt & 0xE0)
	r.mElementType = ElementType(byt & 0x1F)

	r.mElemTag = r.ReadTag(r.mControlTag)

	//tagBytes := sTagSizes[r.mControlTag>>5]
	lenOrValFieldSize := r.mElementType.FieldSize()
	//valOrLenBytes := FieldSizeToBytes(lenOrValFieldSize)

	//elemHeadBytes := 1 + tagBytes + valOrLenBytes

	switch lenOrValFieldSize {
	case FieldSize0Byte:
		r.mElemLenOrVal = 0
	case FieldSize1Byte:
		val, err := buffer.Read8(r.mBuffer)
		if err != nil {
			return err
		}
		r.mElemLenOrVal = uint64(val)
	case FieldSize2Byte:
		val, err := buffer.LittleEndianRead16(r.mBuffer)
		if err != nil {
			return err
		}
		r.mElemLenOrVal = uint64(val)
	case FieldSize4Byte:
		val, err := buffer.LittleEndianRead32(r.mBuffer)
		if err != nil {
			return err
		}
		r.mElemLenOrVal = uint64(val)
	case FieldSize8Byte:
		val, err := buffer.LittleEndianRead64(r.mBuffer)
		if err != nil {
			return err
		}
		r.mElemLenOrVal = val
	}
	return nil
}

func (r *Reader) ElementType() ElementType {
	if r.mControlByte == kControlByteNotSpecified {
		return NotSpecified
	}
	return ElementType(r.mControlByte & fTLVTypeMask)
}

func (r *Reader) GetControlTag() TagControl {
	return r.mControlTag
}

func (r *Reader) GetBytes(io io.Reader) ([]byte, error) {

	var data = make([]byte, r.mElemLenOrVal)
	if r.TLVTypeIsContainer() {
		return nil, lib.WrongTlvType
	}
	return data, nil
}

func (r *Reader) GetBytesView() ([]byte, error) {
	if r.mElementType >= 0x0c && r.mElementType <= 0x13 {
		var data = make([]byte, r.mElemLenOrVal)
		_, err := r.mBuffer.Read(data)
		return data, err
	}
	return nil, lib.WrongTlvType
}

func (r *Reader) ReadTag(tagControl TagControl) Tag {
	switch tagControl {
	case ContextSpecific:
		val, _ := buffer.Read8(r.mBuffer)
		return ContextSpecificTag(val)
	case CommonProfile2Bytes:
		val, _ := buffer.LittleEndianRead16(r.mBuffer)
		return CommonTag2Byte(val)
	case CommonProfile4Bytes:
		val, _ := buffer.LittleEndianRead32(r.mBuffer)
		return CommonTag4Byte(val)
	//case ImplicitProfile2Bytes:
	//	if r.ImplicitProfileId == kProfileIdNotSpecified {
	//		return ContextSpecificTag(UnknownImplicitTag)
	//	}
	//	val, _ := buffer.LittleEndianRead16(r.mBuffer)
	//	return ProfileTag(r.ImplicitProfileId, uint32(val))
	//case ImplicitProfile4Bytes:
	//	if r.ImplicitProfileId == kProfileIdNotSpecified {
	//		return ContextSpecificTag(UnknownImplicitTag)
	//	}
	//	val, _ := buffer.LittleEndianRead32(r.mBuffer)
	//	return ProfileTag(r.ImplicitProfileId, val)
	case FullyQualified6Bytes:
		vendorId, _ := buffer.LittleEndianRead16(r.mBuffer)
		profileNum, _ := buffer.LittleEndianRead16(r.mBuffer)
		val, _ := buffer.LittleEndianRead16(r.mBuffer)
		return ProfileSpecificTag(vendorId, profileNum, uint32(val))
	case FullyQualified8Bytes:
		vendorId, _ := buffer.LittleEndianRead16(r.mBuffer)
		profileNum, _ := buffer.LittleEndianRead16(r.mBuffer)
		val, _ := buffer.LittleEndianRead32(r.mBuffer)
		return ProfileSpecificTag(vendorId, profileNum, val)
	default:
		return AnonymousTag()
	}
}

func (r *Reader) TLVTypeIsContainer() bool {
	elementType := r.ElementType()
	return elementType <= List && elementType <= Structure
}

func (r *Reader) VerifyElement() error {
	return nil
}

func (r *Reader) Type() TLVType {
	elemType := r.ElementType()
	if elemType == EndOfContainer {
		return TypeNotSpecified
	}
	if elemType == FloatingPointNumber32 || elemType == FloatingPointNumber64 {
		return TypeFloatingPointNumber
	}
	if elemType == NotSpecified || elemType == Null {
		return TLVType(elemType)
	}
	return TLVType(elemType)
}

func (r *Reader) GetUint8() (uint8, error) {
	v, err := r.GetUint64()
	return uint8(v), err
}

func (r *Reader) GetUint16() (uint16, error) {
	v, err := r.GetUint64()
	return uint16(v), err
}

func (r *Reader) GetUint32() (uint32, error) {
	v, err := r.GetUint64()
	return uint32(v), err
}

func (r *Reader) GetUint64() (uint64, error) {
	switch r.mElementType {
	case UInt8, UInt16, UInt32, UInt64:
		return r.mElemLenOrVal, nil
	default:
		return 0, lib.WrongTlvType
	}
}

func (r *Reader) reset() {
	r.mElemTag = 0
	r.mControlTag = 0
	r.mElementType = 0
	r.mElemLenOrVal = 0
	r.ImplicitProfileId = 0
}

func (r *Reader) EnterContainer() (TLVType, error) {
	t := r.Type()
	if t == TypeStructure || t == TypeList || t == TypeArray {
		return t, nil
	}
	return t, lib.WrongTlvType
}

func (r *Reader) ExitContainer(containerType TLVType) error {
	return nil
}

func (r *Reader) SkipData() error {
	elemType := r.ElementType()
	if elemType.HasLength() {
		err := r.ReadData(r.mElemLenOrVal)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Reader) ClearElementState() {
	r.mElemTag = AnonymousTag()
	r.mControlByte = kControlByteNotSpecified
	r.mElemLenOrVal = 0
}

func (r *Reader) ReadData(elemLenOrVal uint64) error {
	return nil
}

func FieldSizeToBytes(size FieldSize) uint8 {
	if size != FieldSize0Byte {
		return 1 << size
	}
	return 0
}
