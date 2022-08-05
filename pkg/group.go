package pkg

import "github.com/galenliu/chip/messageing/transport/raw"

type GroupDataProviderListener struct {
	mTransports raw.TransportBase
}

func IntGroupDataProviderListener(transport raw.TransportBase) (*GroupDataProviderListener, error) {
	ins := &GroupDataProviderListener{mTransports: transport}
	return ins, nil
}
