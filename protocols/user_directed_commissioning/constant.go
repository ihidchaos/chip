package user_directed_commissioning

type MsgType uint8

const (
	IdentificationDeclaration MsgType = 0x00

	ProtocolId   uint16 = 0x0003
	ProtocolName        = "UserDirectedCommissioning"
)

func (m MsgType) String() string {
	switch m {
	case IdentificationDeclaration:
		return "IdentificationDeclaration"
	default:
		return "____"
	}
}
