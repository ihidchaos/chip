package discovery

type Uint interface {
	~uint8 | ~uint16 | ~uint32 | ~uint64
}

type Filter[T Uint] struct {
	Type         FilterType
	Code         T
	InstanceName string
}
