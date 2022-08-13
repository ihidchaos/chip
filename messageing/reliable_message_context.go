package messageing

const (
	kFlagInitiator                uint16 = 1 << iota
	kFlagResponseExpected         uint16 = 1 << 1
	kFlagAutoRequestAck           uint16 = 1 << 2
	kFlagMessageNotAcked          uint16 = 1 << 3
	kFlagAckPending               uint16 = 1 << 4
	kFlagAckMessageCounterIsValid uint16 = 1 << 5
	kFlagWillSendMessage          uint16 = 1 << 6
	kFlagClosed                   uint16 = 1 << 7
	kFlagActiveMode               uint16 = 1 << 8
	kFlagEphemeralExchange        uint16 = 1 << 9
	kFlagIgnoreSessionRelease     uint16 = 1 << 10
)

type ReliableMessageContext struct {
	mFlags uint16
}
