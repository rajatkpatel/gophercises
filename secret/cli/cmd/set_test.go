package cmd

import (
	"crypto/cipher"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/gophercises/secret/vault"
	"github.com/stretchr/testify/assert"
)

var setCmdInput = []struct {
	input    []string
	expected string
}{
	{[]string{"test-google-key", "test-google-value"}, "value set successfully"},
	{[]string{"test_decrypt_err", "decrypt_msg"}, "Decrypt reader failed"},
}

func TestSetCmd(t *testing.T) {
	outputFile := "test_output.txt"
	file, err := os.OpenFile(outputFile, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		t.Error("Failed to Open file:", outputFile)
	}
	oldOutput := os.Stdout
	os.Stdout = file
	SecretFile = ".test_secret"
	tempDecryptReader := vault.DecryptReaderVar
	for _, item := range setCmdInput {
		setCmd.Run(setCmd, item.input)
		file.Seek(0, 0)
		contetnts, _ := ioutil.ReadAll(file)
		output := string(contetnts)
		value := strings.Contains(output, item.expected)
		assert.Equalf(t, true, value, "they should be equal")
		file.Truncate(0)
		file.Seek(0, 0)
		vault.DecryptReaderVar = func(key string, r io.Reader) (*cipher.StreamReader, error) {
			return nil, errors.New("Decrypt reader failed")
		}
	}
	os.Stdout = oldOutput
	file.Close()
	os.Remove(outputFile)
	vault.DecryptReaderVar = tempDecryptReader
}
