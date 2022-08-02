package secure_channel

// Certificate-based session establishment Message Types
const (
	CaseSigma1       uint8 = 0x30
	CaseSigma2       uint8 = 0x31
	CaseSigma3       uint8 = 0x32
	CaseSigma2resume uint8 = 0x33
)

// Password-based session establishment Message Types
const (
	PBKDFParamRequest  = 0x20
	PBKDFParamResponse = 0x21
	PASE_Pake1         = 0x22
	PASE_Pake2         = 0x23
	PASE_Pake3         = 0x24
)
