package transport

import "time"

type SecureSessionImpl struct {
	*SessionBaseImpl
	mTable *SecureSessionTable
}

func NewSecureSessionImpl() *SecureSessionImpl {
	return &SecureSessionImpl{
		SessionBaseImpl: NewSessionBaseImpl(),
		mTable:          NewSecureSessionTable(),
	}
}

func (s SecureSessionImpl) GetSessionType() uint8 {
	return kSessionTypeSecure
}

func (s SecureSessionImpl) GetSessionTypeString() string {
	return "secure"
}

func (s SecureSessionImpl) Retain() {
	//TODO implement me
	panic("implement me")
}

func (s SecureSessionImpl) IsGroupSession() bool {
	//TODO implement me
	panic("implement me")
}

func (s SecureSessionImpl) IsSecureSession() bool {
	return s.GetSessionType() == kSessionTypeSecure
}

func (s SecureSessionImpl) DispatchSessionEvent(delegate SessionDelegate) {
	//TODO implement me
	panic("implement me")
}

func (s SecureSessionImpl) ComputeRoundTripTimeout(duration time.Duration) time.Duration {
	//TODO implement me
	panic("implement me")
}

func (s SecureSessionImpl) SessionReleased() {
	//TODO implement me
	panic("implement me")
}

func (s SecureSessionImpl) AsUnauthenticatedSession() *UnauthenticatedSession {
	//TODO implement me
	panic("implement me")
}

func (s SecureSessionImpl) ClearValue() {
	//TODO implement me
	panic("implement me")
}
