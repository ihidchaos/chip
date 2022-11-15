package transport

type SecureSessionType uint8
type SecureSessionState uint8
type SessionType uint8

const (
	kEstablishing SecureSessionState = iota
	kActive
	kDefunct
	kPendingEviction
)
const (
	kUndefined SessionType = iota
	kUnauthenticated
	kSecure
	kGroupIncoming
	kGroupOutgoing
)

const (
	SecureSessionTypePASE SecureSessionType = iota
	SecureSessionTypeCASE
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
