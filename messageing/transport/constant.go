package transport

type TSecureSessionType uint8

const (
	K_PASE TSecureSessionType = 1
	K_CASE TSecureSessionType = 2
)

func (t TSecureSessionType) Uint8() uint8 {
	return uint8(t)
}

type TSecureSessionState uint8

const (
	KEstablishing TSecureSessionState = iota
	KActive
	KDefunct
	KPendingEviction
)

func (t TSecureSessionState) Str() string {
	return [...]string{
		"Establishing", "Active", "Defunct", "PendingEviction",
	}[t]
}

func (t TSecureSessionState) Uint8() uint8 {
	return uint8(t)
}

type TSessionType uint8

const (
	kUndefined       TSessionType = 0
	kUnauthenticated TSessionType = 1
	kSecure          TSessionType = 2
	kGroupIncoming   TSessionType = 3
	kGroupOutgoing   TSessionType = 4
)

func (t TSessionType) Str() string {
	return [...]string{
		"Undefined", "kUnauthenticated", "Secure", "GroupIncoming", "GroupOutgoing",
	}[t]
}

func (t TSessionType) Uint8() uint8 {
	return uint8(t)
}
