package cmd

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/gophercises/secret/vault"
	"github.com/stretchr/testify/assert"
)

var getCmdInput = []struct {
	input    []string
	expected string
}{
	{[]string{"test-google-key"}, "test-google-key = test-google-value"},
	{[]string{"test_key_not_present"}, "no value for the key"},
}

func TestGetCmd(t *testing.T) {
	outputFile := "test_output.txt"
	file, err := os.OpenFile(outputFile, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		t.Error("Failed to Open file:", outputFile)
	}
	oldOutput := os.Stdout
	os.Stdout = file
	SecretFile = ".test_secret"
	tempDecryptReader := vault.DecryptReaderVar
	for _, item := range getCmdInput {
		getCmd.Run(getCmd, item.input)
		file.Seek(0, 0)
		contetnts, _ := ioutil.ReadAll(file)
		output := string(contetnts)
		value := strings.Contains(output, item.expected)
		assert.Equalf(t, true, value, "they should be equal::: value %v output %v", item.expected, output)
		file.Truncate(0)
		file.Seek(0, 0)
	}
	os.Stdout = oldOutput
	file.Close()
	os.Remove(outputFile)
	vault.DecryptReaderVar = tempDecryptReader
}
