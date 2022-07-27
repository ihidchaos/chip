package transport

type TransportManager interface {
	Init() error
}

type TransportImpl struct {
}

func NewTransportImpl() *TransportImpl {
	return &TransportImpl{}
}

func (t TransportImpl) Init() error {
	return nil
}
