package config

const (
	RendezvousInformationFlagNone      uint8 = 0
	RendezvousInformationFlagSoftAP    uint8 = 1 << 0
	RendezvousInformationFlagBLE       uint8 = 1 << 1
	RendezvousInformationFlagOnNetwork uint8 = 1 << 2
	CommissioningFlowStandard                = iota
	CommissioningFlowUserActionRequired
	CommissioningFlowCustom
)

const kShortBits uint8 = 4
const kLongBits uint16 = 12

type SetupDiscriminator struct {
	mDiscriminator        uint16
	mIsShortDiscriminator bool
}

func (m *SetupDiscriminator) SetShorValue(discriminator uint8) {
	m.mDiscriminator = uint16(discriminator & kShortBits)
	m.mIsShortDiscriminator = true
}

func (m *SetupDiscriminator) SetLongValue(discriminator uint16) {
	m.mDiscriminator = discriminator & kLongBits
	m.mIsShortDiscriminator = false
}

type PayloadContents struct {
	Version               uint8
	VendorID              uint16
	ProductID             uint16
	CommissioningFlow     uint8
	RendezvousInformation uint8
	Discriminator         SetupDiscriminator
	SetUpPINCode          uint32

	IsValidQRCodePayload bool
	IsValidManualCode    bool
	IsShortDiscriminator bool
}

func (p PayloadContents) CheckPayloadCommonConstraints() bool {
	return true
}
