package access

type DeviceTypeResolver struct {
}

type Controler interface {
	Init(delegate Delegate, resolver DeviceTypeResolver) error
}

type Control struct {
	mDelegate Delegate
}

func NewAccessControl() *Control {
	return &Control{}
}

func (c *Control) Init(delegate Delegate, d DeviceTypeResolver) error {
	c.mDelegate = delegate
	return nil
}

func SetAccessControl(a Controler) {

}
