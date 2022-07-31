package transport

type Transport interface {
	Init(parameters *UdpListenParameters) error
}

type TransportImpl struct {
}

func (t TransportImpl) Init() error {
	return nil
}
