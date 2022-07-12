package CommissioningMode

const (
	Disabled        uint8 = iota // Commissioning Mode is disabled, CM=0 in DNS-SD key/value pairs
	EnableBasic                  // Basic Commissioning Mode, CM=1 in DNS-SD key/value pairs
	EnabledEnhanced              // Enhanced Commissioning Mode, CM=2 in DNS-SD key/value pairs
)
