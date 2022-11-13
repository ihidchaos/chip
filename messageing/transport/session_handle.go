package transport

import "github.com/galenliu/chip/lib"

type SessionHandle struct {
	Session
	*lib.ReferenceCounted
}

func NewSessionHandle(session Session) *SessionHandle {
	return &SessionHandle{
		Session:          session,
		ReferenceCounted: lib.NewReferenceCounted(1, nil),
	}
}
