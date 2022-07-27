package credentials

type ServerDelegate interface {
}

type ServerFabricDelegate interface {
	Init(s ServerDelegate) error
}

type ServerFabricDelegateImpl struct {
}

func (s2 ServerFabricDelegateImpl) Init(s ServerDelegate) error {
	return nil
}

func NewServerFabricDelegateImpl() *ServerFabricDelegateImpl {
	return &ServerFabricDelegateImpl{}
}
