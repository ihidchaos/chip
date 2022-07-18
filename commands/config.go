package commands

import (
	"github.com/galenliu/gateway/pkg/log"
	"github.com/spf13/cobra"
)

type configCommand struct {
	*baseBuilderCommand
}

func (b *commandsBuilder) newConfigCommand() *configCommand {
	cc := &configCommand{}
	cmd := &cobra.Command{
		Use:  "config",
		RunE: cc.printConfig,
	}
	cc.baseBuilderCommand = b.newBuilderBasicCmd(cmd)
	return cc

}

func (c configCommand) printConfig(cmd *cobra.Command, args []string) error {
	log.Infof("print config")
	return nil
}
