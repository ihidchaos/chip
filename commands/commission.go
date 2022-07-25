package commands

import (
	"github.com/galenliu/chip/config"
	"github.com/galenliu/gateway/pkg/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type commissionCommand struct {
	*baseBuilderCommand
	config *viper.Viper
}

func (b *commandsBuilder) newCommssionCommand() *commissionCommand {
	cmd := &cobra.Command{
		Use: "commission",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			config.HandleCHIPConfig(cmd)
			return nil
		},
	}
	cc := &commissionCommand{
		baseBuilderCommand: b.newBuilderCmd(cmd),
		config:             new(viper.Viper),
	}
	cmd.RunE = cc.newDevice
	return cc
}

func (c *commissionCommand) newDevice(cmd *cobra.Command, args []string) error {
	err := c.config.BindPFlags(c.getCommand().Flags())
	options := config.GetDeviceOptions(c.config)
	log.Infof("device options : %v", options)
	return err
}
