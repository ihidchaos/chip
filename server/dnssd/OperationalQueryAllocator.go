package DnssdServer

type OperationalQueryAllocator struct {
	mAllocator QueryResponderAllocator
}

func (a *OperationalQueryAllocator) GetAllocator() QueryResponderAllocator {
	return a.mAllocator
}
