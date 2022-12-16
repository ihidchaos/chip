package messageing

const (
	fNone uint16 = 0x0000
	/**< Used to indicate that a response is expected within a specified timeout. */
	fExpectResponse uint16 = 0x0001
	/**< Suppress the auto-request acknowledgment feature when sending a message. */
	fNoAutoRequestAck uint16 = 0x0002
)
