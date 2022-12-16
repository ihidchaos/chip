package secure_channel

import (
	"github.com/galenliu/chip/crypto"
	"github.com/galenliu/chip/protocols"
	"time"
)

type state uint8

var mod = "SecureChanel"

var Id = protocols.New(protocolId, nil)

// Message Counter Synchronization Protocol Message Types
const (
	MsgCounterSyncReq MsgType = 0x00
	MsgCounterSyncRsp MsgType = 0x01
)

// StandaloneAck Reliable Messaging Protocol Message Types
const (
	StandaloneAck MsgType = 0x10
)

// Password-based session establishment Message Types
const (
	PBKDFParamRequest  MsgType = 0x20
	PBKDFParamResponse MsgType = 0x21
	PASEPake1          MsgType = 0x22
	PASEPake2          MsgType = 0x23
	PASEPake3          MsgType = 0x24
)

// Certificate-based session establishment Message Types
const (
	CASESigma1       MsgType = 0x30
	CASESigma2       MsgType = 0x31
	CASESigma3       MsgType = 0x32
	CASESigma2Resume MsgType = 0x33
	CASEStatusReport MsgType = 0x40
)

const (
	initialized       state = 0
	sentSigma1        state = 1
	sentSigma2        state = 2
	sentSigma3        state = 3
	sentSigma1Resume  state = 4
	sentSigma2Resume  state = 5
	finished          state = 6
	finishedViaResume state = 7

	sigmaParamRandomNumberSize = 32
	ipkSize                    = crypto.SymmetricKeyLengthBytes

	TagTBEDataSenderNOC    = 1
	TagTBEDataSenderICAC   = 2
	TagTBEDataSignature    = 3
	TagTBEDataResumptionID = 4

	kExpectedHighProcessingTime = time.Duration(30 * time.Second)

	resumptionIdSize        = 16
	ProtocolName            = "SecureChannel"
	protocolId       uint16 = 0x0000
	/* "NCASE_Sigma2N" */
)

type MsgType uint8

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
	case PASEPake1:
		return "PASEPake1"
	case PASEPake2:
		return "PASEPake2"
	case PASEPake3:
		return "PASEPake3"
	case CASESigma1:
		return "CASESigma1"
	case CASESigma2:
		return "CASESigma2"
	case CASESigma3:
		return "CASESigma3"
	case CASESigma2Resume:
		return "CASESigma2Resume"
	case CASEStatusReport:
		return "CASEStatusReport"
	default:
		return "----"
	}
}

func (m MsgType) MessageType() uint8 {
	return uint8(m)
}

func (m MsgType) ProtocolId() protocols.Id {
	return protocols.New(protocolId, nil)
}

var (
	KDFSR2Info     = []byte{0x53, 0x69, 0x67, 0x6d, 0x61, 0x32}
	KDFSR3Info     = []byte{0x53, 0x69, 0x67, 0x6d, 0x61, 0x33}
	kTBEData2Nonce = []byte{0x4e, 0x43, 0x41, 0x53, 0x45, 0x5f, 0x53, 0x69, 0x67, 0x6d, 0x61, 0x32, 0x4e}
)

const (
	protocolCodeSuccess         uint16 = 0x0000
	protocolCodeNoSharedRoot    uint16 = 0x0001
	protocolCodeInvalidParam    uint16 = 0x0002
	protocolCodeCloseSession    uint16 = 0x0003
	protocolCodeBusy            uint16 = 0x0004
	protocolCodeSessionNotFound uint16 = 0x0005
)

type generalStatusCode uint16

const (
	kSuccess           generalStatusCode = 0 /**< Operation completed successfully. */
	kFailure           generalStatusCode = 1 /**< Generic failure, additional details may be included in the protocol specific status. */
	kBadPrecondition   generalStatusCode = 2 /**< Operation was rejected by the system because the system is in an invalid state. */
	kOutOfRange        generalStatusCode = 3 /**< A value was out of a required range. */
	kBadRequest        generalStatusCode = 4 /**< A request was unrecognized or malformed. */
	kUnsupported       generalStatusCode = 5 /**< An unrecognized or unsupported request was received. */
	kUnexpected        generalStatusCode = 6 /**< A request was not expected at this time. */
	kResourceExhausted generalStatusCode = 7 /**< Insufficient resources to process the given request. */
	kBusy              generalStatusCode = 8 /**< Device is busy and cannot handle this request at this time. */
	kTimeout           generalStatusCode = 9 /**< A timeout occurred. */
	kContinue          generalStatusCode = 1 /**< Context-specific signal to proceed. */
	kAborted           generalStatusCode = 1 /**< Failure, often due to a concurrency error. */
	kInvalidArgument   generalStatusCode = 1 /**< An invalid/unsupported argument was provided. */
	kNotFound          generalStatusCode = 1 /**< Some requested entity was not found. */
	kAlreadyExists     generalStatusCode = 1 /**< The caller attempted to create something that already exists. */
	kPermissionDenied  generalStatusCode = 1 /**< Caller does not have sufficient permissions to execute the requested operations. */
	kDataLoss          generalStatusCode = 1 /**< Unrecoverable data loss or corruption has occurred. */
)

type ErrorType string

var ErrorTimeOut ErrorType = "Time Out"
var ErrorNoMemory ErrorType = "ErrorNoMemory"

func (e ErrorType) Error() string {
	return string(e)
}
