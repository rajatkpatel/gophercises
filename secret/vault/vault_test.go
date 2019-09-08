package vault

import (
	"crypto/cipher"
	"errors"
	"io"
	"os"
	"path/filepath"
	"testing"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/stretchr/testify/assert"
)

var (
	encodingKey       = "test_key123"
	tempDecryptReader = DecryptReaderVar
	tempEncryptWriter = EncryptWriterVar
)

func setSecretsPath(secretFile string) string {
	dir, _ := homedir.Dir()
	return filepath.Join(dir, secretFile)
}

func setDecryptReader() {
	DecryptReaderVar = func(key string, r io.Reader) (*cipher.StreamReader, error) {
		return nil, errors.New("Decrypt reader failed")
	}
}

func setEncryptWriter() {
	EncryptWriterVar = func(key string, w io.Writer) (*cipher.StreamWriter, error) {
		return nil, errors.New("Encrypt Writer failed")
	}
}

func TestFile(t *testing.T) {
	filePath := setSecretsPath("test_secrets")
	v := File(encodingKey, filePath)
	assert.Equalf(t, encodingKey, v.encodingKey, "expected %v got %v", encodingKey, v.encodingKey)
}

func TestSave(t *testing.T) {
	filePath := setSecretsPath("test_secrets")
	v := File(encodingKey, filePath)
	err := v.save()
	assert.Equalf(t, false, err != nil, "expected %v got %v", false, err != nil)

	setEncryptWriter()
	err = v.save()
	assert.Equalf(t, true, err != nil, "expected %v got %v", true, err != nil)
	EncryptWriterVar = tempEncryptWriter

	val := File(encodingKey, "")
	err = val.save()
	assert.Equalf(t, true, err != nil, "expected %v got %v", true, err != nil)
}

func TestLoad(t *testing.T) {
	filePath := setSecretsPath("test_secrets")
	v := File(encodingKey, filePath)

	err := v.load()
	assert.Equalf(t, nil, err, "expected %v got %v", nil, err)

	setDecryptReader()
	err = v.load()
	assert.Equalf(t, "Decrypt reader failed", err.Error(), "expected %v got %v", "Decrypt reader failed", err.Error())
	DecryptReaderVar = tempDecryptReader

	val := File(encodingKey, "")
	err = val.load()
	assert.Equalf(t, nil, err, "expected %v got %v", nil, err)
}

func TestSet(t *testing.T) {
	filePath := setSecretsPath("test_secrets")
	v := File(encodingKey, filePath)
	os.Remove(filePath)

	err := v.Set("test_google_key", "test_asdf")
	assert.Equalf(t, nil, err, "expected %v got %v", nil, err)

	setDecryptReader()
	err = v.Set("test_google_key", "test_asdf")
	assert.Equalf(t, "Decrypt reader failed", err.Error(), "expected %v got %v", "Decrypt reader failed", err.Error())
	DecryptReaderVar = tempDecryptReader

	setEncryptWriter()
	err = v.Set("test_google_key", "test_asdf")
	assert.Equalf(t, "Encrypt Writer failed", err.Error(), "expected %v got %v", "Encrypt Writer failed", err.Error())
	EncryptWriterVar = tempEncryptWriter
}

func TestGet(t *testing.T) {
	filePath := setSecretsPath("test_secrets")
	v := File(encodingKey, filePath)

	value, _ := v.Get("test_google_key")
	assert.Equalf(t, "test_asdf", value, "expected %v got %v", "test_asdf", value)

	setDecryptReader()
	value, err := v.Get("test_google_key")
	assert.Equalf(t, "Decrypt reader failed", err.Error(), "expected %v got %v", "Decrypt reader failed", err.Error())
	DecryptReaderVar = tempDecryptReader

	value, err = v.Get("test_key_not_Present")
	assert.Equalf(t, "Vault has no value for the key provided", err.Error(), "expected %v got %v", "Vault has no value for the key provided", err.Error())

}
