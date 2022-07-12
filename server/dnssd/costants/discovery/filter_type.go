package discovery

type FilterType uint8

const (
	None FilterType = iota
	ShortDiscriminator
	LongDiscriminator
	VendorId
	DeviceType
	CommissioningMode
	InstanceName
	Commissioner
	CompressedFabricId
)
