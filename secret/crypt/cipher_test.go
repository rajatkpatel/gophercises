package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	ivValue          = make([]byte, aes.BlockSize)
	cipherBlockErr   = errors.New("TestMode: Cipher block not created")
	ioReadFullErr    = errors.New("TestMode: Failed to fill rand values to iv")
	tempAesNewCipher = aesNewCipherVar
)

var streamInput = []struct {
	input       string
	expectedErr error
}{
	{"key123", nil},
	{"key123", cipherBlockErr},
}

func setAesnewCipherBlock() {
	aesNewCipherVar = func(key []byte) (cipher.Block, error) {
		return nil, cipherBlockErr
	}
}
func TestEncryptStream(t *testing.T) {
	for _, item := range streamInput {
		_, err := encryptStream(item.input, ivValue)
		assert.Equalf(t, item.expectedErr, err, "Expected: %v got %v", item.expectedErr, err)
		defer func() {
			aesNewCipherVar = tempAesNewCipher
		}()
		setAesnewCipherBlock()
	}
}

func TestDecryptStream(t *testing.T) {
	for _, item := range streamInput {
		_, err := decryptStream(item.input, ivValue)
		assert.Equalf(t, item.expectedErr, err, "Expected: %v got %v", item.expectedErr, err)
		defer func() {
			aesNewCipherVar = tempAesNewCipher
		}()
		setAesnewCipherBlock()
	}
}

func TestEncryptWriter(t *testing.T) {
	tempIoReadFull := ioReadFullVar
	file, err := os.OpenFile("test_secret_file", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		t.Error("TestEncryptWriter failed to open file:", err)
	}

	_, err = EncryptWriter("key123", file)
	assert.Equalf(t, nil, err, "Expected: %v got %v", nil, err)

	ioReadFullVar = func(r io.Reader, buf []byte) (n int, err error) {
		return -1, ioReadFullErr
	}
	_, err = EncryptWriter("key123", file)
	assert.Equalf(t, ioReadFullErr, err, "Expected: %v got %v", ioReadFullErr, err)
	ioReadFullVar = tempIoReadFull

	setAesnewCipherBlock()
	_, err = EncryptWriter("key123", file)
	assert.Equalf(t, cipherBlockErr, err, "Expected: %v got %v", cipherBlockErr, err)
	aesNewCipherVar = tempAesNewCipher

	file.Close()
	_, err = EncryptWriter("key123", file)
	assert.Equalf(t, "Unable to write full iv to writer", err.Error(), "Expected: %v got %v", nil, err)
}

func TestDecryptReader(t *testing.T) {
	file, err := os.OpenFile("test_secret_file", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		t.Error("TestEncryptWriter failed to open file:", err)
	}

	_, err = DecryptReader("key123", file)
	assert.Equalf(t, nil, err, "Expected: %v got %v", nil, err)

	_, err = DecryptReader("key123", file)
	assert.NotEqualf(t, nil, err, "Expected: %v got %v", nil, err)

	file.Close()
	file, err = os.OpenFile("test_secret_file", os.O_RDWR|os.O_CREATE, 0755)
	aesNewCipherVar = func(key []byte) (cipher.Block, error) {
		return nil, cipherBlockErr
	}
	_, err = DecryptReader("key123", file)
	assert.Equalf(t, cipherBlockErr, err, "Expected: %v got %v", cipherBlockErr, err)
	aesNewCipherVar = tempAesNewCipher
}
