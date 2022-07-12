package costant

import (
	"flag"
)

var (
	ChipPort                                     uint64
	ChipUdcPort                                  uint64
	ConfigNetworkLayerBle                        bool
	ChipDeviceConfigEnablePairingAutostart       bool
	ChipDeviceConfigEnableCommissionerDiscovery  bool
	ChipDeviceConfigEnableExtendedDiscovery      bool
	ChipDeviceConfigExtendedDiscoveryTimeoutSecs uint64
	ChipDeviceConfigDiscoveryDisabled            uint
	ChipDeviceConfigDiscoveryNoTimeout           int
	InetConfigEnableTcpEndpoint                  bool
)

func Flags() {
	flag.Parse()
	flag.Uint64Var(&ChipPort, "port", 5540, "chip port")
	flag.Uint64Var(&ChipUdcPort, "port", ChipPort+10, "chip port")
	flag.BoolVar(&ConfigNetworkLayerBle, "ble", false, "Config Network Layer Ble")
	flag.BoolVar(&ChipDeviceConfigEnablePairingAutostart, "ble", false, "Chip Device Config EnableP airing Autostart")
	flag.BoolVar(&ChipDeviceConfigEnableCommissionerDiscovery, "commissioner", false, "Chip Device Config Enable Commissioner Discovery")
	flag.BoolVar(&ChipDeviceConfigEnableExtendedDiscovery, "extended", false, "Chip Device Config Enable Commissioner Discovery")
	flag.Uint64Var(&ChipDeviceConfigExtendedDiscoveryTimeoutSecs, "extended-timeout", 15, "Chip DeviceConfig Extended Discovery TimeoutS ecs")
	flag.UintVar(&ChipDeviceConfigDiscoveryDisabled, " device-config-discovery-disabled", 0, "Chip Device Config Discovery Disabled")
	flag.IntVar(&ChipDeviceConfigDiscoveryNoTimeout, " device-config-discovery-no-timeout", 0, "Chip Device Config Discovery No Timeout")
	flag.BoolVar(&InetConfigEnableTcpEndpoint, "Inet Config Enable Tcp Endpoint", true, "Inet Config Enable Tcp Endpoint")
}
