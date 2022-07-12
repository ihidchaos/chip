package cmd

import "github.com/spf13/cobra"

func (c *command) initVendorIdCmd() {
	v := &cobra.Command{
		Use:   "ve",
		Short: "Print device VendorID",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println(c.deviceOptions.Payload.VendorID)
		},
	}
	v.SetOut(c.root.OutOrStdout())
	c.root.AddCommand(v)
}
