package pkg

import "github.com/galenliu/chip/transport"

type GroupDataProviderListener struct {
	mTransports transport.Transport
}

func IntGroupDataProviderListener(transport transport.Transport) (*GroupDataProviderListener, error) {
	ins := &GroupDataProviderListener{mTransports: transport}
	return ins, nil
}
