package tlv

import (
	"bytes"
	"testing"
)

func TestTLV(t *testing.T) {

	sigma1 := []byte{
		0x15,

		0x30, 0x01, 0x20,
		0x56, 0x8d, 0xfa, 0x17, 0x20, 0x56, 0x8d, 0xfa,
		0x6b, 0x3a, 0xcc, 0xf8, 0xfa, 0x80, 0x11, 0x31,
		0x56, 0x8d, 0xfa, 0x17, 0x20, 0x56, 0x8d, 0xfa,
		0x6b, 0x3a, 0xcc, 0xf8, 0xfa, 0x80, 0x11, 0x31,

		0x24, 0x02, 0xcf,

		0x30, 0x03, 0x20,
		0xa6, 0x8d, 0xfa, 0x17, 0x20, 0x6b, 0x3a, 0xcc,
		0x6b, 0x3a, 0xcc, 0xf8, 0xfa, 0xec, 0x2f, 0x4d,
		0xec, 0x2f, 0x4d, 0x21, 0xb5, 0x80, 0x11, 0xf4,
		0x80, 0x11, 0x31, 0x96, 0xf4, 0x31, 0x96, 0xf4,

		0x30, 0x04, 0x41,
		0xc6, 0x8d, 0xfa, 0x17, 0x20, 0x30, 0x20, 0x30,
		0x6b, 0x3a, 0xcc, 0xf8, 0xfa, 0x8d, 0x6b, 0x3a,
		0xec, 0x2f, 0x4d, 0x21, 0xb5, 0x80, 0x11, 0xb3,
		0x80, 0x11, 0x31, 0x96, 0xf4, 0x80, 0x11, 0x31,
		0x56, 0x8d, 0xfa, 0x17, 0x20, 0x8f, 0x11, 0x3f,
		0x6b, 0x3a, 0xcc, 0xf8, 0xfa, 0x3d, 0x11, 0xba,
		0xec, 0x2f, 0x4d, 0x21, 0xb5, 0xa0, 0x11, 0xa1,
		0x80, 0x11, 0x31, 0x96, 0xf4, 0xc2, 0x11, 0x3F,
		0x32,
		0x18,
	}

	var kInitiatorRandomTag uint8 = 1
	var kInitiatorSessionIdTag uint8 = 2
	var kDestinationIdTag uint8 = 3
	var kInitiatorPubKeyTag uint8 = 4
	//var kInitiatorMRPParamsTag uint8 = 5
	//var kResumptionIDTag uint8 = 6
	//var kResume1MICTag uint8 = 7

	// Sigma1，这里应该读取到Structure 0x15

	buf := bytes.NewBuffer(sigma1)
	var err error
	var sessionId uint16
	var initiatorRandom, destinationId, initiatorEphPubKey []byte

	tlvReader := NewReader(buf)
	err = tlvReader.NextE(AnonymousTag(), Type_Structure)
	if err != nil {
		t.Log(err.Error())
		return
	}
	// Sigma1，Tag = 1 initiatorRandom  20个字节的随机数
	err = tlvReader.NextE(ContextTag(kInitiatorRandomTag))
	if err != nil {
		t.Log(err.Error())
	}
	initiatorRandom, err = tlvReader.GetBytesView()
	if err != nil {
		t.Log(err.Error())
	}

	//Sigma1， Tag =2 Session id
	err = tlvReader.NextE(ContextTag(kInitiatorSessionIdTag), Type_UnsignedInteger)
	if err != nil {
		t.Log(err.Error())
	}
	sessionId, err = tlvReader.GetUint16()

	//Sigma1，Tag=3	destination id 20个字节的认证码
	err = tlvReader.NextE(ContextTag(kDestinationIdTag))
	destinationId, err = tlvReader.GetBytesView()
	if err != nil {
		t.Log(err.Error())
	}

	//Sigma1，Tag=4	 Initiator PubKey 1个字节的公钥
	err = tlvReader.NextE(ContextTag(kInitiatorPubKeyTag))
	initiatorEphPubKey, err = tlvReader.GetBytesView()
	if err != nil {
		t.Log(err.Error())
	}

	t.Logf("initiatorRandom: %X", initiatorRandom)
	t.Logf("sessionId: %X", sessionId)
	t.Logf("destinationId: %X", destinationId)
	t.Logf("initiatorEphPubKey: %X", initiatorEphPubKey)

	return
}
