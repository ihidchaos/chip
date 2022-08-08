package commission

import (
	"github.com/galenliu/chip/app"
	"github.com/galenliu/chip/config"
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

func (c *command) intCommission() (err error) {

	cmd := &cobra.Command{
		Use:   "commission",
		Short: "commission mode",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			log.SetReportCaller(true)

			deviceOption := config.NewDeviceOptions()
			deviceOption, _ = deviceOption.Init(c.config)
			err = app.Init(deviceOption)
			if err != nil {
				log.Infof(err.Error())
				return err
			}

			err = app.MainLoop(deviceOption)
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
			config.SetDeviceOptions(cmd)
			_ = c.config.BindPFlags(cmd.Flags())
			return err
		},
	}

	c.root.AddCommand(cmd)
	return nil
}
