package secure_channel

import (
	"github.com/galenliu/chip/crypto"
	"github.com/galenliu/chip/lib/tlv"
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

func ParseSigma1(tlvDecoder *tlv.Decoder, sessionResumptionRequested bool) (sigma1 Sigma1, err error) {
	sigma1 = Sigma1{
		initiatorResumeMICSize: crypto.AEADMicLengthBytes,
	}
	var kInitiatorRandomTag uint8 = 1
	var kInitiatorSessionIdTag uint8 = 2
	var kDestinationIdTag uint8 = 3
	var kInitiatorPubKeyTag uint8 = 4
	var kResumptionIDTag uint8 = 6
	//var kResume1MICTag uint8 = 7

	// Sigma1，这里应该读取到Structure 0x15

	containerType := tlv.TypeStructure
	err = tlvDecoder.NextType(containerType, tlv.AnonymousTag())
	if err != nil {
		return sigma1, err
	}
	containerType, err = tlvDecoder.EnterContainer()
	if err != nil {
		return sigma1, err
	}

	// Sigma1，NextTag = 1 initiatorRandom  32个字节的随机数
	err = tlvDecoder.NextTag(tlv.ContextTag(kInitiatorRandomTag))
	sigma1.initiatorRandom, err = tlvDecoder.GetBytes()
	if err != nil {
		return sigma1, err
	}

	//Sigma1， NextTag =2 Session id
	err = tlvDecoder.NextType(tlv.TypeUnsignedInteger, tlv.ContextTag(kInitiatorSessionIdTag))
	sigma1.initiatorSessionId, err = tlvDecoder.GetU16()
	if err != nil {
		return sigma1, err
	}

	//Sigma1，NextTag=3
	err = tlvDecoder.NextTag(tlv.ContextTag(kDestinationIdTag))
	sigma1.destinationId, err = tlvDecoder.GetBytes()
	if err != nil {
		return sigma1, err
	}

	//Sigma1，NextTag=4	 Initiator PubKey 65个字节的公钥
	err = tlvDecoder.NextTag(tlv.ContextTag(kInitiatorPubKeyTag))
	sigma1.initiatorEphPubKey, err = tlvDecoder.GetBytes()
	if err != nil {
		return sigma1, err
	}

	err = tlvDecoder.NextTag(tlv.ContextTag(kResumptionIDTag))
	sigma1.initiatorEphPubKey, err = tlvDecoder.GetBytes()
	if err != nil {
		return sigma1, err
	}

	//err = tlvDecoder.NextTag(tlv.ContextTag(kResumptionIDTag))
	//tlvDecoder.GetF32()

	return sigma1, tlvDecoder.ExitContainer(containerType)
}