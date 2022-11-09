package core

type CASEClientPoolDelegate interface {
}

type CASEClientPool struct {
	mClientPool []*CASEClient
}

func NewCASEClientPool(size int) *CASEClientPool {
	return &CASEClientPool{
		mClientPool: make([]*CASEClient, size),
	}
}
