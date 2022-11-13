package device

import "github.com/spf13/cobra"

type LightCommand struct {
	*cobra.Command
}

func NewLightCommand() *LightCommand {
	return &LightCommand{
		Command: &cobra.Command{
			Use:                    "light",
			Aliases:                nil,
			SuggestFor:             nil,
			Short:                  "light",
			GroupID:                "",
			Long:                   "",
			Example:                "",
			ValidArgs:              nil,
			ValidArgsFunction:      nil,
			Args:                   nil,
			ArgAliases:             nil,
			BashCompletionFunction: "",
			Deprecated:             "",
			Annotations:            nil,
			Version:                "",
			PersistentPreRun:       nil,
			PersistentPreRunE:      nil,
			PreRun:                 nil,
			PreRunE:                nil,
			Run:                    nil,
			RunE: func(cmd *cobra.Command, args []string) error {
				return nil
			},
			PostRun:                    nil,
			PostRunE:                   nil,
			PersistentPostRun:          nil,
			PersistentPostRunE:         nil,
			FParseErrWhitelist:         cobra.FParseErrWhitelist{},
			CompletionOptions:          cobra.CompletionOptions{},
			TraverseChildren:           false,
			Hidden:                     false,
			SilenceErrors:              false,
			SilenceUsage:               false,
			DisableFlagParsing:         false,
			DisableAutoGenTag:          false,
			DisableFlagsInUseLine:      false,
			DisableSuggestions:         false,
			SuggestionsMinimumDistance: 0,
		},
	}
}
