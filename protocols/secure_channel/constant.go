package secure_channel

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

type MsgType uint8

const (

	// Message Counter Synchronization Protocol Message Types
	MsgCounterSyncReq MsgType = 0x00
	MsgCounterSyncRsp MsgType = 0x01

	// Reliable Messaging Protocol Message Types
	StandaloneAck MsgType = 0x10

	// Password-based session establishment Message Types
	PBKDFParamRequest  MsgType = 0x20
	PBKDFParamResponse MsgType = 0x21
	PASE_Pake1         MsgType = 0x22
	PASE_Pake2         MsgType = 0x23
	PASE_Pake3         MsgType = 0x24

	// Certificate-based session establishment Message Types
	CASE_Sigma1       MsgType = 0x30
	CASE_Sigma2       MsgType = 0x31
	CASE_Sigma3       MsgType = 0x32
	CASE_Sigma2Resume MsgType = 0x33

	StatusReport = 0x40
)

var (
	KDFSR2Info     = []byte{0x53, 0x69, 0x67, 0x6d, 0x61, 0x32}
	KDFSR3Info     = []byte{0x53, 0x69, 0x67, 0x6d, 0x61, 0x33}
	kTBEData2Nonce = []byte{0x4e, 0x43, 0x41, 0x53, 0x45, 0x5f, 0x53, 0x69, 0x67, 0x6d, 0x61, 0x32, 0x4e}
)
