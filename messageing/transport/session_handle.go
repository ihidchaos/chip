package transport

type SessionHandleBase interface {
	Session
	IsActiveSession() bool
}

type SessionHandle struct {
	Session
	mReferenceCounted int
}

func NewSessionHandle(session Session) *SessionHandle {
	return &SessionHandle{
		Session:           session,
		mReferenceCounted: 0,
	}
}
