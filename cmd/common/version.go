package common

import (
	"github.com/spf13/cobra"
)

type VersionCommand struct {
	*cobra.Command
}

func VersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version number",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
}
