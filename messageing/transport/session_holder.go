package transport

import (
	"github.com/galenliu/chip/lib"
	"github.com/galenliu/chip/messageing/transport/session"
)

type SessionHolder struct {
	session.Session
	*lib.ReferenceCounted
}

func NewSessionHolder(session session.Session) *SessionHolder {
	session.Retain()
	return &SessionHolder{
		Session: session,
	}
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

func (s *SessionHolder) SessionHandler() *SessionHandle {
	if s.Session != nil {
		return NewSessionHandle(s.Session)
	}
	return nil
}

func (s *SessionHolder) Release() {
	if s.Session != nil {
		s.Session.RemoveHolder(s)
	}
	s.Session = nil
}

func (s *SessionHolder) GrabPairingSession(ss *SessionHandle) bool {
	s.Release()
	if !ss.IsSecure() {
		return false
	}
	secureSession, ok := ss.Session.(*session.Secure)
	if !ok || !secureSession.IsEstablishing() {
		return false
	}
	s.GrabUnchecked(ss)
	return true
}

func (s *SessionHolder) Grad(session *SessionHandle) bool {
	s.Release()
	if !session.IsActive() {
		return false
	}
	s.GrabUnchecked(session)
	return true
}

func (s *SessionHolder) GrabUnchecked(session *SessionHandle) {
	s.Session = session.Session
	session.AddHolder(s)
}

type SessionHolderWithDelegate struct {
	*SessionHolder
	mDelegate session.Delegate
}

func NewSessionHolderWithDelegateImpl(delegate session.Delegate) *SessionHolderWithDelegate {
	return &SessionHolderWithDelegate{
		SessionHolder: &SessionHolder{},
		mDelegate:     delegate,
	}
}

func (s *SessionHolderWithDelegate) DispatchSessionEvent(event session.DelegateEvent) {
	event()
}

func (s *SessionHolderWithDelegate) SessionReleased() {
	//TODO implement me
	panic("implement me")
}

func (s *SessionHolderWithDelegate) ShiftToSession(session SessionHandle) {
	//TODO implement me
	panic("implement me")
}

func (s *SessionHolderWithDelegate) Release() {
	//TODO implement me
	panic("implement me")
}
