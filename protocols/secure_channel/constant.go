package secure_channel

import "github.com/galenliu/chip/crypto"

const (
	SInitialized                = 0
	SSentSigma1                 = 1
	SSentSigma2                 = 2
	SSentSigma3                 = 3
	SSentSigma1Resume           = 4
	kSentSigma2Resume           = 5
	kFinished                   = 6
	kFinishedViaResume          = 7
	kSigmaParamRandomNumberSize = 32
	kIpkSize                    = crypto.SymmetricKeyLengthBytes

	kTag_TBEData_SenderNOC    = 1
	kTag_TBEData_SenderICAC   = 2
	kTag_TBEData_Signature    = 3
	kTag_TBEData_ResumptionID = 4

	kResumptionIdSize = 16
)

var (
	KDFSR2Info = []byte{0x53, 0x69, 0x67, 0x6d, 0x61, 0x32}
	KDFSR3Info = []byte{0x53, 0x69, 0x67, 0x6d, 0x61, 0x33}
)
