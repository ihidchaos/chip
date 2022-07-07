package platform

type DeviceInstanceInfoProvider interface {
	GetVendorId() (uint16, error)
	GetProductId() (uint16, error)
}

func GetDeviceInstanceInfoProvider() DeviceInstanceInfoProvider {
	return ConfigurationMgr()
}
