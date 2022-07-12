package cmd

import (
	"errors"
	"github.com/galenliu/chip/config"
	"github.com/galenliu/chip/platform/options"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strings"
)

type command struct {
	root                  *cobra.Command
	config                *viper.Viper
	deviceOptions         options.DeviceOptions
	deviceOptionsFile     string
	deviceOptionsFilePath string
}

type option func(*command)

func newCommand(opts ...option) (c *command, err error) {
	c = &command{
		root: &cobra.Command{
			Use:           "commissionable node",
			Short:         "commission",
			SilenceErrors: true,
			SilenceUsage:  true,
			PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
				return c.initConfig()
			},
		},
	}

	for _, o := range opts {
		o(c)
	}

	c.initGlobalFlags()

	err = c.initInitCmd()
	if err != nil {
		return nil, err
	}

	err = c.initTestCmd()
	if err != nil {
		return nil, err
	}

	if err := c.initConfiguratorOptionsCmd(); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *command) initGlobalFlags() {
	globalFlags := c.root.PersistentFlags()
	globalFlags.StringVar(&c.deviceOptionsFile, "config", "", "config file (default is $HOME/.chip.yaml)")
}

func (c *command) Execute() (err error) {
	return c.root.Execute()
}

// Execute parses command line arguments and runs appropriate functions.
func Execute() (err error) {
	c, err := newCommand()
	if err != nil {
		return err
	}
	return c.Execute()
}

func (c *command) initConfig() (err error) {
	conf := viper.New()
	configName := ".chip"
	if c.deviceOptionsFile != "" {
		// Use conf file from the flag.
		conf.SetConfigFile(c.deviceOptionsFile)
	} else {
		// Search conf in home directory with name ".gateway" (without extension).
		conf.AddConfigPath(c.deviceOptionsFilePath)
		conf.SetConfigName(configName)
		conf.SetConfigType("yaml")
	}

	// Environment
	conf.SetEnvPrefix("chip")
	conf.AutomaticEnv() // read in environment variables that match
	conf.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	// If a conf file is found, read it in.
	if err := conf.ReadInConfig(); err != nil {
		var e viper.ConfigFileNotFoundError
		if !errors.As(err, &e) {
			return err
		}
	}
	c.config = conf
	return nil
}

func (c *command) initDeviceOptionsFlags(cmd *cobra.Command) {
	cmd.Flags().Uint8(config.KDeviceOption_Version, 0,
		"       The version indication provides versioning of the setup payload.\n")
	cmd.Flags().Uint64(config.KDeviceOption_VendorID, 0,
		"       The Vendor ID is assigned by the Connectivity Standards Alliance.\n")
	cmd.Flags().Uint64(config.KDeviceOption_ProductID, 0xFFFF,
		"       The Product ID is specified by vendor.\n")
	cmd.Flags().Uint8(config.KDeviceOption_CustomFlow, 0,
		"       A 2-bit unsigned enumeration specifying manufacturer-specific custom flow options.\n")
	cmd.Flags().Uint8(config.KDeviceOption_Capabilities, 0,
		"       Discovery Capabilities Bitmask which contains information about Device’s available technologies for device discovery.\n")
	cmd.Flags().Uint16(config.KDeviceOption_Discriminator, 0xFF12,
		"       A 12-bit unsigned integer match the value which a device advertises during commissioning.\n")
	cmd.Flags().Uint32(config.KDeviceOption_Passcode, 0xFFFFFFF,
		"       If not provided to compute a verifier, the --spake2p-verifier-base64 must be provided. \n")
	cmd.Flags().Uint32(config.KDeviceOption_Spake2pVerifierBase64, 0xFFFFF,
		"       A raw concatenation of 'W0' and 'L' (67 bytes) as base64 to override the verifier\n "+
			"auto-computed from the passcode, if provided.\n")
	cmd.Flags().Uint32(config.KDeviceOption_Spake2pSaltBase64, 0x12121,
		"       16-32 bytes of salt to use for the PASE verifier, as base64. If omitted, will be generated\n "+
			"randomly. If a --spake2p-verifier-base64 is passed, it must match against the salt otherwise\n "+
			"failure will arise.\n ")
	cmd.Flags().Uint64(config.KDeviceOption_Spake2pIterations, 0xffffff, "       Number of PBKDF iterations to use. If omitted, will be 1000. If a --spake2p-verifier-base64 is\n"+
		"passed, the iteration counts must match that used to generate the verifier otherwise failure will\n "+
		"arise.\n")
	cmd.Flags().Uint16(config.KDeviceOption_SecuredDevicePort, 5540,
		"       A 16-bit unsigned integer specifying the listen port to use for secure device messages (default is 5540).\n")
	cmd.Flags().Uint16(config.KDeviceOption_SecuredCommissionerPort, 5542,
		"       A 16-bit unsigned integer specifying the listen port to use for secure commissioner messages (default is 5542). Only ")
	cmd.Flags().Uint16(config.KDeviceOption_UnsecuredCommissionerPort, 5550,
		"       A 16-bit unsigned integer specifying the port to use for unsecured commissioner messages (default is 5550).\n")
	cmd.Flags().String(config.KDeviceOption_Command, "", "       A name for a command to execute during startup.\n")
	cmd.Flags().String(config.KDeviceOption_PICS, "",
		"       A file containing PICS items.\n")
	cmd.Flags().String(config.KDeviceOption_KVS, "/Users/liuguilin/chip.ini",
		"       A file containing PICS items.\n")
	cmd.Flags().String(config.KDeviceOption_InterfaceId, "", "       A interface id to advertise on.\n")
}
