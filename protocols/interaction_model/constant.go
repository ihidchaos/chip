package interaction_model

const ProtocolName = "IM"

type MsgType uint8

const (
	StatusResponse        MsgType = 0x01
	ReadRequest           MsgType = 0x02
	SubscribeRequest      MsgType = 0x03
	SubscribeResponse     MsgType = 0x04
	ReportData            MsgType = 0x05
	WriteRequest          MsgType = 0x06
	WriteResponse         MsgType = 0x07
	InvokeCommandRequest  MsgType = 0x08
	InvokeCommandResponse MsgType = 0x09
	TimedRequest          MsgType = 0x0a
)

func (m MsgType) String() string {
	return ""
}
