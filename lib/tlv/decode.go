package tlv

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
)

var EndOfTLVError error = errors.New("end of tlv")
var ByteLenError error = errors.New("tlv data length err")

// A Decoder reads and decodes JSON values from an input stream.
type Decoder struct {
	r   io.Reader
	buf []byte

	p       int   // start of unread data in buf
	scanned int64 // amount of data already scanned
	err     error

	containerType  Type
	containerStack []Type

	containerOpen bool

	controlByte uint8
	elementType elementType
	tagControl  tagControl
	elemTag     Tag

	elemLenOrVal uint64

	implicitProfileId commonProfilesU32
}

// NewDecoder returns a new decoder that reads from r.
//
// The decoder introduces its own buffering and may
// readData data from r beyond the JSON values requested.
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{r: r, containerType: TypeUnknownContainer}
}

func (dec *Decoder) GetU8() (uint8, error) {
	v, err := dec.GetUint()
	if err != nil {
		return 0, err
	}
	if v > math.MaxUint8 {
		return 0, dec.valueError(v)
	}
	return uint8(v), err
}

func (dec *Decoder) GetU16() (uint16, error) {
	v, err := dec.GetUint()
	if err != nil {
		return 0, err
	}
	if v > math.MaxUint16 {
		return 0, dec.valueError(v)
	}
	return uint16(v), err
}

func (dec *Decoder) GetU32() (uint32, error) {
	v, err := dec.GetUint()
	if err != nil {
		return 0, err
	}
	if v > math.MaxUint32 {
		return 0, dec.valueError(v)
	}
	return uint32(v), err
}

func (dec *Decoder) GetUint() (uint64, error) {
	switch dec.elementType {
	case u8, u16, u32, u64:
		return dec.elemLenOrVal, nil
	default:
		return 0, dec.elemTypeError(dec.elementType)
	}
}

func (dec *Decoder) GetI8() (i8 int8, err error) {
	var v int64 = 0
	if v, err = dec.GetInt(); err != nil {
		return
	}
	if math.MinInt8 > v && v > math.MaxInt8 {
		return 0, dec.valueError(v)
	}
	i8 = int8(v)
	return
}

func (dec *Decoder) GetI16() (i16 int16, err error) {
	var v int64 = 0
	if v, err = dec.GetInt(); err != nil {
		return
	}
	if math.MinInt16 > v && v > math.MaxInt16 {
		return 0, dec.valueError(v)
	}
	i16 = int16(v)
	return
}

func (dec *Decoder) GetI32() (i32 int32, err error) {
	var v int64 = 0
	if v, err = dec.GetInt(); err != nil {
		return
	}
	if math.MinInt32 > v && v > math.MaxInt32 {
		return 0, dec.valueError(v)
	}
	i32 = int32(v)
	return
}

func (dec *Decoder) GetInt() (int64, error) {
	switch dec.elementType {
	case i8, i16, i32, i64:
		return int64(dec.elemLenOrVal), nil
	default:
		return 0, dec.elemTypeError(dec.elementType)
	}
}

func (dec *Decoder) GetF64() (f64 float64, err error) {
	switch dec.elementType {
	case floatingPointNumber64:
		f64 = math.Float64frombits(dec.elemLenOrVal)
		return
	default:
		return 0, dec.elemTypeError(dec.elementType)
	}
}

func (dec *Decoder) GetF32() (f32 float32, err error) {
	switch dec.elementType {
	case floatingPointNumber32:
		f32 = math.Float32frombits(uint32(dec.elemLenOrVal))
		return
	default:
		return 0, dec.elemTypeError(dec.elementType)
	}
}

func (dec *Decoder) GetBytes() (data []byte, err error) {
	if !dec.elementType.hasLength() || !dec.elementType.isValid() {
		return nil, dec.elemTypeError(dec.elementType)
	}
	data = make([]byte, dec.elemLenOrVal)
	_, err = dec.readData(data)
	return data, err
}

func (dec *Decoder) EnterContainer() (Type, error) {
	if !dec.elementType.isContainer() {
		return TypeUnknownContainer, dec.elemTypeError(dec.elementType)
	}
	outContainerType := dec.containerType
	dec.containerType = Type(dec.elementType)
	dec.clearElementState()
	return outContainerType, nil
}

func (dec *Decoder) clearElementState() {
	dec.elemTag = AnonymousTag()
	dec.controlByte = 0xFF
	dec.elemLenOrVal = 0
}

func (dec *Decoder) ExitContainer(outerContainerType Type) error {

	if err := dec.skipToEndContainer(); err != nil {
		return err
	}
	dec.containerType = outerContainerType
	dec.clearElementState()
	return nil
}

func (dec *Decoder) skipToEndContainer() error {
	return nil

	//outContainer := dec.containerType
	//dec.containerOpen = false
	//nestLevel := 0
	//for {
	//	if dec.elementType == endOfContainer {
	//		if nestLevel == 0 {
	//			return nil
	//		}
	//		nestLevel--
	//		if nestLevel == 0 {
	//			dec.containerType = outContainer
	//		} else {
	//			dec.containerType = TypeUnknownContainer
	//		}
	//
	//	} else if dec.elementType.isContainer() {
	//		nestLevel++
	//		dec.containerType = Type(dec.elementType)
	//	}
	//
	//	if err := dec.skipData(); err != nil {
	//		return err
	//	}
	//	if _, err := dec.readElement(dec.elementType); err != nil {
	//		return err
	//	}
	//}
}

func (dec *Decoder) Decode(v any) error {
	return nil
}

func (dec *Decoder) Tag(tag Tag) error {
	if err := dec.Next(); err != nil {
		return err
	}
	if dec.elemTag != tag {
		return dec.TagError(dec.elemTag)
	}
	return nil
}

func (dec *Decoder) Type(expectedType Type, expectedTag Tag) error {
	err := dec.Tag(expectedTag)
	if err != nil {
		return err
	}
	if dec.elementType != elementType(expectedType) {
		return dec.elemTypeError(dec.elementType)
	}
	return nil
}

func (dec *Decoder) Next() error {
	c, err := dec.byte()
	if err != nil {
		return err
	}
	tagControl := tagCtl(c)
	elementType := elemType(c)
	if !elementType.isValid() {
		return dec.elemTypeError(elementType)
	}
	if elementType == endOfContainer {
		return EndOfTLVError
	}
	if tag, err := dec.readTag(tagControl); err != nil {
		return err
	} else {
		dec.elemTag = tag
	}
	if elementType.hasValue() || elementType.hasLength() {
		data, err := dec.readElement(elementType)
		if err != nil {
			return err
		}
		dec.elemLenOrVal = data
	}
	dec.elementType = elementType
	dec.tagControl = tagControl
	return nil
}

func (dec *Decoder) readTag(tagControl tagControl) (tag Tag, err error) {
	var p []byte
	tag = AnonymousTag()
	err = nil
	switch tagControl {
	case ContextSpecific:
		p = make([]byte, 1)
		_, err = dec.readData(p)
		tag = ContextTag(p[0])
	case CommonProfile2Bytes:
		p = make([]byte, 2)
		_, err = dec.readData(p)
		val := binary.LittleEndian.Uint16(p)
		tag = commonTag(val)
	case CommonProfile4Bytes:
		p = make([]byte, 4)
		_, err = dec.readData(p)
		val := binary.LittleEndian.Uint32(p)
		tag = commonTag(val)
	case ImplicitProfile2Bytes:
		p = make([]byte, 2)
		_, err = dec.readData(p)
		val := binary.LittleEndian.Uint16(p)
		if dec.implicitProfileId == profileIdNotSpecified {
			tag = unknownImplicitTag
		} else {
			tag = profileTag(uint32(dec.implicitProfileId), val)
		}
	case ImplicitProfile4Bytes:
		p = make([]byte, 4)
		_, err = dec.readData(p)
		val := binary.LittleEndian.Uint32(p)
		if dec.implicitProfileId == profileIdNotSpecified {
			tag = unknownImplicitTag
		} else {
			tag = profileTag(uint32(dec.implicitProfileId), val)
		}
	case FullyQualified6Bytes:
		p = make([]byte, 2)
		_, err = dec.readData(p)
		vendorId := binary.LittleEndian.Uint16(p)
		_, err = dec.readData(p)
		profileNum := binary.LittleEndian.Uint16(p)
		_, err = dec.readData(p)
		val := binary.LittleEndian.Uint16(p)
		tag = profileSpecificTag(vendorId, profileNum, val)
	case FullyQualified8Bytes:
		p = make([]byte, 2)
		_, err = dec.readData(p)
		vendorId := binary.LittleEndian.Uint16(p)
		_, err = dec.readData(p)
		profileNum := binary.LittleEndian.Uint16(p)
		p = make([]byte, 4)
		_, err = dec.readData(p)
		val := binary.LittleEndian.Uint32(p)
		tag = profileSpecificTag(vendorId, profileNum, val)
	default:
		tag = AnonymousTag()
	}
	return
}

func (dec *Decoder) GetTag() Tag {
	return dec.elemTag
}

func (dec *Decoder) readData(p []byte) (int, error) {
	return dec.r.Read(p)
}

func (dec *Decoder) byte() (byte, error) {
	data := make([]byte, 1)
	_, err := dec.r.Read(data)
	return data[0], err
}

func (dec *Decoder) refill() error {
	// Make room to readData more into the buffer.
	// First slide down data already consumed.
	if dec.p > 0 {
		dec.scanned += int64(dec.p)
		n := copy(dec.buf, dec.buf[dec.p:])
		dec.buf = dec.buf[:n]
		dec.p = 0
	}

	// Grow buffer if not large enough.
	const minRead = 512
	if cap(dec.buf)-len(dec.buf) < minRead {
		newBuf := make([]byte, len(dec.buf), 2*cap(dec.buf)+minRead)
		copy(newBuf, dec.buf)
		dec.buf = newBuf
	}

	// readData. Delay error for next iteration (after scan).
	n, err := dec.r.Read(dec.buf[len(dec.buf):cap(dec.buf)])
	dec.buf = dec.buf[0 : len(dec.buf)+n]

	return err
}

func (dec *Decoder) readElement(elementType elementType) (val uint64, err error) {
	p := make([]byte, 1<<elementType.fieldSize())
	_, err = dec.readData(p)
	switch elementType.fieldSize() {
	case fieldSize0Byte:
		return 0, nil
	case fieldSize1Byte:
		return uint64(p[0]), nil
	case fieldSize2Byte:
		return uint64(binary.LittleEndian.Uint16(p)), nil
	case fieldSize4Byte:
		return uint64(binary.LittleEndian.Uint32(p)), nil
	case fieldSize8Byte:
		return binary.LittleEndian.Uint64(p), nil
	}
	return
}

func (dec *Decoder) skipData() error {
	if dec.elementType.hasLength() {
		_, _ = dec.readData(make([]byte, dec.elemLenOrVal))
	}
	return nil
}

func (dec *Decoder) elemTypeError(val any) error {
	return fmt.Errorf("element type err:%d", val)
}

func (dec *Decoder) valueError(val any) error {
	return fmt.Errorf("wrong value :%d", val)
}

func (dec *Decoder) TagError(val any) error {
	return fmt.Errorf("tag err :%d", val)
}

func (dec *Decoder) GetBoolean() (bool, error) {
	if dec.elementType == booleanFalse {
		return false, nil
	}
	if dec.elementType == booleanTrue {
		return true, nil
	}
	return false, dec.elemTypeError(dec.elementType)
}

func (dec *Decoder) GetData(data []byte) error {
	buf, err := dec.GetBytes()
	if err != nil {
		return err
	}
	if len(buf) != len(data) {
		return ByteLenError
	}
	copy(data, buf)
	return err
}
