package transport

type SessionHolder interface {
	SessionReleased()
	ShiftToSession(session SessionHandleBase)
	GrabUnchecked(session SessionHandleBase)
	GrabPairingSession(session SessionHandleBase) bool
	Release()
	DispatchSessionEvent(delegate SessionDelegate)
	Contains(session SessionHandleBase) bool
	SessionDelegate
	SessionHandle() SessionHandleBase
}

type SessionHolderImpl struct {
	mSession          Session
	mReferenceCounted int
}

func (s *SessionHolderImpl) SessionReleased() {
	//TODO implement me
	panic("implement me")
}

func (s *SessionHolderImpl) ShiftToSession(session SessionHandleBase) {
	//TODO implement me
	panic("implement me")
}

func (s *SessionHolderImpl) GrabPairingSession(session SessionHandleBase) bool {
	if !session.IsSecureSession() {
		return false
	}

	if session().IsEstablishing() {
		return false
	}
	s.GrabUnchecked(session)
	return true
}

func (s *SessionHolderImpl) Release() {
	//TODO implement me
	panic("implement me")
}

func (s *SessionHolderImpl) DispatchSessionEvent(delegate SessionDelegate) {
	//TODO implement me
	panic("implement me")
}

func (s *SessionHolderImpl) SessionHandle() SessionHandleBase {
	if s.mSession == nil {
		return nil
	}
	return NewSessionHandle(s.mSession)
}

func NewSessionHolderImpl() *SessionHolderImpl {
	return &SessionHolderImpl{}
}

func (s *SessionHolderImpl) Grad(session SessionHandleBase) bool {
	if !session.IsActiveSession() {
		return false
	}
	s.GrabUnchecked(session)
	return true
}

func (s *SessionHolderImpl) GrabUnchecked(handle SessionHandleBase) {
	handle.AddHolder(s)
}

type SessionHolderWithDelegate interface {
	SessionHolder
}

type SessionHolderWithDelegateImpl struct {
	*SessionHolderImpl
	mDelegate SessionDelegate
}

func NewSessionHolderWithDelegateImpl(delegate SessionDelegate) *SessionHolderWithDelegateImpl {
	return &SessionHolderWithDelegateImpl{
		SessionHolderImpl: NewSessionHolderImpl(),
		mDelegate:         delegate,
	}
}

func (s *SessionHolderWithDelegateImpl) SessionReleased() {
	//TODO implement me
	panic("implement me")
}

func (s *SessionHolderWithDelegateImpl) ShiftToSession(session SessionHandleBase) {
	//TODO implement me
	panic("implement me")
}

func (s *SessionHolderWithDelegateImpl) GrabPairingSession(session SessionHandleBase) {
	//TODO implement me
	panic("implement me")
}

func (s *SessionHolderWithDelegateImpl) Release() {
	//TODO implement me
	panic("implement me")
}

func (s *SessionHolderWithDelegateImpl) DispatchSessionEvent(delegate SessionDelegate) {
	//TODO implement me
	panic("implement me")
}

func (s *SessionHolderImpl) Contains(session SessionHandleBase) bool {
	return s.mSession != nil && session == s.mSession
}
