package transport

type SecureSessionType uint8
type SecureSessionState uint8
type SessionType uint8

const (
	Establishing SecureSessionState = iota
	Active
	Defunct
	PendingEviction

	PASE SecureSessionType = iota - 4
	CASE

	Undefined SessionType = iota - 2
	Unauthenticated
	Secure
	GroupIncoming
	GroupOutgoing
)

func (t SecureSessionState) String() string {
	return [...]string{
		"Establishing", "Active", "Defunct", "PendingEviction",
	}[t]
}

func (t SecureSessionType) String() string {
	return [...]string{"Pase", "Case"}[t]
}

func (t SessionType) String() string {
	return [...]string{
		"Undefined", "Unauthenticated", "Secure", "GroupIncoming", "GroupOutgoing",
	}[t]
}
