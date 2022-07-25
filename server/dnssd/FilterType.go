package dnssd

import "fmt"

type filterType uint8

type Uint interface {
	~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uint | ~int
}

const (
	Subtype_None filterType = iota
	Subtype_ShortDiscriminator
	Subtype_LongDiscriminator
	Subtype_VendorId
	Subtype_DeviceType
	Subtype_CommissioningMode
	Subtype_InstanceName
	Subtype_Commissioner
	Subtype_CompressedFabricId
)

func makeServiceSubtype[T Uint](filter filterType, values ...T) string {
	var val T = 0
	if len(values) != 0 {
		val = values[0]
	}
	switch filter {
	case Subtype_ShortDiscriminator:
		return fmt.Sprintf("_S%d", val)
	case Subtype_LongDiscriminator:
		return fmt.Sprintf("_L%d", val)
	case Subtype_VendorId:
		return fmt.Sprintf("_V%d", val)
	case Subtype_DeviceType:
		return fmt.Sprintf("_T%d", val)
	case Subtype_CommissioningMode:
		return "_CM"
	case Subtype_Commissioner:
		return fmt.Sprintf("_D%d", val)
	case Subtype_CompressedFabricId:
		return fmt.Sprintf("_I%016X", val)
	case Subtype_InstanceName:
		return fmt.Sprintf("%016X", val)
	case Subtype_None:
		return ""
	default:
		return fmt.Sprintf("%s", val)
	}
}
