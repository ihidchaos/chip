package credentials

type ServerDelegate interface {
}

type ServerFabricDelegate interface {
	Init(s ServerDelegate) error
}

type ServerFabricDelegateImpl struct {
}
