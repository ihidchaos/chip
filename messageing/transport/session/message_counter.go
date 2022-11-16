package session

type MessageCounter struct {
	mPeerMessageCounter *PeerMessageCounter
}

func (c *MessageCounter) VerifyEncryptedUnicast(counter uint32) error {
	return nil
}

func (c *MessageCounter) PeerMessageCounter() *PeerMessageCounter {
	return c.mPeerMessageCounter
}
