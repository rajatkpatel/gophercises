package cmd

import (
	"fmt"

	"github.com/gophercises/task/dbconnect"
	"github.com/spf13/cobra"
)

var allTasks = dbconnect.AllTasks

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all of your incomplete tasks",
	Run: func(cmd *cobra.Command, args []string) {
		tasks, err := allTasks()
		if err != nil {
			fmt.Println("Failed to get tasks list: ", err)
			return
		}
		if len(tasks) == 0 {
			fmt.Println("You have no tasks to complete!")
			return
		}
		fmt.Println("You have the following tasks:")
		for i, task := range tasks {
			fmt.Printf("%d. %s\n", i+1, task.Value)
		}
	},
}

func init() {
	RootCmd.AddCommand(listCmd)
}
