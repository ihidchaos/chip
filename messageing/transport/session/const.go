package session

type SecureType uint8
type State uint8
type Type uint8

type Role uint8

const (
	RoleInitiator Role = iota
	RoleResponder
)

func (t Role) Uint8() uint8 {
	return uint8(t)
}

func (t Role) String() string {
	switch t {
	case RoleResponder:
		return "RoleResponder"
	case RoleInitiator:
		return "RoleInitiator"
	default:
		return "Unknown"
	}
}

const (
	kEstablishing State = iota
	kActive
	kDefunct
	kPendingEviction
)
const (
	kUndefined Type = iota
	kUnauthenticated
	kSecure
	kGroupIncoming
	kGroupOutgoing
)

const (
	SecureSessionTypePASE SecureType = iota
	SecureSessionTypeCASE
)

func (t State) String() string {
	return [...]string{
		"Establishing", "Active", "Defunct", "PendingEviction",
	}[t]
}

func (t SecureType) String() string {
	return [...]string{"Pase", "Case"}[t]
}

func (t Type) String() string {
	return [...]string{
		"Undefined", "Unauthenticated", "Secure", "GroupIncoming", "GroupOutgoing",
	}[t]
}

type ErrorType string

var ErrorMessageCounterExhausted ErrorType = "MessageCounterExhausted"

func (e ErrorType) Error() string {
	return string(e)
}
