package cmd

import (
	"github.com/galenliu/chip/app"
	"github.com/galenliu/chip/config"
	"github.com/galenliu/gateway/pkg/log"
	"github.com/spf13/cobra"
)

func (c *command) newCommission() (err error) {

	cmd := &cobra.Command{
		Use:   "commission",
		Short: "commission mode",
		RunE: func(cmd *cobra.Command, args []string) (err error) {

			err = app.AppMainInit()
			if err != nil {
				log.Infof(err.Error())
				return err
			}

			err = app.AppMainLoop()
			if err != nil {
				log.Infof(err.Error())
				return err
			}

			if len(args) > 0 {
				return cmd.Help()
			}

			return nil
		},

		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			config.FlagsDeviceOptions(cmd)
			config.HandleCHIPConfig(cmd)
			_ = c.config.BindPFlags(cmd.Flags())
			return err
		},
	}

	c.root.AddCommand(cmd)
	return nil
}
