package cmd

import (
	"github.com/galenliu/chip/app"
	"github.com/galenliu/chip/config"
	DeviceLayer "github.com/galenliu/chip/platform/device_layer"
	"github.com/galenliu/gateway/pkg/log"
	"github.com/spf13/cobra"
)

func (c *command) initTestCmd() (err error) {

	cmd := &cobra.Command{
		Use:   "test",
		Short: "chip test",
		RunE: func(cmd *cobra.Command, args []string) (err error) {

			deviceOptions := config.GetDeviceOptions(c.config)
			info := &DeviceLayer.DeviceInfo{
				VendorId:        0xFFF0,
				ProductId:       0x123,
				ProductName:     "34234",
				VendorName:      "234234",
				HardwareVersion: "234234",
			}
			err = app.AppMainInit(info, deviceOptions)
			if err != nil {
				log.Infof(err.Error())
				return err
			}

			err = app.AppMainLoop(deviceOptions, info)
			if err != nil {
				log.Infof(err.Error())
				return err
			}

			if len(args) > 0 {
				return cmd.Help()
			}

			log.Infof("%#v", deviceOptions)

			//if len(args) > 0 {
			//	return cmd.Help()
			//}
			//
			//log.Infof("gateway version %v", constant.Version)
			//
			//ctx, cancelFunc := context.WithCancel(context.Background())
			//defer cancelFunc()
			//deviceOption := options.DeviceOptions{
			//	PICS: "",
			//	KVS:  "",
			//}
			//
			//if err != nil {
			//	cancelFunc()
			//	return err
			//}
			//
			//// Wait for termination or interrupt signals.
			//// We want to clean up things at the end.
			//interruptChannel := make(chan os.Signal, 1)
			//signal.Notify(interruptChannel, syscall.SIGINT, syscall.SIGTERM)

			return nil
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			initDeviceOptionsFlags(cmd)
			config.SetCHIPConfig(cmd)
			err := c.config.BindPFlags(cmd.Flags())
			return err
		},
	}
	c.root.AddCommand(cmd)
	return nil
}
