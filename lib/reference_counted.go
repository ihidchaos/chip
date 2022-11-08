package lib

import "log"

type DeleteDelegate interface {
	DoRelease()
}

type ReferenceCounted struct {
	mRefCounted int
	deleter     DeleteDelegate
}

func NewReferenceCounted(initRefCount int, deleter DeleteDelegate) *ReferenceCounted {
	return &ReferenceCounted{mRefCounted: initRefCount, deleter: deleter}
}

func (r *ReferenceCounted) Retain() {
	if r.mRefCounted == 0 {
		log.Panicln("ReferenceCounted error")
	}
	r.mRefCounted = r.mRefCounted + 1
}

func (r *ReferenceCounted) Release() {
	if r.mRefCounted == 0 {
		log.Panicln("ReferenceCounted error")
	}
	r.mRefCounted = r.mRefCounted - 1
	if r.mRefCounted == 0 {
		r.deleter.DoRelease()
	}
}
