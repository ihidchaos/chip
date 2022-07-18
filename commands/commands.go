package commands

import "github.com/spf13/cobra"

type commandsBuilder struct {
	chipBuilderCommon
	commands []cmder
}

func newCommandsBuilder() *commandsBuilder {
	return &commandsBuilder{}
}

func (b *commandsBuilder) addAll() *commandsBuilder {
	b.commands = append(b.commands, newVersionCommand(), b.newConfigCommand(), b.newCommssionCommand())
	return b
}

func (b *commandsBuilder) build() *chipCommand {
	h := b.newChipCommand()
	for _, cmd := range b.commands {
		h.getCommand().AddCommand(cmd.getCommand())
	}
	return h
}

func (b *commandsBuilder) newChipCommand() *chipCommand {
	cc := &chipCommand{}
	cc.baseBuilderCommand = b.newBuilderCmd(&cobra.Command{
		Use:                        "",
		Aliases:                    nil,
		SuggestFor:                 nil,
		Short:                      "",
		Long:                       "",
		Example:                    "",
		ValidArgs:                  nil,
		ValidArgsFunction:          nil,
		Args:                       nil,
		ArgAliases:                 nil,
		BashCompletionFunction:     "",
		Deprecated:                 "",
		Annotations:                nil,
		Version:                    "",
		PersistentPreRun:           nil,
		PersistentPreRunE:          nil,
		PreRun:                     nil,
		PreRunE:                    nil,
		Run:                        nil,
		RunE:                       nil,
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
	})
	return cc
}

func (b *commandsBuilder) newBuilderCmd(cmd *cobra.Command) *baseBuilderCommand {
	bcmd := &baseBuilderCommand{commandsBuilder: b, baseCommand: &baseCommand{cmd: cmd}}
	bcmd.chipBuilderCommon.handFlags(cmd)
	return bcmd
}

func (b *commandsBuilder) newBuilderBasicCmd(cmd *cobra.Command) *baseBuilderCommand {
	bcmd := &baseBuilderCommand{commandsBuilder: b, baseCommand: &baseCommand{cmd: cmd}}
	bcmd.chipBuilderCommon.handleCommonBuilderFlags(cmd)
	return bcmd
}

type baseCommand struct {
	cmd *cobra.Command
}

func (b baseCommand) getCommand() *cobra.Command {
	return b.cmd
}

type baseBuilderCommand struct {
	*baseCommand
	*commandsBuilder
}
