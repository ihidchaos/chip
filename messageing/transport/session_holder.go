package transport

type SessionHolder struct {
	Session           *Session
	mReferenceCounted int
}

func (s *SessionHolder) SessionReleased() {
	//TODO implement me
	panic("implement me")
}

func (s *SessionHolder) ShiftToSession(session SessionHandleBase) {
	//TODO implement me
	panic("implement me")
}

//func (s *SessionHolderImpl) GrabPairingSession(session SessionHandleBase) bool {
//	if !session.IsSecureSession() {
//		return false
//	}
//
//	if session.IsEstablishing() {
//		return false
//	}
//	s.GrabUnchecked(session)
//	return true
//}

func (s *SessionHolder) Release() {
	//TODO implement me
	panic("implement me")
}

func (s *SessionHolder) DispatchSessionEvent(delegate SessionDelegate) {
	//TODO implement me
	panic("implement me")
}

//func (s *SessionHolderImpl) Grad(session SessionHandleBase) bool {
//	if !session.IsActiveSession() {
//		return false
//	}
//	s.GrabUnchecked(session)
//	return true
//}
//
//func (s *SessionHolderImpl) GrabUnchecked(handle SessionHandleBase) {
//	handle.AddHolder(s)
//}

type SessionHolderWithDelegate struct {
	*SessionHolder
	mDelegate SessionDelegate
}

func NewSessionHolderWithDelegateImpl(delegate SessionDelegate) *SessionHolderWithDelegate {
	return &SessionHolderWithDelegate{
		SessionHolder: &SessionHolder{},
		mDelegate:     delegate,
	}
}

func (s *SessionHolderWithDelegate) SessionReleased() {
	//TODO implement me
	panic("implement me")
}

func (s *SessionHolderWithDelegate) ShiftToSession(session SessionHandleBase) {
	//TODO implement me
	panic("implement me")
}

func (s *SessionHolderWithDelegate) GrabPairingSession(session SessionHandleBase) {
	//TODO implement me
	panic("implement me")
}

func (s *SessionHolderWithDelegate) Release() {
	//TODO implement me
	panic("implement me")
}

func (s *SessionHolderWithDelegate) DispatchSessionEvent(delegate SessionDelegate) {
	//TODO implement me
	panic("implement me")
}

//func (s *SessionHolder) Contains(session SessionHandleBase) bool {
//	return s.Session != nil && session == s.Session
//}
