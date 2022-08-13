package transport

import (
	"github.com/galenliu/chip/lib"
	"time"
)

type NewSessionHandlingPolicy uint8

const (
	KShiftToNewSession NewSessionHandlingPolicy = 0
	KStayAtOldSession  NewSessionHandlingPolicy = 1
)

type SessionDelegate interface {
}

type SessionBase interface {
	AddHolder(handle SessionHolder)
	RemoveHolder(holder SessionHolder)
	SetFabricIndex(index lib.FabricIndex)
	GetFabricIndex() lib.FabricIndex
	NotifySessionReleased()
	DispatchSessionEvent(delegate SessionDelegate)
}

type Session interface {
	SessionBase
	GetSessionType() uint8
	GetSessionTypeString() string
	AddHolder(s SessionHolder)

	//NotifySessionReleased()

	Retain()
	Release()
	IsActiveSession() bool

	IsGroupSession() bool
	IsSecureSession() bool
	IsEstablishing() bool

	ComputeRoundTripTimeout(duration time.Duration) time.Duration

	SessionReleased()
	AsUnauthenticatedSession() *UnauthenticatedSessionImpl
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

func NewSessionBaseImpl() *SessionBaseImpl {
	return &SessionBaseImpl{
		mFabricIndex: lib.UndefinedFabricIndex,
		mHolders:     make([]SessionHolder, 0),
	}
}

func (s *SessionBaseImpl) GetFabricIndex() lib.FabricIndex {
	return s.mFabricIndex
}

func (s *SessionBaseImpl) NotifySessionReleased() {
	//TODO implement me
	panic("implement me")
}

func (s *SessionBaseImpl) SetFabricIndex(index lib.FabricIndex) {
	s.mFabricIndex = index
}

func (s *SessionBaseImpl) DispatchSessionEvent(event SessionDelegate) {
	for _, holder := range s.mHolders {
		holder.DispatchSessionEvent(event)
	}
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
