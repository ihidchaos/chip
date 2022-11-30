package tlv

import (
	"github.com/galenliu/chip/lib"
	"github.com/galenliu/chip/platform/system/buffer"
	"io"
	"math"
)

type TLVReader interface {
	io.Reader
	Len() int
}

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
	mBuffer           TLVReader
	mControlTag       TagControl
	mElemTag          Tag
	mContainerType    TLVType
	mElemLenOrVal     uint64
	mControlByte      uint16
	mContainerOpen    bool
	ImplicitProfileId CommonProfiles
}

func NewReader(buf TLVReader) *Reader {
	r := &Reader{
		mBuffer:           buf,
		mControlTag:       0,
		mElemTag:          AnonymousTag(),
		mContainerType:    TypeNotSpecified,
		mContainerOpen:    false,
		mControlByte:      kControlByteNotSpecified,
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
	elemType := r.ElementType()
	err := r.Skip()
	if err != nil {
		return err
	}
	err = r.readElement()
	if err != nil {
		return err
	}
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

func (r *Reader) ElementType() ElementType {
	if r.mControlByte == kControlByteNotSpecified {
		return NotSpecified
	}
	return ElementType(r.mControlByte & fTLVTypeMask)
}

func (r *Reader) GetControlTag() TagControl {
	return r.mControlTag
}

func (r *Reader) GetBytes() ([]byte, error) {

	if r.ElementType().IsContainer() {
		return nil, lib.MATTER_ERROR_WRONG_TLV_TYPE
	}
	return r.readData(r.mBuffer, r.mElemLenOrVal)
}

func (r *Reader) VerifyElement() error {
	if r.ElementType() == EndOfContainer {
		if r.mContainerType == TypeNotSpecified {
			return lib.MATTER_ERROR_INVALID_TLV_ELEMENT
		}
		if r.mElemTag != AnonymousTag() {
			return lib.MATTER_ERROR_INVALID_TLV_TAG
		}
	} else {
		if r.mElemTag == UnknownImplicitTag {
			return lib.MATTER_ERROR_UNKNOWN_IMPLICIT_TLV_TAG
		}
		switch r.mContainerType {
		case TypeNotSpecified:
			if r.mElemTag.IsContext() {
				return lib.MATTER_ERROR_UNKNOWN_IMPLICIT_TLV_TAG
			}
		case TypeStructure:
			if r.mElemTag == AnonymousTag() {
				return lib.MATTER_ERROR_UNKNOWN_IMPLICIT_TLV_TAG
			}
		case TypeArray:
			if r.mElemTag != AnonymousTag() {
				return lib.MATTER_ERROR_UNKNOWN_IMPLICIT_TLV_TAG
			}
		case TypeUnknownContainer:
			break
		case TypeList:
			break
		default:
			return lib.MATTER_ERROR_INCORRECT_STATE
		}
	}
	if r.ElementType().HasLength() {
		if r.mBuffer.Len() < int(r.ElementType().FieldSize().ByteSize()) {
			return lib.MATTER_ERROR_TLV_UNDERRUN
		}
	}
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

func (r *Reader) GetU16() (uint16, error) {
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
	t := r.ElementType()
	switch t {
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
	switch r.ElementType() {
	case Int8, Int16, Int32, Int64:
		return int64(r.mElemLenOrVal), nil
	default:
		return 0, lib.MATTER_ERROR_WRONG_TLV_TYPE
	}
}

func (r *Reader) GetF64() (f64 float64, err error) {
	switch r.ElementType() {
	case FloatingPointNumber32:
		f64 = math.Float64frombits(r.mElemLenOrVal)
		return
	default:
		return 0, lib.MATTER_ERROR_WRONG_TLV_TYPE
	}
}

func (r *Reader) GetF32() (f32 float32, err error) {
	switch r.ElementType() {
	case FloatingPointNumber32:
		f32 = math.Float32frombits(uint32(r.mElemLenOrVal))
		return
	default:
		return 0, lib.MATTER_ERROR_WRONG_TLV_TYPE
	}
}

func (r *Reader) EnterContainer() (TLVType, error) {
	elemType := r.ElementType()
	if !elemType.IsContainer() {
		return TypeUnknownContainer, lib.MATTER_ERROR_INCORRECT_STATE
	}
	outContainerType := r.mContainerType
	r.mContainerType = TLVType(elemType)
	r.ClearElementState()
	return outContainerType, nil
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
	//elemType := r.ElementType()
	//if elemType.HasLength() {
	//	if _, err := r.readData(r.mBuffer, r.mElemLenOrVal); err != nil {
	//		return err
	//	}
	//}
	return nil
}

func (r *Reader) ClearElementState() {
	r.mElemTag = AnonymousTag()
	r.mControlByte = kControlByteNotSpecified
	r.mElemLenOrVal = 0
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
		if err := r.readElement(); err != nil {
			return err
		}

	}
}

func (r *Reader) SetContainerOpen(b bool) {
	r.mContainerOpen = b
}

func (r *Reader) ensureData(buf TLVReader) error {
	if buf.Len() == 0 {
		return lib.MATTER_ERROR_TLV_UNDERRUN
	}
	return nil
}

func (r *Reader) readData(buf TLVReader, elemLenOrVal uint64) (data []byte, err error) {
	if err = r.ensureData(buf); err != nil {
		return nil, err
	}
	if elemLenOrVal == 0 {
		return nil, lib.MATTER_ERROR_UNEXPECTED_TLV_ELEMENT
	}
	data = make([]byte, elemLenOrVal)
	if i, err := buf.Read(data); err != nil || i != int(elemLenOrVal) {
		return nil, lib.MATTER_ERROR_TLV_UNDERRUN
	}
	return data, nil
}

func (r *Reader) readElement() (err error) {

	if err = r.ensureData(r.mBuffer); err != nil {
		return
	}
	if u8, err := buffer.Read8(r.mBuffer); err != nil {
		return err
	} else {
		r.mControlByte = uint16(u8)
	}

	elemType := r.ElementType()
	if !elemType.IsValid() {
		return lib.MATTER_ERROR_INVALID_TLV_ELEMENT
	}

	tagControl := ParseTagControl(r.mControlByte)

	if r.mElemTag, err = r.readTag(tagControl); err != nil {
		return err
	}
	switch elemType.FieldSize() {
	case FieldSize0Byte:
		r.mElemLenOrVal = 0
	case FieldSize1Byte:
		if val, err := buffer.Read8(r.mBuffer); err != nil {
			return err
		} else {
			r.mElemLenOrVal = uint64(val)
		}
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
	return r.VerifyElement()
}

func (r *Reader) readTag(tagControl TagControl) (tag Tag, err error) {

	switch tagControl {
	case ContextSpecific:
		val, err := buffer.Read8(r.mBuffer)
		return ContextTag(val), err
	case CommonProfile2Bytes:
		val, err := buffer.LittleEndianRead16(r.mBuffer)
		return CommonTag(val), err
	case CommonProfile4Bytes:
		val, err := buffer.LittleEndianRead32(r.mBuffer)
		return CommonTag(val), err
	case ImplicitProfile2Bytes:
		if r.ImplicitProfileId == kProfileIdNotSpecified {
			return UnknownImplicitTag, err
		}
		val, err := buffer.LittleEndianRead16(r.mBuffer)
		return ProfileTag(uint32(r.ImplicitProfileId), val), err
	case ImplicitProfile4Bytes:
		if r.ImplicitProfileId == kProfileIdNotSpecified {
			return UnknownImplicitTag, err
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
