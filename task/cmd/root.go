package cmd

import "github.com/spf13/cobra"

//RootCmd is the initial task command.
var RootCmd = &cobra.Command{
	Use:   "task",
	Short: "Task is a CLI for managing your TODOs.",
}
