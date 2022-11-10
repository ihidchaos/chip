package transport

import "github.com/galenliu/chip/lib"

type SessionHolder struct {
	Session
	*lib.ReferenceCountedHandle
}

func (s *SessionHolder) SessionReleased() {
	//TODO implement me
	panic("implement me")
}

func (s *SessionHolder) ShiftToSession(session *SessionHandle) {
	//TODO implement me
	panic("implement me")
}

func (s *SessionHolder) Contains(session *SessionHandle) bool {
	return s.Session != nil && s.Session == session.Session
}

func (s *SessionHolder) Get() *SessionHandle {
	if s.Session != nil {
		return NewSessionHandle(s.Session)
	}
	return nil
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
	if s.Session != nil {
		s.Session.RemoveHolder(s)
		s.Session.ClearValue()
	}
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

func (s *SessionHolderWithDelegate) ShiftToSession(session SessionHandle) {
	//TODO implement me
	panic("implement me")
}

func (s *SessionHolderWithDelegate) GrabPairingSession(session SessionHandle) {
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
