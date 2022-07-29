package dac

type ExampleDACProvider interface {
	DeviceAttestationCredentialsProvider
}

type ExampleDACProviderImpl struct {
}

func (e ExampleDACProviderImpl) GetCertificationDeclaration() {
	//TODO implement me
	panic("implement me")
}

func (e ExampleDACProviderImpl) GetFirmwareInformation() {
	//TODO implement me
	panic("implement me")
}

func (e ExampleDACProviderImpl) GetDeviceAttestationCert() {
	//TODO implement me
	panic("implement me")
}

func (e ExampleDACProviderImpl) GetProductAttestationIntermediateCert() {
	//TODO implement me
	panic("implement me")
}

func (e ExampleDACProviderImpl) SignWithDeviceAttestationKey() {
	//TODO implement me
	panic("implement me")
}

func (e ExampleDACProviderImpl) IsDeviceAttestationCredentialsProviderSet() bool {
	//TODO implement me
	panic("implement me")
}
