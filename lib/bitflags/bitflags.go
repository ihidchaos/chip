package bitflags

import (
	"golang.org/x/exp/constraints"
)

type Flags[T constraints.Unsigned] struct {
	mFlags T
}

func Some[T constraints.Unsigned](f T) Flags[T] {
	return Flags[T]{mFlags: f}
}

func (b *Flags[T]) Has(f T) bool {
	if f != b.mFlags&f {
		return false
	}
	return true
}

func (b *Flags[T]) HasAll(flags ...T) bool {
	for _, flag := range flags {
		if !b.Has(flag) {
			return false
		}
	}
	return true
}

func (b *Flags[T]) Set(isSet bool, f T) {
	if !isSet {
		b.Sets(f)
	}
}

func (b *Flags[T]) Sets(flags ...T) {
	for _, flag := range flags {
		b.mFlags = b.mFlags | flag
	}
}

func (b *Flags[T]) Value() T {
	return b.mFlags
}
