package cmd

import (
	"errors"
	"github.com/galenliu/chip/config"
	"github.com/galenliu/gateway/pkg/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"strings"
)

type command struct {
	root           *cobra.Command
	config         *viper.Viper
	deviceOptions  config.DeviceOptions
	configName     string
	configPath     string
	configFileType string
}

type option func(*command)

func newCommand(opts ...option) (c *command, err error) {
	c = &command{
		root: &cobra.Command{
			Use:           "commissionable",
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

	err = c.initConfig()
	if err != nil {
		log.Infof(err.Error())
	}

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
	globalFlags.StringVar(&c.configPath, "config", findChipHomeEnv(), "config file (default is $HOME/)")
	globalFlags.StringVar(&c.configName, "filename", "chip", "config file name (default is .chip)")
	globalFlags.StringVar(&c.configFileType, "config-type", "yaml", "config file type (default is yaml")
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
	if c.configPath != "" {
		conf.AddConfigPath(c.configPath)
		conf.SetConfigName(c.configName)
		conf.SetConfigType(c.configFileType)
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

func findChipHomeEnv() string {
	if p := os.Getenv("CHIP_HOME"); p != "" {
		return p
	}
	dir, err := os.UserHomeDir()
	if err == nil {
		return dir
	}
	return ""
}
