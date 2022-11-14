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

	cmd.Flags().BoolVar(&PairingAutostart,
		PairingAutostartName,
		PairingAutostart,
		PairingAutostartUsage)

	cmd.Flags().IntVar(&CommissionerDiscovery,
		CommissionerDiscoveryName,
		CommissionerDiscovery,
		CommissionerDiscoveryUsage)

	cmd.Flags().BoolVar(&ExtendedDiscovery,
		ExtendedDiscoveryName,
		ExtendedDiscovery,
		ExtendedDiscoveryUsage)

	cmd.Flags().Uint64Var(&ExtendedDiscoveryTimeoutSecs,
		ExtendedDiscoveryTimeoutSecsName,
		ExtendedDiscoveryTimeoutSecs,
		ExtendedDiscoveryTimeoutSecsUsage)

	cmd.Flags().Uint8Var(&DiscoveryDisabled,
		DiscoveryDisabledName,
		DiscoveryDisabled,
		DiscoveryDisabledUsage)

	cmd.Flags().Int8Var(&DiscoveryNoTimeout,
		DiscoveryNoTimeoutName,
		DiscoveryNoTimeout,
		DiscoveryNoTimeoutUsage)

	cmd.Flags().Int8Var(&InetConfigEnableTcpEndpoint,
		InetConfigEnableTcpEndpointName,
		DiscoveryNoTimeout,
		InetConfigEnableTcpEndpointUsage)
}
