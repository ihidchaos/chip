package config

import (
	"github.com/spf13/cobra"
)

var (
	NetworkLayerBle      = false
	NetworkLayerBleName  = "ble"
	NetworkLayerBleUsage = "Chip Config Network Layer Ble"

	RendezvousMode = false

	ChipDeviceConfigEnablePairingAutostart      = false
	ChipDeviceConfigEnablePairingAutostartName  = "pairing-autostart"
	ChipDeviceConfigEnablePairingAutostartUsage = "Chip Device Config Enable Pairing Autostart"

	ChipDeviceConfigEnableCommissionerDiscovery      = 1
	ChipDeviceConfigEnableCommissionerDiscoveryName  = "enable-commissioner-discover"
	ChipDeviceConfigEnableCommissionerDiscoveryUsage = "Chip Device Config Enable Commissioner Discovery"

	ChipDeviceConfigEnableExtendedDiscovery      = false
	ChipDeviceConfigEnableExtendedDiscoveryName  = "enable-extended-discovery"
	ChipDeviceConfigEnableExtendedDiscoveryUsage = "ChipDeviceConfigEnableExtendedDiscovery"

	ChipDeviceConfigExtendedDiscoveryTimeoutSecs      uint64 = 15 * 60
	ChipDeviceConfigExtendedDiscoveryTimeoutSecsName         = "extended-discovery-timeout"
	ChipDeviceConfigExtendedDiscoveryTimeoutSecsUsage        = "Chip Device Config Extended Discovery TimeoutSecs"

	ChipDeviceConfigDiscoveryDisabled      uint8 = 0
	ChipDeviceConfigDiscoveryDisabledName        = "discovery-disabled"
	ChipDeviceConfigDiscoveryDisabledUsage       = "Chip Device Config Discovery Disabled Name"

	ChipDeviceConfigDiscoveryNoTimeout      int8 = -1
	ChipDeviceConfigDiscoveryNoTimeoutName       = "discover-no-timeout"
	ChipDeviceConfigDiscoveryNoTimeoutUsage      = "Chip Device Config Discovery NoT imeout"

	InetConfigEnableTcpEndpoint      int8 = 0
	InetConfigEnableTcpEndpointName       = "enable-tcp"
	InetConfigEnableTcpEndpointUsage      = "Inet Config Enable Tcp Endpoint 0/1(able/disable)"
)

func HandleCHIPConfig(cmd *cobra.Command) {
	cmd.Flags().BoolVar(&NetworkLayerBle,
		NetworkLayerBleName,
		NetworkLayerBle,
		NetworkLayerBleUsage)

	cmd.Flags().BoolVar(&ChipDeviceConfigEnablePairingAutostart,
		ChipDeviceConfigEnablePairingAutostartName,
		ChipDeviceConfigEnablePairingAutostart,
		ChipDeviceConfigEnablePairingAutostartUsage)

	cmd.Flags().IntVar(&ChipDeviceConfigEnableCommissionerDiscovery,
		ChipDeviceConfigEnableCommissionerDiscoveryName,
		ChipDeviceConfigEnableCommissionerDiscovery,
		ChipDeviceConfigEnableCommissionerDiscoveryUsage)

	cmd.Flags().BoolVar(&ChipDeviceConfigEnableExtendedDiscovery,
		ChipDeviceConfigEnableExtendedDiscoveryName,
		ChipDeviceConfigEnableExtendedDiscovery,
		ChipDeviceConfigEnableExtendedDiscoveryUsage)

	cmd.Flags().Uint64Var(&ChipDeviceConfigExtendedDiscoveryTimeoutSecs,
		ChipDeviceConfigExtendedDiscoveryTimeoutSecsName,
		ChipDeviceConfigExtendedDiscoveryTimeoutSecs,
		ChipDeviceConfigExtendedDiscoveryTimeoutSecsUsage)

	cmd.Flags().Uint8Var(&ChipDeviceConfigDiscoveryDisabled,
		ChipDeviceConfigDiscoveryDisabledName,
		ChipDeviceConfigDiscoveryDisabled,
		ChipDeviceConfigDiscoveryDisabledUsage)

	cmd.Flags().Int8Var(&ChipDeviceConfigDiscoveryNoTimeout,
		ChipDeviceConfigDiscoveryNoTimeoutName,
		ChipDeviceConfigDiscoveryNoTimeout,
		ChipDeviceConfigDiscoveryNoTimeoutUsage)

	cmd.Flags().Int8Var(&InetConfigEnableTcpEndpoint,
		InetConfigEnableTcpEndpointName,
		ChipDeviceConfigDiscoveryNoTimeout,
		InetConfigEnableTcpEndpointUsage)
}
