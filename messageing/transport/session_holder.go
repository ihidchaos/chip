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
}

type SessionHolderImpl struct {
	mSession Session
}

func NewSessionHolderImpl() *SessionHolderImpl {
	return &SessionHolderImpl{}
}

type SessionHolderWithDelegate interface {
	SessionHolder
}

func NewSessionHolderWithDelegateImpl() *SessionHolderWithDelegateImpl {
	return &SessionHolderWithDelegateImpl{}
}

type SessionHolderWithDelegateImpl struct {
}

func (s SessionHolderWithDelegateImpl) SessionReleased() {
	//TODO implement me
	panic("implement me")
}

func (s SessionHolderWithDelegateImpl) ShiftToSession(session SessionHandle) {
	//TODO implement me
	panic("implement me")
}

func (s SessionHolderWithDelegateImpl) GrabPairingSession(session SessionHandle) {
	//TODO implement me
	panic("implement me")
}

func (s SessionHolderWithDelegateImpl) Release() {
	//TODO implement me
	panic("implement me")
}

func (s SessionHolderWithDelegateImpl) Get() SessionHandle {
	return nil
}

func (s SessionHolderWithDelegateImpl) DispatchSessionEvent(delegate SessionDelegate) {
	//TODO implement me
	panic("implement me")
}

func (s SessionHolderWithDelegateImpl) Contains(session SessionHandle) bool {
	//TODO implement me
	panic("implement me")
}
