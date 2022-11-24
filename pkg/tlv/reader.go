package tlv

import (
	"github.com/galenliu/chip/lib"
	"github.com/galenliu/chip/platform/system/buffer"
	"io"
	"math"
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
		mControlByte:      kControlByteNotSpecified,
		mElemTag:          0,
		mElemLenOrVal:     0,
		ImplicitProfileId: kProfileIdNotSpecified,
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
		outerContainerType, err := r.EnterContainer()
		if err != nil {
			return err
		}
		if err := r.ExitContainer(outerContainerType); err != nil {
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

func (r *Reader) ReadElement() (err error) {

	if err = r.EnsureData(); err != nil {
		return
	}

	if r.mControlByte, err = buffer.Read8(r.mBuffer); err != nil {
		return
	}

	byt, err := buffer.Read8(r.mBuffer)
	if err != nil {
		return err
	}
	r.mControlTag = TagControl(byt & 0xE0)
	r.mElementType = ElementType(byt & 0x1F)

	r.mElemTag = r.readTag(r.mControlTag)

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

func (r *Reader) readTag(tagControl TagControl) (tag Tag, err error) {

	switch tagControl {
	case ContextSpecific:
		val, err := buffer.Read8(r.mBuffer)
		return ContextSpecificTag(val), err
	case CommonProfile2Bytes:
		val, err := buffer.LittleEndianRead16(r.mBuffer)
		return CommonTag2Byte(val), err
	case CommonProfile4Bytes:
		val, err := buffer.LittleEndianRead32(r.mBuffer)
		return CommonTag4Byte(val), err
	case ImplicitProfile2Bytes:
		if r.ImplicitProfileId == kProfileIdNotSpecified {
			return ContextSpecificTag(UnknownImplicitTag), err
		}
		val, err := buffer.LittleEndianRead16(r.mBuffer)
		return ProfileTag(uint32(r.ImplicitProfileId), val), err
	case ImplicitProfile4Bytes:
		if r.ImplicitProfileId == kProfileIdNotSpecified {
			return ContextSpecificTag(UnknownImplicitTag), err
		}
		val, err := buffer.LittleEndianRead32(r.mBuffer)
		return ProfileTag(uint32(r.ImplicitProfileId), val), err
	case FullyQualified6Bytes:
		vendorId, _ := buffer.LittleEndianRead16(r.mBuffer)
		profileNum, _ := buffer.LittleEndianRead16(r.mBuffer)
		val, err := buffer.LittleEndianRead16(r.mBuffer)
		return ProfileSpecificTag(vendorId, profileNum, val), err
	case FullyQualified8Bytes:
		vendorId, err := buffer.LittleEndianRead16(r.mBuffer)
		profileNum, err := buffer.LittleEndianRead16(r.mBuffer)
		val, err := buffer.LittleEndianRead32(r.mBuffer)
		return ProfileSpecificTag(vendorId, profileNum, val), err
	default:
		return AnonymousTag(), err
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

func (r *Reader) GetU8() (uint8, error) {
	v, err := r.GetUint()
	if err != nil {
		return 0, err
	}
	if v > math.MaxUint8 {
		return 0, lib.MATTER_ERROR_INVALID_INTEGER_VALUE
	}
	return uint8(v), err
}

func (r *Reader) GetU6() (uint16, error) {
	v, err := r.GetUint()
	if err != nil {
		return 0, err
	}
	if v > math.MaxUint16 {
		return 0, lib.MATTER_ERROR_INVALID_INTEGER_VALUE
	}
	return uint16(v), err
}

func (r *Reader) GetU32() (uint32, error) {
	v, err := r.GetUint()
	if err != nil {
		return 0, err
	}
	if v > math.MaxUint32 {
		return 0, lib.MATTER_ERROR_INVALID_INTEGER_VALUE
	}
	return uint32(v), err
}

func (r *Reader) GetUint() (uint64, error) {
	switch r.mElementType {
	case UInt8, UInt16, UInt32, UInt64:
		return r.mElemLenOrVal, nil
	default:
		return 0, lib.MATTER_ERROR_WRONG_TLV_TYPE
	}
}

func (r *Reader) GetI8() (i8 int8, err error) {
	var v int64 = 0
	if v, err = r.GetInt(); err != nil {
		return
	}
	if math.MinInt8 > v && v > math.MaxInt8 {
		err = lib.MATTER_ERROR_INVALID_INTEGER_VALUE
		return
	}
	i8 = int8(v)
	return
}

func (r *Reader) GetI16() (i16 int16, err error) {
	var v int64 = 0
	if v, err = r.GetInt(); err != nil {
		return
	}
	if math.MinInt16 > v && v > math.MaxInt16 {
		err = lib.MATTER_ERROR_INVALID_INTEGER_VALUE
		return
	}
	i16 = int16(v)
	return
}

func (r *Reader) GetI32() (i32 int32, err error) {
	var v int64 = 0
	if v, err = r.GetInt(); err != nil {
		return
	}
	if math.MinInt32 > v && v > math.MaxInt32 {
		err = lib.MATTER_ERROR_INVALID_INTEGER_VALUE
		return
	}
	i32 = int32(v)
	return
}

func (r *Reader) GetInt() (int64, error) {
	switch r.mElementType {
	case Int8, Int16, Int32, Int64:
		return int64(r.mElemLenOrVal), nil
	default:
		return 0, lib.MATTER_ERROR_WRONG_TLV_TYPE
	}
}

func (r *Reader) EnterContainer() (TLVType, error) {
	elemType := r.ElementType()
	if !elemType.IsContainer() {
		return TypeUnknownContainer, lib.MATTER_ERROR_INCORRECT_STATE
	}
	r.mContainerType = TLVType(elemType)
	r.ClearElementState()
	return r.mContainerType, lib.WrongTlvType
}

func (r *Reader) ExitContainer(outerContainerType TLVType) error {

	if err := r.SkipToEndContainer(); err != nil {
		return err
	}
	r.mContainerType = outerContainerType
	r.ClearElementState()
	return nil
}

func (r *Reader) SkipData() error {
	elemType := r.ElementType()
	if elemType.HasLength() {
		if _, err := r.ReadData(r.mElemLenOrVal); err != nil {
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

func (r *Reader) ReadData(elemLenOrVal uint64) ([]byte, error) {
	if elemLenOrVal == 0 {
		return nil, lib.MATTER_ERROR_UNEXPECTED_TLV_ELEMENT
	}
	data := make([]byte, elemLenOrVal)
	if i, err := r.mBuffer.Read(data); err != nil || i != int(elemLenOrVal) {
		return nil, lib.MATTER_ERROR_TLV_UNDERRUN
	}
	return data, nil
}

func (r *Reader) SkipToEndContainer() error {
	outContainer := r.mContainerType
	r.SetContainerOpen(false)
	nestLevel := 0
	for {
		elemType := r.ElementType()
		if elemType == EndOfContainer {
			if nestLevel == 0 {
				return nil
			}
			nestLevel--
			if nestLevel == 0 {
				r.mContainerType = outContainer
			} else {
				r.mContainerType = TypeUnknownContainer
			}

		} else if elemType.IsContainer() {
			nestLevel++
			r.mContainerType = TLVType(elemType)
		}

		if err := r.SkipData(); err != nil {
			return err
		}
		if err := r.ReadElement(); err != nil {
			return err
		}

	}
}

func (r *Reader) SetContainerOpen(b bool) {
	r.mContainerOpen = b
}

func (r *Reader) EnsureData() error {
	return nil
}

func FieldSizeToBytes(size FieldSize) uint8 {
	if size != FieldSize0Byte {
		return 1 << size
	}
	return 0
}
