package transport

import (
	"github.com/galenliu/chip/lib"
	"golang.org/x/exp/slices"
	"sync"
	"time"
)

type NewSessionHandlingPolicy uint8

const (
	KShiftToNewSession NewSessionHandlingPolicy = 0
	KStayAtOldSession  NewSessionHandlingPolicy = 1
)

type SessionDelegate interface {
	OnSessionReleased()
}

type SessionBase interface {
	AddHolder(handle *SessionHolder)
	RemoveHolder(holder *SessionHolder)
	SetFabricIndex(index lib.FabricIndex)
	FabricIndex() lib.FabricIndex
	NotifySessionReleased()
	DispatchSessionEvent(delegate SessionDelegate)
}

type Session interface {
	SessionBase
	SessionType() SessionType
	Retain()
	Release()
	IsActiveSession() bool
	IsGroupSession() bool
	IsSecureSession() bool
	IsEstablishing() bool
	ComputeRoundTripTimeout(duration time.Duration) time.Duration
	Released()
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
	locker       sync.Mutex
	mFabricIndex lib.FabricIndex
	mHolders     []*SessionHolder
	*lib.ReferenceCountedHandle
}

func NewSessionBaseImpl(counter int, releaseHandler lib.ReleasedHandler) *SessionBaseImpl {
	return &SessionBaseImpl{
		mFabricIndex:           lib.FabricIndexUndefined,
		mHolders:               make([]*SessionHolder, 0),
		locker:                 sync.Mutex{},
		ReferenceCountedHandle: lib.NewReferenceCountedHandle(counter, releaseHandler),
	}
}

func (s *SessionBaseImpl) FabricIndex() lib.FabricIndex {
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

func (s *SessionBaseImpl) AddHolder(holder *SessionHolder) {
	s.locker.Lock()
	defer s.locker.Unlock()
	s.mHolders = append(s.mHolders, holder)
}

func (s *SessionBaseImpl) RemoveHolder(holder *SessionHolder) {
	s.locker.Lock()
	defer s.locker.Unlock()
	index := slices.Index(s.mHolders, holder)
	if index >= 0 {
		s.mHolders = slices.Delete(s.mHolders, index, index+1)
	}
}
