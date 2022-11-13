package config

import (
	"github.com/spf13/cobra"
)

var (
	NetworkLayerBle      = false
	NetworkLayerBleName  = "ble"
	NetworkLayerBleUsage = "Chip Config Network Layer Ble"

	RendezvousMode = false

	InetConfigEnableTcpEndpoint      int8 = 0
	InetConfigEnableTcpEndpointName       = "enable-tcp"
	InetConfigEnableTcpEndpointUsage      = "Inet Config Enable Tcp Endpoint 0/1(able/disable)"
)

func MatterFlags(cmd *cobra.Command) {
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
