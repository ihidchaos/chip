package transport

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
