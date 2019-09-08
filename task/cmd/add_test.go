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

var addTaskInput = []struct {
	input    []string
	expected string
}{
	{[]string{}, "No argument provided with add"},
	{[]string{"task1"}, "Added \"task1\" to your task list."},
	{[]string{"errorTask"}, "\"errorTask\" not added to your task list"},
	{[]string{"Learn Go lang"}, "Added \"Learn Go lang\" to your task list."},
}

func TestAddCmd(t *testing.T) {
	temp := createTask
	outputFile := "test_output.txt"
	file, err := os.OpenFile(outputFile, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		t.Error("Failed to Open file:", outputFile)
	}
	oldOutput := os.Stdout
	os.Stdout = file
	for _, item := range addTaskInput {
		initDbPath()
		if len(item.input) > 0 {
			if item.input[0] == "errorTask" {
				createTask = func(task string) (int, error) {
					return -1, errors.New("Db create task throw error")
				}
			}
		}
		addCmd.Run(addCmd, item.input)
		dbconnect.CloseDB()
		file.Seek(0, 0)
		contetnts, _ := ioutil.ReadAll(file)
		output := string(contetnts)
		value := strings.Contains(output, item.expected)
		assert.Equalf(t, true, value, "expected %v got %v", item.expected, output)
		file.Truncate(0)
		file.Seek(0, 0)
		createTask = temp
	}
	os.Stdout = oldOutput
	file.Close()
	os.Remove(outputFile)
}
