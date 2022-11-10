package transport

import "github.com/galenliu/chip/lib"

type SessionHandle struct {
	Session
	*lib.ReferenceCountedHandle
}

func NewSessionHandle(session Session) *SessionHandle {
	return &SessionHandle{
		Session:                session,
		ReferenceCountedHandle: lib.NewReferenceCountedHandle(1, nil),
	}
}
