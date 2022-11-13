package transport

type SessionMessageCounter struct {
	mPeerMessageCounter *PeerMessageCounter
}

func (c SessionMessageCounter) VerifyEncryptedUnicast(counter uint32) error {
	return nil
}

func (c *SessionMessageCounter) PeerMessageCounter() *PeerMessageCounter {
	return c.mPeerMessageCounter
}
