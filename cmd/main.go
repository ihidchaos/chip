package main

import (
	"fmt"
	"github.com/galenliu/chip/cmd/common"
	"github.com/galenliu/chip/cmd/device"
	"github.com/galenliu/chip/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

func main() {
	rootCmd := initRootCommand()
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

type RootCommand struct {
	*cobra.Command
	config         *viper.Viper
	deviceOptions  config.DeviceOptions
	configName     string
	configPath     string
	configFileType string
}

func initRootCommand() *RootCommand {
	rootCommand := &RootCommand{
		Command: &cobra.Command{
			Use:           "matter",
			Short:         "matter",
			SilenceErrors: true,
			SilenceUsage:  true,
		},
	}

	rootCommand.AddCommand(common.VersionCmd())
	rootCommand.AddCommand()

	lightCmd := device.NewLightCommand()
	rootCommand.AddCommand(lightCmd.Command)
	return rootCommand
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
