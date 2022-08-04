package transport

type SessionHolder interface {
	SessionReleased()
	ShiftToSession(session SessionHandle)
	GrabPairingSession(session SessionHandle)
	Release()
	Get() SessionHandle
	DispatchSessionEvent(delegate SessionDelegate)
}

type SessionHolderWithDelegate interface {
	SessionHolder
}
