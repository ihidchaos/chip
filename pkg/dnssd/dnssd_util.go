package dnssd

import (
	"fmt"
	"time"
)

type filterType uint8

const (
	KSubtypeServiceNamePart    = "_sub"
	KCommissionableServiceName = "_matterc"
	KCommissionerServiceName   = "_matterd"
	KOperationalServiceName    = "_matter"
	KCommissionProtocol        = "_udp"
	KLocalDomain               = "local"
	KOperationalProtocol       = "_tcp"
)

const kMaxRetryInterval = time.Millisecond * 3600000

type CommissioningModeProvider interface {
	GetCommissioningMode() int
}

// The mode of a Node in which it allows Commissioning.
const (
	CommissioningMode_Disabled        = iota // Commissioning Mode is disabled, CM=0 in DNS-SD key/value pairs
	CommissioningMode_EnableBasic            // Basic Commissioning Mode, CM=1 in DNS-SD key/value pairs
	CommissioningMode_EnabledEnhanced        // Enhanced Commissioning Mode, CM=2 in DNS-SD key/value pairs

	AdvertiseMode_CommissionableNode = iota
	AdvertiseMode_Commissioner

	BroadcastAdvertiseType_RemovingAll = iota
	BroadcastAdvertiseType_Started
)

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
