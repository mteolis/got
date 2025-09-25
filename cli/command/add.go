package command

import (
	"github.com/mteolis/got/internal/repo"
	"github.com/spf13/cobra"
)

var Add = &cobra.Command{
	Use:   "add",
	Short: "Add file contents to the index staged for the next commit",
	Run: func(cmd *cobra.Command, args []string) {
		repo.AddFile()
	},
}
