package transport

import (
	"github.com/galenliu/chip/lib"
	"time"
)

const (
	kSessionTypeUndefined       uint8 = 0
	kSessionTypeUnauthenticated uint8 = 1
	kSessionTypeSecure          uint8 = 2
	kSessionTypeGroupIncoming   uint8 = 3
	kSessionTypeGroupOutgoing   uint8 = 4
)

type SessionDelegate interface {
}

type SessionBase interface {
	AddHolder(handle SessionHolder)
	RemoveHolder(holder SessionHolder)
	SetFabricIndex(index lib.FabricIndex)
	GetFabricIndex() lib.FabricIndex
	NotifySessionReleased()
}

type Session interface {
	GetSessionType() uint8
	GetSessionTypeString() string

	NotifySessionReleased()

	//NotifySessionReleased()

	Retain()

	IsGroupSession() bool
	IsSecureSession() bool

	DispatchSessionEvent(delegate SessionDelegate)
	ComputeRoundTripTimeout(duration time.Duration) time.Duration

	SessionReleased()
	AsUnauthenticatedSession() *UnauthenticatedSession
	ClearValue()

	//virtual void Retain()  = 0;
	//virtual void Release() = 0;
	//virtual bool IsActiveSession() const = 0;
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
	mHolders     []SessionHolder
}

func (s *SessionBaseImpl) GetFabricIndex() lib.FabricIndex {
	return s.mFabricIndex
}

func NewSessionBaseImpl() *SessionBaseImpl {
	return &SessionBaseImpl{
		mFabricIndex: lib.UndefinedFabricIndex,
		mHolders:     make([]SessionHolder, 0),
	}
}

func (s *SessionBaseImpl) NotifySessionReleased() {
	//TODO implement me
	panic("implement me")
}

func (s *SessionBaseImpl) SetFabricIndex(index lib.FabricIndex) {
	s.mFabricIndex = index
}

func (s *SessionBaseImpl) AddHolder(holder SessionHolder) {
	s.mHolders = append(s.mHolders, holder)
}

func (s *SessionBaseImpl) RemoveHolder(holder SessionHolder) {
	for i, h := range s.mHolders {
		if h == holder {
			s.mHolders = append(s.mHolders[:i], s.mHolders[i+1:]...)
		}
	}
}