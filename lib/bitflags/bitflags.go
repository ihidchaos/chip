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

func (b *Flags[T]) Has(flags ...T) bool {
	for _, flag := range flags {
		if flag != b.mFlags&flag {
			return false
		}
	}
	return true
}

func (b *Flags[T]) Sets(isSet bool, flags ...T) {
	if !isSet {
		return
	}
	for _, flag := range flags {
		b.Set(flag)
	}
}

func (b *Flags[T]) Set(f T) {
	b.mFlags = b.mFlags | f
}

func (b *Flags[T]) Unwrap() T {
	return b.mFlags
}
