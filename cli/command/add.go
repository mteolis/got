package command

import (
	"github.com/mteolis/got/internal/repo"
	"github.com/spf13/cobra"
)

var AddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add file contents to the index staged for the next commit",
	Args:  addArgs,
	RunE:  addRunE,
}

func addArgs(cmd *cobra.Command, args []string) error {
	return nil
}

func addRunE(cmd *cobra.Command, args []string) error {
	return repo.AddFile()
}
