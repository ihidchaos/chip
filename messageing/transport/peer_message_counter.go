package transport

type PeerMessageCounter struct {
}

func (c PeerMessageCounter) VerifyUnencrypted(counter uint32) error {
	return nil
}

func (c PeerMessageCounter) CommitUnencrypted(counter uint32) {

}
