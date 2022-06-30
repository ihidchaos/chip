package server

import "github.com/galenliu/chip/transport"

type GroupDataProviderListener struct {
	mTransports transport.TransportManager
}

func IntGroupDataProviderListener(transport transport.TransportManager) (*GroupDataProviderListener, error) {
	ins := &GroupDataProviderListener{mTransports: transport}
	return ins, nil
}
