package secure_channel

import (
	"github.com/galenliu/chip/crypto"
	"github.com/galenliu/chip/lib"
	"github.com/galenliu/chip/lib/tlv"
)

type Sigma1 struct {
	initiatorRandom        []byte
	initiatorSessionId     uint16
	destinationId          []byte
	destinationIdSize      int
	initiatorEphPubKey     []byte
	initiatorSEDParams     any      //optional
	resumptionID           [16]byte //optional
	initiatorResumeMIC     []byte   //optional
	initiatorResumeMICSize int
}

func ParseSigma1(tlvReader tlv.Reader) (sigma1 Sigma1, err error) {
	sigma1 = Sigma1{
		initiatorResumeMICSize: crypto.AeadMicLengthBytes,
		destinationIdSize:      crypto.HashLenBytes,
	}
	var kInitiatorRandomTag uint8 = 1
	var kInitiatorSessionIdTag uint8 = 2
	var kDestinationIdTag uint8 = 3
	var kInitiatorPubKeyTag uint8 = 4
	var kInitiatorMRPParamsTag uint8 = 5
	//var kResumptionIDTag uint8 = 6
	//var kResume1MICTag uint8 = 7

	// Sigma1，这里应该读取到Structure 0x15

	err = tlvReader.NextE(tlv.AnonymousTag(), tlv.TypeStructure)
	containerType, err := tlvReader.EnterContainer()
	if err != nil {
		return
	}
	// Sigma1，Tag = 1 initiatorRandom  20个字节的随机数
	err = tlvReader.NextE(tlv.ContextTag(kInitiatorRandomTag))
	sigma1.initiatorRandom, err = tlvReader.GetBytesView()
	if err != nil && len(sigma1.initiatorRandom) != 32 {
		err = lib.ChipErrorInvalidCaseParameter
		return
	}

	//Sigma1， Tag =2 Session id
	err = tlvReader.NextE(tlv.ContextTag(kInitiatorSessionIdTag), tlv.TypeUnsignedInteger)
	sigma1.initiatorSessionId, err = tlvReader.GetUint16()
	if err != nil {
		return
	}

	//Sigma1，Tag=3	destination id 20个字节的认证码
	err = tlvReader.NextE(tlv.ContextTag(kDestinationIdTag))
	sigma1.destinationId, err = tlvReader.GetBytesView()
	if err != nil && len(sigma1.destinationId) != crypto.KSha256HashLength {
		err = lib.ChipErrorInvalidCaseParameter
		return
	}

	//Sigma1，Tag=4	 Initiator PubKey 1个字节的公钥
	err = tlvReader.NextE(tlv.ContextTag(kInitiatorPubKeyTag))
	sigma1.initiatorEphPubKey, err = tlvReader.GetBytesView()
	if err != nil && len(sigma1.initiatorEphPubKey) != crypto.KP256PublicKeyLength {
		err = lib.ChipErrorInvalidCaseParameter
		return
	}

	tlvReader.Next()
	if tlvReader.GetTag() == tlv.ContextTag(kInitiatorMRPParamsTag) {
		//s.DecodeMRPParametersIfPresent(tlv.ContextTag(kInitiatorMRPParamsTag), tlvReader)
	}
	err = tlvReader.ExitContainer(containerType)
	return
}
