package commands

import (
	"github.com/galenliu/chip/config"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Response struct {
	Err error
	Cmd *cobra.Command
}

func (r Response) IsUserError() bool {
	return true
}

func Execute(args []string) Response {
	commissionCmd := newCommandsBuilder().addAll().build()
	cmd := commissionCmd.getCommand()
	cmd.SetArgs(args)
	c, err := cmd.ExecuteC()

	var resp Response
	resp.Err = err
	resp.Cmd = c

	return resp

}

type chipCommand struct {
	*baseBuilderCommand
}

type chipBuilderCommon struct {
	commonConfig viper.Viper
}

func (c *chipBuilderCommon) handleCommonBuilderFlags(cmd *cobra.Command) {
	config.HandleCHIPConfig(cmd)
}

func (c *chipBuilderCommon) handFlags(cmd *cobra.Command) {
	c.handleCommonBuilderFlags(cmd)
}
