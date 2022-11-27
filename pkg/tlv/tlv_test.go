package tlv

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func WriteIntMinMax() (data []byte, err error) {
	tlvWriter := NewWriter()
	//err = tlvWriter.StartContainer(AnonymousTag(), TypeStructure)
	//if err != nil {
	//	return nil, err
	//}
	_, err = tlvWriter.StartContainer(AnonymousTag(), TypeStructure)
	if err != nil {
		return
	}
	err = tlvWriter.PutUint(ContextTag(1), math.MaxUint8)
	if err != nil {
		return
	}
	err = tlvWriter.PutUint(ContextTag(2), math.MaxUint8)
	if err != nil {
		return
	}

	err = tlvWriter.EndContainer(TypeNotSpecified)
	if err != nil {
		return
	}

	err = tlvWriter.PutUint(AnonymousTag(), math.MaxUint16)
	if err != nil {
		return
	}
	err = tlvWriter.PutUint(AnonymousTag(), math.MaxUint16)
	if err != nil {
		return
	}

	err = tlvWriter.PutUint(AnonymousTag(), math.MaxUint32)
	if err != nil {
		return
	}
	err = tlvWriter.PutUint(AnonymousTag(), math.MaxUint32)
	if err != nil {
		return
	}

	err = tlvWriter.PutUint(AnonymousTag(), math.MaxUint64)
	if err != nil {
		return
	}
	err = tlvWriter.PutUint(AnonymousTag(), math.MaxUint64)
	if err != nil {
		return
	}

	err = tlvWriter.PutInt(AnonymousTag(), math.MaxInt8)
	if err != nil {
		return
	}
	err = tlvWriter.PutInt(AnonymousTag(), math.MaxInt8)
	if err != nil {
		return
	}

	err = tlvWriter.PutInt(AnonymousTag(), math.MaxInt16)
	if err != nil {
		return
	}
	err = tlvWriter.PutInt(AnonymousTag(), math.MaxInt16)
	if err != nil {
		return
	}

	err = tlvWriter.PutInt(AnonymousTag(), math.MaxInt32)
	if err != nil {
		return
	}
	err = tlvWriter.PutInt(AnonymousTag(), math.MaxInt32)
	if err != nil {
		return
	}

	err = tlvWriter.PutInt(AnonymousTag(), math.MaxInt64)
	if err != nil {
		return
	}
	err = tlvWriter.PutInt(AnonymousTag(), math.MinInt64)
	if err != nil {
		return
	}

	fmt.Printf("%0X", tlvWriter.Bytes())
	return tlvWriter.Bytes(), nil
}

func TestWriterMinMaxInt_Text(t *testing.T) {
	data, err := WriteIntMinMax()
	if err != nil {
		t.Error(err)
	}
	tlvReader := NewReader(bytes.NewReader(data))

	getUint8, err := tlvReader.GetU8()
	if err != nil {
		return
	}
	assert.Equal(t, math.MaxUint8, getUint8)
	assert.Equal(t, math.MaxUint8, getUint8)
}
