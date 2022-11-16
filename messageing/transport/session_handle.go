package transport

import "github.com/galenliu/chip/messageing/transport/session"

type SessionHandle struct {
	session.Session
}

func NewSessionHandle(session session.Session) *SessionHandle {
	session.Retain()
	return &SessionHandle{
		Session: session,
	}
}
