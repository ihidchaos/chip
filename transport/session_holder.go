package transport

type SessionHolder interface {
	SessionReleased()
	DispatchSessionEvent(delegate SessionDelegate)
}
