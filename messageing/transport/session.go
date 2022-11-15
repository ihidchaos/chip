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
	IsSecureSession() bool
	SessionType() SessionType
}

type Session interface {
	SessionBase
	Retain()
	Release()
	IsActiveSession() bool
	IsGroupSession() bool
	IsEstablishing() bool
	GetPeer() lib.ScopedNodeId
	ComputeRoundTripTimeout(duration time.Duration) time.Duration
	Released()
	ClearValue()
}

type SessionBaseImpl struct {
	locker       sync.Mutex
	mFabricIndex lib.FabricIndex
	mHolders     []*SessionHolder
	mSessionType SessionType
	*lib.ReferenceCounted
}

func NewSessionBaseImpl(initCounter int, sessionType SessionType, releaseHandler lib.ReleasedHandler) *SessionBaseImpl {
	return &SessionBaseImpl{
		mFabricIndex:     lib.FabricIndexUndefined,
		locker:           sync.Mutex{},
		mSessionType:     sessionType,
		ReferenceCounted: lib.NewReferenceCounted(initCounter, releaseHandler),
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
	if s.mHolders == nil {
		s.mHolders = make([]*SessionHolder, 0)
	}
	s.mHolders = append(s.mHolders, holder)
}

func (s *SessionBaseImpl) RemoveHolder(holder *SessionHolder) {
	s.locker.Lock()
	defer s.locker.Unlock()
	if s.mHolders == nil || len(s.mHolders) == 0 {
		return
	}
	index := slices.Index(s.mHolders, holder)
	if index >= 0 {
		s.mHolders = slices.Delete(s.mHolders, index, index+1)
	}
}

func (s *SessionBaseImpl) IsSecureSession() bool {
	return s.SessionType() == kSecure
}

func (s *SessionBaseImpl) SessionType() SessionType {
	return s.mSessionType
}
