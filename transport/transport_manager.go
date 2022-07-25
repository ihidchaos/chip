package transport

type TransportManager interface {
	Init() error
}

type TransportImpl struct {
}
