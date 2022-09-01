package transport

type TypeSecureSession uint8

const (
	PASE TypeSecureSession = 1
	CASE TypeSecureSession = 2
)

func (t TypeSecureSession) Uint8() uint8 {
	return uint8(t)
}

type StateSecureSession uint8

const (
	Establishing StateSecureSession = iota
	Active
	Defunct
	PendingEviction
)

func (t StateSecureSession) Str() string {
	return [...]string{
		"Establishing", "Active", "Defunct", "PendingEviction",
	}[t]
}

func (t StateSecureSession) Uint8() uint8 {
	return uint8(t)
}

type TypeSession uint8

const (
	Undefined       TypeSession = 0
	Unauthenticated TypeSession = 1
	Secure          TypeSession = 2
	GroupIncoming   TypeSession = 3
	GroupOutgoing   TypeSession = 4
)

func (t TypeSession) Str() string {
	return [...]string{
		"Undefined", "Unauthenticated", "Secure", "GroupIncoming", "GroupOutgoing",
	}[t]
}

func (t TypeSession) Uint8() uint8 {
	return uint8(t)
}
