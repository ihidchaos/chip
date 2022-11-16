package config

import "github.com/spf13/cobra"

var (
	PairingAutostart      = false
	PairingAutostartName  = "pairing-autostart"
	PairingAutostartUsage = "Matter Device Config Enable Pairing Autostart"

	CommissionerDiscovery      = 1
	CommissionerDiscoveryName  = "enable-commissioner-discover"
	CommissionerDiscoveryUsage = "Matter Device Config Enable Commissioner Discovery"

	ExtendedDiscovery      = false
	ExtendedDiscoveryName  = "enable-extended-discovery"
	ExtendedDiscoveryUsage = "ExtendedDiscovery"

	ExtendedDiscoveryTimeoutSecs      uint64 = 15 * 60
	ExtendedDiscoveryTimeoutSecsName         = "extended-discovery-timeout"
	ExtendedDiscoveryTimeoutSecsUsage        = "Matter Device Config Extended Discovery TimeoutSecs"

	DiscoveryDisabled      uint8 = 0
	DiscoveryDisabledName        = "discovery-disabled"
	DiscoveryDisabledUsage       = "Matter Device Config Discovery Disabled ProtocolName"

	DiscoveryNoTimeout      int8 = -1
	DiscoveryNoTimeoutName       = "discover-no-timeout"
	DiscoveryNoTimeoutUsage      = "Matter Device Config Discovery NoTimeout"
)

func DeviceFlags(cmd *cobra.Command) {
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
