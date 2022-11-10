package lib

import (
	"errors"
	log "golang.org/x/exp/slog"
	"sync"
)

type ReleasedHandler interface {
	Released()
}

type ReferenceCountedHandle struct {
	locker      sync.Mutex
	mRefCounted int
	handler     ReleasedHandler
}

func NewReferenceCountedHandle(initRefCount int, deleter ReleasedHandler) *ReferenceCountedHandle {
	return &ReferenceCountedHandle{mRefCounted: initRefCount, handler: deleter, locker: sync.Mutex{}}
}

func (r *ReferenceCountedHandle) Retain() {
	r.locker.Lock()
	defer r.locker.Unlock()
	if r.mRefCounted == 0 {
		log.Error("ReferenceCountedHandle Retain", errors.New("ReferenceCounted is zero"))
		return
	}
	r.mRefCounted = r.mRefCounted + 1
}

func (r *ReferenceCountedHandle) Release() {
	r.locker.Lock()
	defer r.locker.Unlock()
	if r.mRefCounted == 0 {
		log.Error("ReferenceCountedHandle Release", errors.New("ReferenceCounted is zero"))
	}
	r.mRefCounted = r.mRefCounted - 1
	if r.mRefCounted == 0 {
		r.handler.Released()
	}
}

type ReferenceCounted struct {
	locker      sync.Mutex
	mRefCounted int
	delegate    ReleasedHandler
}

func NewReferenceCounted(initRefCount int, deleter ReleasedHandler) *ReferenceCounted {
	return &ReferenceCounted{mRefCounted: initRefCount, delegate: deleter, locker: sync.Mutex{}}
}
