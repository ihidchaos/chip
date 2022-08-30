package secure

import "github.com/galenliu/chip/crypto"

const (
	StateInitialized            = 0
	StateSentSigma1             = 1
	StateSentSigma2             = 2
	StateSentSigma3             = 3
	StateSentSigma1Resume       = 4
	StateSentSigma2Resume       = 5
	StateFinished               = 6
	StateFinishedViaResume      = 7
	kSigmaParamRandomNumberSize = 32
	kIpkSize                    = crypto.SymmetricKeyLengthBytes

	TagTBEDataSenderNOC    = 1
	TagTBEDataSenderICAC   = 2
	TagTBEDataSignature    = 3
	TagTBEDataResumptionID = 4

	ResumptionIdSize = 16

	/* "NCASE_Sigma2N" */

)

var (
	KDFSR2Info     = []byte{0x53, 0x69, 0x67, 0x6d, 0x61, 0x32}
	KDFSR3Info     = []byte{0x53, 0x69, 0x67, 0x6d, 0x61, 0x33}
	kTBEData2Nonce = []byte{0x4e, 0x43, 0x41, 0x53, 0x45, 0x5f, 0x53, 0x69, 0x67, 0x6d, 0x61, 0x32, 0x4e}
)
