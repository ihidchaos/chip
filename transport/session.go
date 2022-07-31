package transport

import "github.com/galenliu/chip/lib"

const (
	kSessionTypeUndefined       = 0
	kSessionTypeUnauthenticated = 1
	kSessionTypeSecure          = 2
	kSessionTypeGroupIncoming   = 3
	kSessionTypeGroupOutgoing   = 4
)

type SessionDelegate interface {
}

type Session interface {
	GetFabricIndex()
	SetFabricIndex(index lib.FabricIndex)
	NotifySessionReleased()
	DispatchSessionEvent(delegate SessionDelegate)
	Retain()

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

type SessionImpl struct {
	mFabricIndex lib.FabricId
}
