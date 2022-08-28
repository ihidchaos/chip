package credentials

type DeviceAttestationCredentialsProvider interface {
	GetCertificationDeclaration()
	GetFirmwareInformation()
	GetDeviceAttestationCert()
	GetProductAttestationIntermediateCert()
	SignWithDeviceAttestationKey()
	IsDeviceAttestationCredentialsProviderSet() bool
}

type UnimplementedDACProvider struct {
}

var gDacProvider DeviceAttestationCredentialsProvider

func SetDeviceAttestationCredentialsProvider(provider DeviceAttestationCredentialsProvider) {
	gDacProvider = provider
}

func GetDeviceAttestationCredentialsProvider() DeviceAttestationCredentialsProvider {
	return gDacProvider
}
