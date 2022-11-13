package common

import (
	"github.com/spf13/cobra"
)

type InitCommand struct {
	*cobra.Command
}

func (c *InitCommand) initInitCmd() (err error) {
	_ = &cobra.Command{
		Use:   "init",
		Short: "init chip",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			if len(args) > 0 {
				return cmd.Help()
			}
			return nil
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
	return nil
}
