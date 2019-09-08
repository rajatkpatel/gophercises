package cmd

import (
	"fmt"
	"strconv"

	"github.com/gophercises/task/dbconnect"

	"github.com/spf13/cobra"
)

var (
	allTaskDo  = dbconnect.AllTasks
	deleteTask = dbconnect.DeleteTask
)

// doCmd represents the do command
var doCmd = &cobra.Command{
	Use:   "do",
	Short: "Mark a task on your TODO list as complete",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("No argument provided with do")
			return
		}
		var ids []int
		for _, arg := range args {
			id, err := strconv.Atoi(arg)
			if err != nil {
				fmt.Println("Failed to parse the argument: ", arg)
			} else {
				ids = append(ids, id)
			}
		}

		tasks, err := allTaskDo()
		if err != nil {
			fmt.Println("Failed to get tasks list: ", err)
			return
		}
		for _, id := range ids {
			if id <= 0 || id > len(tasks) {
				fmt.Println("Invalid task number:", id)
				continue
			}
			err = deleteTask(tasks[id-1].Key)
			if err != nil {
				fmt.Printf("Failed to mark \"%d\" as completed. Error: %s\n", id, err)
			} else {
				fmt.Printf("Marked \"%d\" as completed.\n", id)
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(doCmd)
}
