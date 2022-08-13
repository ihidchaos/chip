package transport

type SessionHandle interface {
	Session
}

type SessionHandlerImpl struct {
	*SessionBaseImpl
}
