package config

var (
	ChipPort                                    uint16 = 5540
	ChipUdcPort                                 uint16 = ChipPort + 10
	ConfigNetworkLayerBle                       bool
	ChipDeviceConfigEnablePairingAutostart      bool
	ChipDeviceConfigEnableCommissionerDiscovery bool
)

const (
	InetConfigEnableTcpEndpoint = true
)
