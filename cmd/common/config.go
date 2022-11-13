package common

import (
	"github.com/spf13/cobra"
)

type ConfigCommand struct {
	*cobra.Command
}

func (c *ConfigCommand) initConfiguratorOptionsCmd() (err error) {

	_ = &cobra.Command{
		Use:   "config",
		Short: "Print default or provided configuration in yaml format",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			//if len(args) > 0 {
			//	return cmd.Help()
			//}
			//d := c.config.AllSettings()
			//ym, err := yaml.Marshal(d)
			//if err != nil {
			//	return err
			//}
			//cmd.Println(string(ym))
			return nil

		},
		PreRunE: func(i *cobra.Command, args []string) error {
			return nil
			//config.HandleCHIPConfig(i)
			//err := c.config.BindPFlags(i.Flags())
			//return err
		},
	}

	return nil

}
