package transport

type SessionHandle interface {
	IsActiveSession() bool
}

type SessionHandleImpl struct {
	Session
	mReferenceCounted int
}

func (s SessionHandleImpl) IsActiveSession() bool {
	//TODO implement me
	panic("implement me")
}

func NewSessionHandle(session Session) SessionHandleImpl {
	return SessionHandleImpl{
		Session:           session,
		mReferenceCounted: 0,
	}
}
