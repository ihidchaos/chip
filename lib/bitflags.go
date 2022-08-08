package lib

type Uint interface {
	~uint8 | ~uint16 | ~uint32 | ~uint64 | uint
}

func HasFlags[T Uint](value T, flags ...T) (b bool) {
	b = true
	for _, f := range flags {
		if value&f == 0 {
			b = false
		}
	}
	return
}

func SetFlags[T Uint](value T, flags ...T) T {
	for _, f := range flags {
		value = f | value
	}
	return value
}
