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

var listTaskInput = []struct {
	testCase string
	expected string
}{
	{"no task", "You have no tasks to complete!"},
	{"single task", "1. task1"},
	{"retrieve failed", "Failed to get tasks list"},
}

func TestListCmd(t *testing.T) {
	temp := allTasks
	outputFile := "test_output.txt"
	file, err := os.OpenFile(outputFile, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		t.Error("Failed to Open file:", outputFile)
	}
	oldOutput := os.Stdout
	os.Stdout = file
	for _, item := range listTaskInput {
		removeDb()
		initDbPath()
		if item.testCase == "single task" {
			dbconnect.CreateTask("task1")
		}
		if item.testCase == "retrieve failed" {
			allTasks = func() ([]dbconnect.Task, error) {
				return nil, errors.New("Db all tasks throw error")
			}
		}
		listCmd.Run(listCmd, nil)
		dbconnect.CloseDB()
		file.Seek(0, 0)
		contetnts, _ := ioutil.ReadAll(file)
		output := string(contetnts)
		value := strings.Contains(output, item.expected)
		assert.Equalf(t, true, value, "expected %v got %v", item.expected, output)
		file.Truncate(0)
		file.Seek(0, 0)
		allTasks = temp
	}
	os.Stdout = oldOutput
	file.Close()
	os.Remove(outputFile)
}
