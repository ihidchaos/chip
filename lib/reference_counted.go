package lib

import (
	"errors"
	log "golang.org/x/exp/slog"
	"sync"
)

type ReleasedHandler interface {
	Released()
}

type ReferenceCounted struct {
	locker    sync.Mutex
	mRefCount int
	handler   ReleasedHandler
}

func NewReferenceCounted(initRefCount int, deleter ReleasedHandler) *ReferenceCounted {
	return &ReferenceCounted{mRefCount: initRefCount, handler: deleter, locker: sync.Mutex{}}
}

func (r *ReferenceCounted) Retain() {
	r.locker.Lock()
	defer r.locker.Unlock()
	if r.mRefCount == 0 {
		log.Error("ReferenceCounted Retain", errors.New("ReferenceCounted is zero"))
		return
	}
	r.mRefCount = r.mRefCount + 1
}

func (r *ReferenceCounted) Release() {
	r.locker.Lock()
	defer r.locker.Unlock()
	if r.mRefCount == 0 {
		log.Error("ReferenceCounted Release", errors.New("ReferenceCounted is zero"))
	}
	r.mRefCount = r.mRefCount - 1
	if r.mRefCount == 0 {
		r.handler.Released()
	}
}

func (r *ReferenceCounted) ReferenceCount() int {
	return r.mRefCount
}
