package secure_channel

import "github.com/galenliu/chip/crypto"

type state uint8

const (
	initialized       state = 0
	sentSigma1        state = 1
	sentSigma2        state = 2
	sentSigma3        state = 3
	sentSigma1Resume  state = 4
	sentSigma2Resume  state = 5
	finished          state = 6
	finishedViaResume state = 7

	kSigmaParamRandomNumberSize = 32
	kIpkSize                    = crypto.SymmetricKeyLengthBytes
	TagTBEDataSenderNOC         = 1
	TagTBEDataSenderICAC        = 2
	TagTBEDataSignature         = 3
	TagTBEDataResumptionID      = 4

	kResumptionIdSize = 16

	ProtocolName        = "SecureChannel"
	ProtocolId   uint16 = 0x0000
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
	StatusReport              = 0x40
)

func (m MsgType) String() string {
	switch m {
	case MsgCounterSyncReq:
		return "MsgCounterSyncReq"
	case MsgCounterSyncRsp:
		return "MsgCounterSyncRsp"
	case StandaloneAck:
		return "StandaloneAck"
	case PBKDFParamRequest:
		return "PBKDFParamRequest"
	case PBKDFParamResponse:
		return "PBKDFParamResponse"
	case PASE_Pake1:
		return "PASE_Pake1"
	case PASE_Pake2:
		return "PASE_Pake2"
	case PASE_Pake3:
		return "PASE_Pake3"
	case CASE_Sigma1:
		return "CASE_Sigma1"
	case CASE_Sigma2:
		return "CASE_Sigma2"
	case CASE_Sigma3:
		return "CASE_Sigma3"
	case CASE_Sigma2Resume:
		return "CASE_Sigma2Resume"
	case StatusReport:
		return "StatusReport"
	default:
		return "----"
	}
}

var (
	KDFSR2Info     = []byte{0x53, 0x69, 0x67, 0x6d, 0x61, 0x32}
	KDFSR3Info     = []byte{0x53, 0x69, 0x67, 0x6d, 0x61, 0x33}
	kTBEData2Nonce = []byte{0x4e, 0x43, 0x41, 0x53, 0x45, 0x5f, 0x53, 0x69, 0x67, 0x6d, 0x61, 0x32, 0x4e}
)
