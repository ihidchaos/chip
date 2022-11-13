package secure_channel

import (
	"github.com/galenliu/chip/crypto"
	"github.com/galenliu/chip/lib"
	tlv2 "github.com/galenliu/chip/pkg/tlv"
)

type Sigma1 struct {
	initiatorRandom    []byte
	initiatorSessionId uint16
	destinationId      []byte

	initiatorEphPubKey     []byte //0x41 bytes
	initiatorSEDParams     any    //optional
	resumptionId           []byte //optional
	initiatorResumeMIC     []byte //optional
	initiatorResumeMICSize int

	sessionResumptionRequested bool
}

func ParseSigma1(tlvReader tlv2.Reader) (sigma1 Sigma1, err error) {
	sigma1 = Sigma1{
		initiatorResumeMICSize: crypto.AEADMicLengthBytes,
	}
	var kInitiatorRandomTag uint8 = 1
	var kInitiatorSessionIdTag uint8 = 2
	var kDestinationIdTag uint8 = 3
	var kInitiatorPubKeyTag uint8 = 4
	var kInitiatorMRPParamsTag uint8 = 5
	//var kResumptionIDTag uint8 = 6
	//var kResume1MICTag uint8 = 7

	// Sigma1，这里应该读取到Structure 0x15

	err = tlvReader.NextE(tlv2.AnonymousTag(), tlv2.Type_Structure)
	containerType, err := tlvReader.EnterContainer()
	if err != nil {
		return
	}
	// Sigma1，Tag = 1 initiatorRandom  20个字节的随机数
	err = tlvReader.NextE(tlv2.ContextTag(kInitiatorRandomTag))
	sigma1.initiatorRandom, err = tlvReader.GetBytesView()
	if err != nil && len(sigma1.initiatorRandom) != 32 {
		err = lib.InvalidCaseParameter
		return
	}

	//Sigma1， Tag =2 Session id
	err = tlvReader.NextE(tlv2.ContextTag(kInitiatorSessionIdTag), tlv2.Type_UnsignedInteger)
	sigma1.initiatorSessionId, err = tlvReader.GetUint16()
	if err != nil {
		return
	}

	//Sigma1，Tag=3	destination id 20个字节的认证码
	err = tlvReader.NextE(tlv2.ContextTag(kDestinationIdTag))
	sigma1.destinationId, err = tlvReader.GetBytesView()
	if err != nil && len(sigma1.destinationId) != crypto.KSha256HashLength {
		err = lib.InvalidCaseParameter
		return
	}

	//Sigma1，Tag=4	 Initiator PubKey 1个字节的公钥
	err = tlvReader.NextE(tlv2.ContextTag(kInitiatorPubKeyTag))
	sigma1.initiatorEphPubKey, err = tlvReader.GetBytesView()
	if err != nil && len(sigma1.initiatorEphPubKey) != crypto.KP256PublicKeyLength {
		err = lib.InvalidCaseParameter
		return
	}

	tlvReader.Next()
	if tlvReader.GetTag() == tlv2.ContextTag(kInitiatorMRPParamsTag) {
		//s.DecodeMRPParametersIfPresent(tlv.ContextTag(kInitiatorMRPParamsTag), tlvReader)
	}
	err = tlvReader.ExitContainer(containerType)
	return
}
