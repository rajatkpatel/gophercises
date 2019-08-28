package cmd

import (
	"errors"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/gophercises/task/dbconnect"
	"github.com/stretchr/testify/assert"
)

var doTaskInput = []struct {
	input    []string
	expected string
}{
	{[]string{}, "No argument provided with do"},
	{[]string{"1"}, "Marked \"1\" as completed."},
	{[]string{"parsingFailed"}, "Failed to parse the argument"},
	{[]string{"100"}, "Invalid task number: 100"},
	{[]string{"999999"}, "Failed to get tasks list"},
	{[]string{"2"}, "Failed to mark \"2\" as completed"},
	{[]string{"1"}, "Marked \"1\" as completed."},
}

func TestDoCmd(t *testing.T) {
	tempAllTask := allTaskDo
	tempDeleteTask := deleteTask
	outputFile := "test_output.txt"
	file, err := os.OpenFile(outputFile, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		t.Error("Failed to Open file:", outputFile)
	}
	oldOutput := os.Stdout
	os.Stdout = file
	for _, item := range doTaskInput {
		initDbPath()
		if len(item.input) > 0 {
			if item.input[0] == "999999" {
				allTaskDo = func() ([]dbconnect.Task, error) {
					return nil, errors.New("Db all tasks throw error")
				}
			} else if item.input[0] == "2" {
				dbconnect.CreateTask("task1")
				dbconnect.CreateTask("task2")
				dbconnect.CreateTask("task3")
				deleteTask = func(key int) error {
					return errors.New("Db delete task throw error")
				}
			}
		}
		doCmd.Run(doCmd, item.input)
		dbconnect.CloseDB()
		file.Seek(0, 0)
		contetnts, _ := ioutil.ReadAll(file)
		output := string(contetnts)
		value := strings.Contains(output, item.expected)
		assert.Equalf(t, true, value, "expected %v got %v", item.expected, output)
		file.Truncate(0)
		file.Seek(0, 0)
		allTaskDo = tempAllTask
		deleteTask = tempDeleteTask
	}
	os.Stdout = oldOutput
	file.Close()
	os.Remove(outputFile)

}
