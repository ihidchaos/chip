package session

type MessageCounter struct {
	MessageCounter     *LocalMessageCounter
	PeerMessageCounter *PeerMessageCounter
}
