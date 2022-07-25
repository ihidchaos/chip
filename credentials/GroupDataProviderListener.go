package credentials

type GroupDataProviderListener interface {
	Init(s ServerDelegate) error
}

type GroupDataProviderListenerImpl struct {
}
