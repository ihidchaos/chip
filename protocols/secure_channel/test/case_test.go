package test

import (
	"bytes"
	"fmt"
	"github.com/galenliu/chip/lib"
	"github.com/galenliu/chip/pkg/tlv"
	"github.com/galenliu/chip/protocols/secure_channel"
	"testing"
)

var kRootPubKeyFromSpec = []byte{
	0x04, 0x4a, 0x9f, 0x42, 0xb1, 0xca, 0x48, 0x40, 0xd3, 0x72, 0x92, 0xbb, 0xc7, 0xf6, 0xa7, 0xe1, 0x1e,
	0x22, 0x20, 0x0c, 0x97, 0x6f, 0xc9, 0x00, 0xdb, 0xc9, 0x8a, 0x7a, 0x38, 0x3a, 0x64, 0x1c, 0xb8, 0x25,
	0x4a, 0x2e, 0x56, 0xd4, 0xe2, 0x95, 0xa8, 0x47, 0x94, 0x3b, 0x4e, 0x38, 0x97, 0xc4, 0xa7, 0x73, 0xe9,
	0x30, 0x27, 0x7b, 0x4d, 0x9f, 0xbe, 0xde, 0x8a, 0x05, 0x26, 0x86, 0xbf, 0xac, 0xfa}

var kIpkOperationalGroupKeyFromSpec = []byte{
	0x9b, 0xc6, 0x1c, 0xd9, 0xc6, 0x2a, 0x2d, 0xf6, 0xd6, 0x4d, 0xfc, 0xaa, 0x9d, 0xc4, 0x72, 0xd4}

var kInitiatorRandomFromSpec = []byte{0x11, 0x11, 0x11, 0x31, 0x56, 0x8d, 0xfa, 0x17,
	0x20, 0x6b, 0x3a, 0xcc, 0xf8, 0xfa, 0xec, 0x2f,
	0x4d, 0x21, 0xb5, 0x80, 0x11, 0x31, 0x96, 0xf4,
	0x7c, 0x7c, 0x4d, 0xeb, 0x81, 0x01, 0x01, 0x01}

var kExpectedDestinationIdFromSpec = []byte{0xdc, 0x35, 0xdd, 0x5f, 0xc9, 0x13, 0x4c, 0xc5,
	0x54, 0x45, 0x38, 0xc9, 0xc3, 0xfc, 0x42, 0x97,
	0xc1, 0xec, 0x33, 0x70, 0xc8, 0x39, 0x13, 0x6a,
	0x80, 0xe1, 0x07, 0x96, 0x45, 0x1d, 0x4c, 0x53}

var kFabricIdFromSpec lib.FabricId = 0x2906C908D115D362
var kNodeIdFromSpec lib.NodeId = 0xCD5544AA7B13EF14

func EncodeSigma1() []byte {
	var kInitiatorRandomTag uint8 = 1
	var kInitiatorSessionIdTag uint8 = 2
	var kDestinationIdTag uint8 = 3
	var kInitiatorPubKeyTag uint8 = 4
	//var kInitiatorMRPParamsTag uint8 = 5
	tlvWriter := tlv.NewWriter()

	if _, err := tlvWriter.StartContainer(tlv.AnonymousTag(), tlv.TypeStructure); err != nil {
		fmt.Printf(err.Error())
	}
	if err := tlvWriter.PutBytes(tlv.ContextTag(kInitiatorRandomTag), kInitiatorRandomFromSpec); err != nil {
		fmt.Printf(err.Error())
	}
	if err := tlvWriter.PutU16(tlv.ContextTag(kInitiatorSessionIdTag), 0x1111); err != nil {
		fmt.Printf(err.Error())
	}

	if err := tlvWriter.PutBytes(tlv.ContextTag(kDestinationIdTag), kExpectedDestinationIdFromSpec); err != nil {
		fmt.Printf(err.Error())
	}

	if err := tlvWriter.PutBytes(tlv.ContextTag(kInitiatorPubKeyTag), kRootPubKeyFromSpec); err != nil {
		fmt.Printf(err.Error())
	}
	return tlvWriter.Bytes()
}

func TestParseSigma1(t *testing.T) {
	tlvReader := tlv.NewReader(bytes.NewBuffer(EncodeSigma1()))
	sigma1, err := secure_channel.ParseSigma1(tlvReader, false)
	if err != nil {
		t.Error(err)
	}
	t.Logf("ParseSigma1: %v", sigma1)

}
