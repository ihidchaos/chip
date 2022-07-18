package commands

import (
	"github.com/galenliu/gateway/pkg/log"
	"github.com/spf13/cobra"
)

type versionCommand struct {
	*baseCommand
}

func newVersionCommand() *versionCommand {
	return &versionCommand{&baseCommand{&cobra.Command{
		Use: "version",
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Infof("chip version ")
			return nil
		},
	}}}
}
