package transport

type SessionHolder interface {
	SessionReleased()
	ShiftToSession(session SessionHandle)
	GrabPairingSession(session SessionHandle)
	Release()
	Get() SessionHandle
	DispatchSessionEvent(delegate SessionDelegate)
	Contains(session SessionHandle) bool
	SessionDelegate
	AsSecureSession() *SecureSession
}

type SessionHolderImpl struct {
	mSession Session
}

func NewSessionHolderImpl() *SessionHolderImpl {
	return &SessionHolderImpl{}
}

type SessionHolderWithDelegateImpl struct {
	*SessionHolderImpl
	mDelegate SessionDelegate
}

type SessionHolderWithDelegate interface {
	SessionHolder
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

func (s *SessionHolderWithDelegateImpl) Get() SessionHandle {
	return nil
}

func (s *SessionHolderWithDelegateImpl) DispatchSessionEvent(delegate SessionDelegate) {
	//TODO implement me
	panic("implement me")
}

func (s *SessionHolderWithDelegateImpl) Grad(session SessionHandle) bool {
	if !session.IsActiveSession() {
		return false
	}
	session.AddHolder(s)
	return true
}

func (s *SessionHolderImpl) Contains(session SessionHandle) bool {
	return s.mSession != nil && session == s.mSession
}
