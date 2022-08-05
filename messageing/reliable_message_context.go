package messageing

const (
	kFlagInitiator                = 1 << iota
	kFlagResponseExpected         = 1 << 1
	kFlagAutoRequestAck           = 1 << 2
	kFlagMessageNotAcked          = 1 << 3
	kFlagAckPending               = 1 << 4
	kFlagAckMessageCounterIsValid = 1 << 5
	kFlagWillSendMessage          = 1 << 6
	kFlagClosed                   = 1 << 7
	kFlagActiveMode               = 1 << 8
	kFlagEphemeralExchange        = 1 << 9
	kFlagIgnoreSessionRelease     = 1 << 10
)

type ReliableMessageContext struct {
	mFlags uint16
}
