package bdx

type MsgType uint8

const (
	ProtocolName = "BDX"
	ProtocolId   = 0x0002
)

func (m MsgType) String() string {
	return ""
}
