package echo

type MsgType uint8

const (
	EchoRequest  MsgType = 0x01
	EchoResponse MsgType = 0x02

	ProtocolId   uint16 = 0x0004
	ProtocolName        = "Echo"
)

func (m MsgType) String() string {
	switch m {
	case EchoRequest:
		return "EchoRequest"
	case EchoResponse:
		return "EchoResponse"
	default:
		return "_____"
	}
}
