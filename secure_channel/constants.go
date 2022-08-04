package secure_channel

// Certificate-based session establishment Message Types
const (
	CASESigma1       uint8 = 0x30
	CASESigma2       uint8 = 0x31
	CASESigma3       uint8 = 0x32
	CASESigma2Resume uint8 = 0x33
)

// Password-based session establishment Message Types
const (
	PBKDFParamRequest  = 0x20
	PBKDFParamResponse = 0x21
	PASEPake1          = 0x22
	PASEPake2          = 0x23
	PASEPake3          = 0x24
)
