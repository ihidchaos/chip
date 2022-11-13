package lib

import "golang.org/x/exp/constraints"

// HasFlags Value中所有Flags是否存在
func HasFlags[T constraints.Unsigned](value T, flags ...T) (b bool) {
	b = true
	for _, f := range flags {
		if value&f == 0 {
			b = false
		}
	}
	return
}

func SetFlags[T constraints.Unsigned](value T, flags ...T) T {
	for _, f := range flags {
		value = f | value
	}
	return value
}

func SetFlag[T constraints.Unsigned](isSet bool, value T, flag T) T {
	if isSet {
		value = SetFlags(value, flag)
	}
	return value
}
