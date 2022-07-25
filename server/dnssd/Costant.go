package dnssd

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
