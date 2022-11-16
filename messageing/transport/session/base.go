package session

import (
	"github.com/galenliu/chip/lib"
	"github.com/galenliu/chip/messageing/transport/raw"
	"golang.org/x/exp/slices"
	log "golang.org/x/exp/slog"
	"sync"
	"time"
)

type NewSessionHandlingPolicy uint8

const (
	KShiftToNewSession NewSessionHandlingPolicy = 0
	KStayAtOldSession  NewSessionHandlingPolicy = 1
	kMinActiveTime     time.Duration            = time.Duration(4 * time.Second)
)

type Delegate interface {
	OnSessionReleased()
}

type Holder interface {
	DispatchSessionEvent(Delegate)
}

type Base interface {
	AddHolder(handle Holder)
	RemoveHolder(holder Holder)
	SetFabricIndex(index lib.FabricIndex)
	FabricIndex() lib.FabricIndex
	NotifySessionReleased()
	DispatchSessionEvent(delegate Delegate)
	IsSecureSession() bool
	IsGroupSession() bool
	SessionType() Type
}

type Session interface {
	Base
	Retain()
	Release()
	IsActiveSession() bool
	IsEstablishing() bool
	AckTimeout() time.Duration
	GetPeer() lib.ScopedNodeId
	ComputeRoundTripTimeout(duration time.Duration) time.Duration
	RemoteMRPConfig() *ReliableMessageProtocolConfig
	Released()
	ClearValue()
}

type BaseImpl struct {
	locker       sync.Mutex
	mFabricIndex lib.FabricIndex
	mHolders     []Holder
	mSessionType Type
	mPeerAddress raw.PeerAddress
	mDelegate    Session
	*lib.ReferenceCounted
}

func NewBaseImpl(initCounter int, sessionType Type, delegate Session) *BaseImpl {
	return &BaseImpl{
		mFabricIndex:     lib.UndefinedFabricIndex(),
		locker:           sync.Mutex{},
		mSessionType:     sessionType,
		mDelegate:        delegate,
		ReferenceCounted: lib.NewReferenceCounted(initCounter, delegate),
	}
}

func (s *BaseImpl) FabricIndex() lib.FabricIndex {
	return s.mFabricIndex
}

func (s *BaseImpl) NotifySessionReleased() {
	//TODO implement me
	panic("implement me")
}

func (s *BaseImpl) SetFabricIndex(index lib.FabricIndex) {
	s.mFabricIndex = index
}

func (s *BaseImpl) DispatchSessionEvent(event Delegate) {
	for _, holder := range s.mHolders {
		holder.DispatchSessionEvent(event)
	}
}

func (s *BaseImpl) AddHolder(holder Holder) {
	s.locker.Lock()
	defer s.locker.Unlock()
	if s.mHolders == nil {
		s.mHolders = make([]Holder, 0)
	}
	s.mHolders = append(s.mHolders, holder)
}

func (s *BaseImpl) RemoveHolder(holder Holder) {
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

func (s *BaseImpl) IsSecureSession() bool {
	return s.SessionType() == kSecure
}

func (s *BaseImpl) IsGroupSession() bool {
	return s.SessionType() == kGroupIncoming || s.SessionType() == kGroupOutgoing
}

func (s *BaseImpl) SessionType() Type {
	return s.mSessionType
}

func (s *BaseImpl) LogValue() log.Value {
	return log.GroupValue(
		log.String("Type", s.mSessionType.String()),
	)
}
