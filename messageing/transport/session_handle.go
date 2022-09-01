package transport

type SessionHandle interface {
	IsActiveSession() bool
	Session() Session
}

type SessionHandleImpl struct {
	mSession          Session
	mReferenceCounted int
}

func (s SessionHandleImpl) IsActiveSession() bool {
	//TODO implement me
	panic("implement me")
}

func (s SessionHandleImpl) Session() Session {
	return s.mSession
}

func NewSessionHandle(session Session) SessionHandleImpl {
	return SessionHandleImpl{
		mSession:          session,
		mReferenceCounted: 0,
	}
}
