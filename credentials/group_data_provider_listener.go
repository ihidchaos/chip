package credentials

type GroupDataProviderListener interface {
	Init(s ServerDelegate) error
}

type GroupDataProviderListenerImpl struct {
}

func (g GroupDataProviderListenerImpl) Init(s ServerDelegate) error {
	return nil
}

func NewGroupDataProviderListenerImpl() *GroupDataProviderListenerImpl {
	return &GroupDataProviderListenerImpl{}
}
