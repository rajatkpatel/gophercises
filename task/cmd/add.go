package cmd

import (
	"fmt"
	"strings"

	"github.com/gophercises/task/dbconnect"
	"github.com/spf13/cobra"
)

var createTask = dbconnect.CreateTask

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new task to your TODO list",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("No argument provided with add")
			return
		}
		task := strings.Join(args, " ")
		_, err := createTask(task)
		if err != nil {
			fmt.Printf("\"%s\" not added to your task list: %v.\n", task, err)
			return
		}
		fmt.Printf("Added \"%s\" to your task list.\n", task)
	},
}

func init() {
	RootCmd.AddCommand(addCmd)
}
