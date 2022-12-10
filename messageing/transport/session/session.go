package session

import (
	"github.com/galenliu/chip/lib"
	"github.com/galenliu/chip/messageing"
	"github.com/galenliu/chip/messageing/transport/raw"
	log "golang.org/x/exp/slog"
	"sync"
	"time"
)

type DelegateEvent func()

type Delegate interface {
	OnSessionReleased()
}

type Holder interface {
	DispatchSessionEvent(DelegateEvent)
}

type Session interface {
	Retain()
	Release()
	IsActive() bool
	IsEstablishing() bool
	AckTimeout() time.Duration
	GetPeer() lib.ScopedNodeId
	ComputeRoundTripTimeout(duration time.Duration) time.Duration
	RemoteMRPConfig() *messageing.ReliableMessageProtocolConfig
	Released()
	ClearValue()
	AddHolder(handle Holder)
	RemoveHolder(holder Holder)
	FabricIndex() lib.FabricIndex
	NotifySessionReleased()
	DispatchSessionEvent(delegate DelegateEvent)
	IsSecure() bool
	IsGroup() bool
	Type() Type
}

type BaseImpl struct {
	locker       sync.Mutex
	mFabricIndex lib.FabricIndex
	mHolders     []Holder
	mSessionType Type
	mPeerAddress raw.PeerAddress
	base         Session
	*lib.ReferenceCounted
}

func (s *BaseImpl) FabricIndex() lib.FabricIndex {
	return s.mFabricIndex
}

func (s *BaseImpl) NotifySessionReleased() {
	//TODO implement me
	panic("implement me")
}

func (s *BaseImpl) DispatchSessionEvent(event DelegateEvent) {
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
	for i, h := range s.mHolders {
		if h == holder {
			s.mHolders = append(s.mHolders[:i], s.mHolders[i+1:]...)
		}
	}
}

func (s *BaseImpl) IsSecure() bool {
	return s.Type() == kSecure
}

func (s *BaseImpl) IsGroup() bool {
	return s.Type() == kGroupIncoming || s.Type() == kGroupOutgoing
}

func (s *BaseImpl) Type() Type {
	return s.mSessionType
}

func (s *BaseImpl) LogValue() log.Value {
	return log.GroupValue(
		log.String("Type", s.mSessionType.String()),
		log.Int("ReferenceCounted", s.ReferenceCount()),
	)
}
