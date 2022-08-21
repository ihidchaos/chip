package tlv

import (
	"github.com/galenliu/chip/lib"
	"github.com/galenliu/chip/lib/buffer"
	"github.com/galenliu/gateway/pkg/log"
	"io"
)

//var sTagSizes = []uint8{0, 1, 2, 4, 2, 4, 6, 8}

/**********************************************
* TAG
  | 111  xxxxx|  3bit  ElementTag Control
  | xxx  11111|  5bit  ElementTag TLVType
**********************************************/

type Reader interface {
	Next()
	NextE(tag ElementTag, tlvTyp ...TLVType) error
	GetBytesView() ([]byte, error)
	EnterContainer() (TLVType, error)
	GetUint8() (uint8, error)
	GetUint16() (uint16, error)
	GetUint32() (uint32, error)
	GetUint64() (uint64, error)
	GetTag() ElementTag
	ExitContainer(tlvType TLVType) error
}

type ReaderImpl struct {
	mBuffer      io.Reader
	mControlTag  TagControl
	mElementType ElementType

	mElemTag ElementTag

	mElemLenOrVal uint64

	ImplicitProfileId uint32
}

func NewReader(reader io.Reader) *ReaderImpl {
	return &ReaderImpl{
		mBuffer:           reader,
		mControlTag:       0,
		mElementType:      0,
		mElemTag:          0,
		mElemLenOrVal:     0,
		ImplicitProfileId: 0,
	}
}

// NextE 读取一个指定tag的TLV
func (r *ReaderImpl) NextE(tag ElementTag, tlvTyp ...TLVType) error {
	r.Next()
	if r.mElemTag != tag {
		return lib.ChipErrorUnexpectedTlvElement
	}
	if len(tlvTyp) > 0 {
		if r.GetType() != tlvTyp[0] {
			return lib.ChipErrorWrongTlvType
		}
	}
	log.Infof("Elem tag:%X", r.mElemTag)
	log.Infof("ControlTag:%X", r.mControlTag)
	log.Infof("Element Type:%X", r.mElementType)
	log.Infof("Elem Len Or Val:%X", r.mElemLenOrVal)
	log.Infof("--------------------")
	return r.VerifyElement()
}

func (r *ReaderImpl) Next() {
	r.reset()
	err := r.ReadElement()
	if err != nil {
		log.Errorf(err.Error())
	}
}

func (r *ReaderImpl) GetTag() ElementTag {
	return r.mElemTag
}

func (r *ReaderImpl) ReadElement() error {
	//var err error
	//stagingBuf := make([]byte, 17)
	//var elemType ElementType
	byt, err := buffer.Read8(r.mBuffer)
	if err != nil {
		return err
	}
	r.mControlTag = TagControl(byt & 0xE0)
	r.mElementType = ElementType(byt & 0x1F)

	r.mElemTag = r.ReadTag(r.mControlTag)

	//tagBytes := sTagSizes[r.mControlTag>>5]
	lenOrValFieldSize := GetTLVFieldSize(r.mElementType)
	//valOrLenBytes := FieldSizeToBytes(lenOrValFieldSize)

	//elemHeadBytes := 1 + tagBytes + valOrLenBytes

	switch lenOrValFieldSize {
	case kTLVFieldSize0Byte:
		r.mElemLenOrVal = 0
	case kTLVFieldSize1Byte:
		val, err := buffer.Read8(r.mBuffer)
		if err != nil {
			return err
		}
		r.mElemLenOrVal = uint64(val)
	case kTLVFieldSize2Byte:
		val, err := buffer.Read16(r.mBuffer)
		if err != nil {
			return err
		}
		r.mElemLenOrVal = uint64(val)
	case kTLVFieldSize4Byte:
		val, err := buffer.Read32(r.mBuffer)
		if err != nil {
			return err
		}
		r.mElemLenOrVal = uint64(val)
	case kTLVFieldSize8Byte:
		val, err := buffer.Read64(r.mBuffer)
		if err != nil {
			return err
		}
		r.mElemLenOrVal = val
	}
	return nil
}

func (r *ReaderImpl) ElementType() ElementType {
	return r.mElementType
}

func (r *ReaderImpl) GetControlTag() TagControl {
	return r.mControlTag
}

func (r *ReaderImpl) GetBytes(io io.Reader) ([]byte, error) {

	var data = make([]byte, r.mElemLenOrVal)
	if r.TLVTypeIsContainer() {
		return nil, lib.ChipErrorWrongTlvType
	}
	return data, nil
}

func (r *ReaderImpl) GetBytesView() ([]byte, error) {
	if r.mElementType >= 0x0c && r.mElementType <= 0x13 {
		var data = make([]byte, r.mElemLenOrVal)
		_, err := r.mBuffer.Read(data)
		return data, err
	}
	return nil, lib.ChipErrorWrongTlvType
}

func (r *ReaderImpl) ReadTag(tagControl TagControl) ElementTag {
	switch tagControl {
	case ContextSpecific:
		val, _ := buffer.Read8(r.mBuffer)
		return ContextTag(val)
	case CommonProfile2Bytes:
		val, _ := buffer.Read16(r.mBuffer)
		return CommonTag2Byte(val)
	case CommonProfile4Bytes:
		val, _ := buffer.Read32(r.mBuffer)
		return CommonTag4Byte(val)
	//case ImplicitProfile2Bytes:
	//	if r.ImplicitProfileId == kProfileIdNotSpecified {
	//		return ContextTag(UnknownImplicitTag)
	//	}
	//	val, _ := buffer.Read16(r.mBuffer)
	//	return ProfileTag(r.ImplicitProfileId, uint32(val))
	//case ImplicitProfile4Bytes:
	//	if r.ImplicitProfileId == kProfileIdNotSpecified {
	//		return ContextTag(UnknownImplicitTag)
	//	}
	//	val, _ := buffer.Read32(r.mBuffer)
	//	return ProfileTag(r.ImplicitProfileId, val)
	case FullQualified6Bytes:
		vendorId, _ := buffer.Read16(r.mBuffer)
		profileNum, _ := buffer.Read16(r.mBuffer)
		val, _ := buffer.Read16(r.mBuffer)
		return ProfileTag4Byte(vendorId, profileNum, uint32(val))
	case FullyQualified8Bytes:
		vendorId, _ := buffer.Read16(r.mBuffer)
		profileNum, _ := buffer.Read16(r.mBuffer)
		val, _ := buffer.Read32(r.mBuffer)
		return ProfileTag4Byte(vendorId, profileNum, val)
	default:
		return AnonymousTag()
	}
}

func (r *ReaderImpl) TLVTypeIsContainer() bool {
	elementType := r.ElementType()
	return elementType <= List && elementType <= Structure
}

func (r *ReaderImpl) VerifyElement() error {
	return nil
}

func (r *ReaderImpl) GetType() TLVType {
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

func (r *ReaderImpl) GetUint8() (uint8, error) {
	v, err := r.GetUint64()
	return uint8(v), err
}

func (r *ReaderImpl) GetUint16() (uint16, error) {
	v, err := r.GetUint64()
	return uint16(v), err
}

func (r *ReaderImpl) GetUint32() (uint32, error) {
	v, err := r.GetUint64()
	return uint32(v), err
}

func (r *ReaderImpl) GetUint64() (uint64, error) {
	switch r.mElementType {
	case UInt8, UInt16, UInt32, UInt64:
		return r.mElemLenOrVal, nil
	default:
		return 0, lib.ChipErrorWrongTlvType
	}
}

func (r *ReaderImpl) reset() {
	r.mElemTag = 0
	r.mControlTag = 0
	r.mElementType = 0
	r.mElemLenOrVal = 0
	r.ImplicitProfileId = 0
}

func (r *ReaderImpl) EnterContainer() (TLVType, error) {
	t := r.GetType()
	if t == TypeStructure || t == TypeList || t == TypeArray {
		return t, nil
	}
	return t, lib.ChipErrorWrongTlvType
}

func (r *ReaderImpl) ExitContainer(containerType TLVType) error {
	return r.NextE(AnonymousTag(), containerType)
}

func FieldSizeToBytes(size FieldSize) uint8 {
	if size != kTLVFieldSize0Byte {
		return 1 << size
	}
	return 0
}

func GetTLVFieldSize(elementType ElementType) FieldSize {
	if elementType.HasValue() {
		return FieldSize(uint8(elementType) & 0x03)
	}
	return kTLVFieldSize0Byte
}
