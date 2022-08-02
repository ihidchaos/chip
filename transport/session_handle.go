package transport

import "net/netip"

type SessionHandle interface {
	SessionReleased()
	DispatchSessionEvent(delegate SessionDelegate)

	AsUnauthenticatedSession() *UnauthenticatedSession
	ClearValue()
}

type SessionHandleImpl struct {
	mDelegate   SessionDelegate
	peerAddress netip.AddrPort
	Session     Session
}

func (s SessionHandleImpl) SessionReleased() {
	//TODO implement me
	panic("implement me")
}

func (s SessionHandleImpl) ClearValue() {
	//TODO implement me
	panic("implement me")
}

func (s SessionHandleImpl) DispatchSessionEvent(delegate SessionDelegate) {
	//TODO implement me
	panic("implement me")
}

func (s SessionHandleImpl) AsUnauthenticatedSession() *UnauthenticatedSession {
	session, _ := s.Session.(*UnauthenticatedSession)
	return session
}
