package transport

type NewSessionHandlingPolicy uint8

const (
	ShiftToNewSession NewSessionHandlingPolicy = 0
	StayAtOldSession  NewSessionHandlingPolicy = 1
)

type SessionDelegate interface {
	OnSessionReleased()
	OnSessionHang()
	GetNewSessionHandlingPolicy() NewSessionHandlingPolicy
}
