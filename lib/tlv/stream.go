package tlv

import (
	"fmt"
	"io"
)

type Token any

// A Decoder reads and decodes JSON values from an input stream.
type Decoder struct {
	r       io.Reader
	buf     []byte
	d       decodeState
	scanp   int   // start of unread data in buf
	scanned int64 // amount of data already scanned
	scan    scanner
	err     error

	tokenState int
	tokenStack []int

	containerType  Type
	containerStack []Type

	controlByte uint8
	elementType ElementType
	tagControl  TagControl

	valOrLength []byte
}

// NewDecoder returns a new decoder that reads from r.
//
// The decoder introduces its own buffering and may
// read data from r beyond the JSON values requested.
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{r: r, containerType: typeUnknownContainer}
}

func (dec *Decoder) Token() (t Type, elementType ElementType, err error) {

	for {
		c, err := dec.peek()
		if err != nil {
			return -1, -1, err
		}

		tagControle := NewTagControl(c)
		elementType := NewElementType(c)

		if !elementType.IsValid() {
			return -1, -1, dec.elemTypeError(c)
		}

		if elementType.IsContainer() {
			if !dec.tokenValueAllowed() {
				return -1, -1, dec.tokenError(c)
			}
			switch elementType {
			case Structure:
				dec.containerStack = append(dec.containerStack, dec.containerType)
				dec.containerType = TypeStructure
			case Array:
				dec.containerStack = append(dec.containerStack, dec.containerType)
				dec.containerType = TypeArray
			case List:
				dec.containerStack = append(dec.containerStack, dec.containerType)
				dec.containerType = TypeList
			}
		}

		if elementType.HasValue() {
			switch elementType {
			case Int8:

			case Int16:
			case Int32:

			case Int64:

			case UInt8:

			case UInt16:

			case UInt32:

			case UInt64:

			case BooleanFalse:

			case BooleanTrue:

			case FloatingPointNumber32:

			case FloatingPointNumber64:

			}
		}

		switch c {
		case '[':
			if !dec.tokenValueAllowed() {
				return -1, -1, dec.tokenError(c)
			}
			dec.scanp++
			dec.tokenStack = append(dec.tokenStack, dec.tokenState)
			dec.tokenState = tokenArrayStart
			return Delim('['), nil

		case ']':
			if dec.tokenState != tokenArrayStart && dec.tokenState != tokenArrayComma {
				return dec.tokenError(c)
			}
			dec.scanp++
			dec.tokenState = dec.tokenStack[len(dec.tokenStack)-1]
			dec.tokenStack = dec.tokenStack[:len(dec.tokenStack)-1]
			dec.tokenValueEnd()
			return Delim(']'), nil

		case '{':
			if !dec.tokenValueAllowed() {
				return dec.tokenError(c)
			}
			dec.scanp++
			dec.tokenStack = append(dec.tokenStack, dec.tokenState)
			dec.tokenState = tokenObjectStart
			return Delim('{'), nil

		case '}':
			if dec.tokenState != tokenObjectStart && dec.tokenState != tokenObjectComma {
				return dec.tokenError(c)
			}
			dec.scanp++
			dec.tokenState = dec.tokenStack[len(dec.tokenStack)-1]
			dec.tokenStack = dec.tokenStack[:len(dec.tokenStack)-1]
			dec.tokenValueEnd()
			return Delim('}'), nil

		case ':':
			if dec.tokenState != tokenObjectColon {
				return dec.tokenError(c)
			}
			dec.scanp++
			dec.tokenState = tokenObjectValue
			continue

		case ',':
			if dec.tokenState == tokenArrayComma {
				dec.scanp++
				dec.tokenState = tokenArrayValue
				continue
			}
			if dec.tokenState == tokenObjectComma {
				dec.scanp++
				dec.tokenState = tokenObjectKey
				continue
			}
			return dec.tokenError(c)

		case '"':
			if dec.tokenState == tokenObjectStart || dec.tokenState == tokenObjectKey {
				var x string
				old := dec.tokenState
				dec.tokenState = tokenTopValue
				err := dec.Decode(&x)
				dec.tokenState = old
				if err != nil {
					return nil, err
				}
				dec.tokenState = tokenObjectColon
				return x, nil
			}
			fallthrough

		default:
			if !dec.tokenValueAllowed() {
				return dec.tokenError(c)
			}
			var x any
			if err := dec.Decode(&x); err != nil {
				return nil, err
			}
			return x, nil
		}
	}
}

func (dec *Decoder) peek() (byte, error) {
	var err error
	for {
		for i := dec.scanp; i < len(dec.buf); i++ {
			c := dec.buf[i]
			if isSpace(c) {
				continue
			}
			dec.scanp = i
			return c, nil
		}
		// buffer has been scanned, now report any error
		if err != nil {
			return 0, err
		}
		err = dec.refill()
	}
}

func (dec *Decoder) refill() error {
	// Make room to read more into the buffer.
	// First slide down data already consumed.
	if dec.scanp > 0 {
		dec.scanned += int64(dec.scanp)
		n := copy(dec.buf, dec.buf[dec.scanp:])
		dec.buf = dec.buf[:n]
		dec.scanp = 0
	}

	// Grow buffer if not large enough.
	const minRead = 512
	if cap(dec.buf)-len(dec.buf) < minRead {
		newBuf := make([]byte, len(dec.buf), 2*cap(dec.buf)+minRead)
		copy(newBuf, dec.buf)
		dec.buf = newBuf
	}

	// Read. Delay error for next iteration (after scan).
	n, err := dec.r.Read(dec.buf[len(dec.buf):cap(dec.buf)])
	dec.buf = dec.buf[0 : len(dec.buf)+n]

	return err
}

func (dec *Decoder) tokenValueAllowed() bool {
	switch dec.tokenState {
	case tokenTopValue, tokenArrayStart, tokenArrayValues, tokenObjectValue:
		return true
	}
	return false
}

func (dec *Decoder) tokenValueEnd() {
	switch dec.tokenState {
	case tokenArrayStart, tokenArrayValues:
		dec.tokenState = tokenArrayComma
	case tokenObjectValue:
		dec.tokenState = tokenObjectComma
	}
}

func (dec *Decoder) tokenError(c uint8) error {
	return fmt.Errorf("elemType error: %d", c)
}

func (dec *Decoder) elemTypeError(val uint8) error {
	return fmt.Errorf("element type err:%d", val)
}

const (
	tokenTopValue = iota
	tokenArrayStart
	tokenArrayValues
	tokenArrayComma
	tokenObjectStart
	tokenObjectKey
	tokenObjectColon
	tokenObjectValue
	tokenObjectComma
)
