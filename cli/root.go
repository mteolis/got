package cli

import (
	"github.com/mteolis/got/cli/command"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "got",
	Short: "Got is a git-like version control system",
}

func Execute() {
	rootCmd.AddCommand(command.Init)
	rootCmd.AddCommand(command.Add)
	rootCmd.Execute()
}
