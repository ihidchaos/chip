package messageing

// MsgType Password-based session establishment Message Types
// Protocol MsgType
type MsgType uint8

const (
	MsgCounterSyncReq MsgType = 0x00 // Message Counter Synchronization Protocol Message Types
	MsgCounterSyncRsp MsgType = 0x01
	StandaloneAck     MsgType = 0x10
	StatusReport      MsgType = 0x40
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
)
