package DeviceLayer

type TestOnlyCommissionableDataProvider struct {
}

func (t TestOnlyCommissionableDataProvider) GetProductId() (uint16, error) {
	//TODO implement me
	panic("implement me")
}

func (t TestOnlyCommissionableDataProvider) GetSetupDiscriminator() (uint16, error) {
	return 0, nil
}

func (t TestOnlyCommissionableDataProvider) SetSetupDiscriminator(uint162 uint16) {
	//TODO implement me
	panic("implement me")
}

func (t TestOnlyCommissionableDataProvider) GetSpake2pIterationCount() (uint32, error) {
	//TODO implement me
	panic("implement me")
}

func (t TestOnlyCommissionableDataProvider) GetSpake2pSalt() ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (t TestOnlyCommissionableDataProvider) GetSpake2pVerifier() ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (t TestOnlyCommissionableDataProvider) GetSetupPasscode() (uint32, error) {
	return 0, nil
}

func (t TestOnlyCommissionableDataProvider) SetSetupPasscode(uint322 uint32) {
	//TODO implement me
	panic("implement me")
}
