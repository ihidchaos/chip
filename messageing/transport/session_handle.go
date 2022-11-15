package transport

type SessionHandle struct {
	Session
}

func NewSessionHandle(session Session) *SessionHandle {
	session.Retain()
	return &SessionHandle{
		Session: session,
	}
}
