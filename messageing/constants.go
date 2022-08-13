package messageing

// Opcode Password-based session establishment Message Types
// Protocol Opcode
type Opcode uint8

const (
	MsgCounterSyncReq Opcode = 0x00 // Message Counter Synchronization Protocol Message Types
	MsgCounterSyncRsp Opcode = 0x01
	StandaloneAck     Opcode = 0x10
	StatusReport      Opcode = 0x40
)

// Password-based session establishment Message Types
const (
	PBKDFParamRequest  Opcode = 0x20
	PBKDFParamResponse Opcode = 0x21
	PASEPake1          Opcode = 0x22
	PASEPake2          Opcode = 0x23
	PASEPake3          Opcode = 0x24
)

// Certificate-based session establishment Message Types
const (
	CASESigma1       Opcode = 0x30
	CASESigma2       Opcode = 0x31
	CASESigma3       Opcode = 0x32
	CASESigma2Resume Opcode = 0x33
)
