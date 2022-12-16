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

func (b *Flags[T]) Raw() T {
	return b.mFlags
}

func (b *Flags[T]) Clear(f T) {

}

func U64To32(val uint64) (h, l uint32) {
	h = uint32((val & 0xFFFF_FFFF_0000_0000) >> 32)
	l = uint32(val & 0x0000_0000_FFFF_FFFF)
	return
}

func U32To16(val uint32) (h, l uint16) {
	h = uint16((val & 0xFFFF_0000) >> 16)
	l = uint16(val & 0x0000_FFFF)
	return
}
