package platform

type CommissionableDataProvider interface {
	GetProductId() (uint16, error)
	GetSetupDiscriminator() (uint16, error)
}

func GetCommissionableDataProvider() CommissionableDataProvider {
	return ConfigurationMgr()
}
