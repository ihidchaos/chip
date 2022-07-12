package DeviceLayer

type CommissionableDataProvider interface {
	GetProductId() (uint16, error)
	GetSetupDiscriminator() (uint16, error)
	SetSetupDiscriminator(uint162 uint16)
	GetSpake2pIterationCount() (uint32, error)
	GetSpake2pSalt() ([]byte, error)
	GetSpake2pVerifier() ([]byte, error)
	GetSetupPasscode() (uint32, error)
	SetSetupPasscode(uint322 uint32)
}
