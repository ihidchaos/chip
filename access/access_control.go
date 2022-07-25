package access

type DeviceTypeResolver struct {
}

type AccessControler interface {
	Init(delegate Delegate, resolver DeviceTypeResolver) error
}

type AccessControl struct {
	mDelegate Delegate
}

func (c *AccessControl) Init(delegate Delegate, d DeviceTypeResolver) error {
	c.mDelegate = delegate
	return nil
}

func SetAccessControl(a AccessControl) {

}
