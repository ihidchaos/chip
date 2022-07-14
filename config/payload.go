package config

type PayloadContents struct {
	Version               uint8
	VendorID              uint16
	ProductID             uint16
	CommissioningFlow     uint8
	RendezvousInformation uint8
	Discriminator         uint16
	SetUpPINCode          uint32

	IsValidQRCodePayload bool
	IsValidManualCode    bool
	IsShortDiscriminator bool
}

func (p PayloadContents) CheckPayloadCommonConstraints() bool {
	return true
}
