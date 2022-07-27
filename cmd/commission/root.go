package commission

import (
	"errors"
	"github.com/galenliu/chip/config"
	log "github.com/sirupsen/logrus"
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
			Use:           "chip",
			Short:         "chip",
			SilenceErrors: true,
			SilenceUsage:  true,
		},
	}
	for _, o := range opts {
		o(c)
	}

	c.initGlobalFlags()

	err = c.initConfig()
	if err != nil {
		log.Infof("read config file err :%s", err.Error())
	}

	err = c.initInitCmd()
	if err != nil {
		return nil, err
	}

	err = c.intCommission()
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

	globalFlags.StringVarP(&c.configPath, "config", "c", findChipHomeEnv(), "config file (default is $HOME/)")
	globalFlags.StringVarP(&c.configName, "filename", "n", "chip", "config file name (default is chip）)")
	globalFlags.StringVarP(&c.configFileType, "fileType", "t", "yaml", "config file type (default is yaml")
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
	c.config = viper.New()
	if c.configPath != "" {
		c.config.AddConfigPath(c.configPath)
		c.config.SetConfigName(c.configName)
		c.config.SetConfigType(c.configFileType)
	}

	// Environment
	c.config.SetEnvPrefix("chip")
	c.config.AutomaticEnv() // read in environment variables that match
	c.config.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	// If a conf file is found, read it in.
	if err := c.config.ReadInConfig(); err != nil {
		var e viper.ConfigFileNotFoundError
		if !errors.As(err, &e) {
			return err
		}
	}
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
