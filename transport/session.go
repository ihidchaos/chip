package transport

import (
	"github.com/galenliu/chip/lib"
	"time"
)

const (
	kSessionTypeUndefined       = 0
	kSessionTypeUnauthenticated = 1
	kSessionTypeSecure          = 2
	kSessionTypeGroupIncoming   = 3
	kSessionTypeGroupOutgoing   = 4
)

type SessionDelegate interface {
}

type SessionBase interface {
}

type Session interface {
	NotifySessionReleased()
	GetFabricIndex()
	SetFabricIndex(index lib.FabricIndex)
	//NotifySessionReleased()

	Retain()

	IsGroupSession() bool
	IsSecureSession() bool
	DispatchSessionEvent(delegate SessionDelegate)
	ComputeRoundTripTimeout(duration time.Duration) time.Duration
	AddHolder(handle SessionHandle)
	RemoveHolder(handle SessionHandle)

	//	virtual void Retain()  = 0;
	//virtual void Release() = 0;
	//
	//virtual bool IsActiveSession() const = 0;
	//
	//virtual ScopedNodeId GetPeer() const                                     = 0;
	//virtual ScopedNodeId GetLocalScopedNodeId() const                        = 0;
	//virtual Access::SubjectDescriptor GetSubjectDescriptor() const           = 0;
	//virtual bool RequireMRP() const                                          = 0;
	//virtual const ReliableMessageProtocolConfig & GetRemoteMRPConfig() const = 0;
	//virtual System::Clock::Timestamp GetMRPBaseTimeout()                     = 0;
	//virtual System::Clock::Milliseconds32 GetAckTimeout() const              = 0;
}

type SessionBaseImpl struct {
	mFabricIndex lib.FabricIndex
	mHolders     []SessionHandle
}

func (s *SessionBaseImpl) NotifySessionReleased() {
	for _, handler := range s.mHolders {
		handler.SessionReleased()
	}
}

func (s *SessionBaseImpl) SetFabricIndex(index lib.FabricIndex) {
	s.mFabricIndex = index
}

func (s *SessionBaseImpl) DispatchSessionEvent(delegate SessionDelegate) {
	for _, handler := range s.mHolders {
		handler.DispatchSessionEvent(delegate)
	}
}

func (s SessionBaseImpl) ComputeRoundTripTimeout(duration time.Duration) time.Duration {
	//TODO implement me
	panic("implement me")
}

func (s SessionBaseImpl) AddHolder(handle SessionHandle) {
	//TODO implement me
	panic("implement me")
}

func (s SessionBaseImpl) RemoveHolder(handle SessionHandle) {
	//TODO implement me
	panic("implement me")
}
