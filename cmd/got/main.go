package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "got",
	Short: "Got is a git-like version control system",
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initiliaze an empty Got repository",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Respository initialized")
	},
}

func main() {
	rootCmd.AddCommand(initCmd)
	rootCmd.Execute()
}
