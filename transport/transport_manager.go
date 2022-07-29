package transport

type Transport interface {
}

type TransportImpl struct {
}

func (t TransportImpl) Init() error {
	return nil
}
