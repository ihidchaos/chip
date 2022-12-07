package session

type SecureSessionType uint8
type SecureState uint8
type Type uint8

type Role uint8

const (
	Initiator Role = iota
	Responder
)

func (t Role) Uint8() uint8 {
	return uint8(t)
}

func (t Role) String() string {
	switch t {
	case Responder:
		return "Responder"
	case Initiator:
		return "Initiator"
	default:
		return "Unknown"
	}
}

const (
	kEstablishing SecureState = iota
	kActive
	kDefunct
	kPendingEviction
)
const (
	kUndefined Type = iota
	kUnauthenticated
	SecureType
	kGroupIncoming
	kGroupOutgoing
)

const (
	SecureSessionTypePASE SecureSessionType = iota
	SecureSessionTypeCASE
)

func (t SecureState) String() string {
	return [...]string{
		"Establishing", "Active", "Defunct", "PendingEviction",
	}[t]
}

func (t SecureSessionType) String() string {
	return [...]string{"Pase", "Case"}[t]
}

func (t Type) String() string {
	return [...]string{
		"Undefined", "Unauthenticated", "Secure", "GroupIncoming", "GroupOutgoing",
	}[t]
}
