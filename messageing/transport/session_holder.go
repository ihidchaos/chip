package transport

type SessionHolder interface {
	SessionReleased()
	ShiftToSession(session SessionHandle)
	GrabUnchecked(session SessionHandle)
	GrabPairingSession(session SessionHandle) bool
	Release()
	DispatchSessionEvent(delegate SessionDelegate)
	Contains(session SessionHandle) bool
	SessionDelegate
	SessionHandle() SessionHandle
}

type SessionHolderImpl struct {
	mSession          Session
	mReferenceCounted int
}

func (s *SessionHolderImpl) SessionReleased() {
	//TODO implement me
	panic("implement me")
}

func (s *SessionHolderImpl) ShiftToSession(session SessionHandle) {
	//TODO implement me
	panic("implement me")
}

func (s *SessionHolderImpl) GrabPairingSession(session SessionHandle) bool {
	if !session.Session.IsSecureSession() {
		return false
	}

	if session.Session().IsEstablishing() {
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

func (s *SessionHolderImpl) SessionHandle() SessionHandle {
	if s.mSession == nil {
		return nil
	}
	return NewSessionHandle(s.mSession)
}

func NewSessionHolderImpl() *SessionHolderImpl {
	return &SessionHolderImpl{}
}

func (s *SessionHolderImpl) Grad(session SessionHandle) bool {
	if !session.IsActiveSession() {
		return false
	}
	s.GrabUnchecked(session)
	return true
}

func (s *SessionHolderImpl) GrabUnchecked(handle SessionHandle) {
	handle.Session().AddHolder(s)
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

func (s *SessionHolderWithDelegateImpl) AsSecureSession() *SecureSession {
	ss, ok := s.mSession.(*SecureSession)
	if ok {
		return ss
	}
	return nil
}

func (s *SessionHolderWithDelegateImpl) SessionReleased() {
	//TODO implement me
	panic("implement me")
}

func (s *SessionHolderWithDelegateImpl) ShiftToSession(session SessionHandle) {
	//TODO implement me
	panic("implement me")
}

func (s *SessionHolderWithDelegateImpl) GrabPairingSession(session SessionHandle) {
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

func (s *SessionHolderImpl) Contains(session SessionHandle) bool {
	return s.mSession != nil && session.Session() == s.mSession
}
