package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"log"
)

var (
	aesNewCipherVar = aes.NewCipher
	ioReadFullVar   = io.ReadFull
)

//encryptStream take encoding key and a byte slice.
//It return a cipher stream using the block and iv which is use for encryption.
//If cipher block failed to generate then it will throw error.
func encryptStream(key string, iv []byte) (cipher.Stream, error) {
	block, err := newCipherBlock(key)
	if err != nil {
		log.Println("Failed to generate cipher block:", err)
		return nil, err
	}
	return cipher.NewCFBEncrypter(block, iv), nil
}

//EncryptWriter take two parameters encoding key and a io writer
//encrypt it using random cryptographic number and return cipher streamWriter and error as nil.
//If encryption failed then it will return nil and error.
func EncryptWriter(key string, w io.Writer) (*cipher.StreamWriter, error) {
	iv := make([]byte, aes.BlockSize)
	if _, err := ioReadFullVar(rand.Reader, iv); err != nil {
		log.Println("Failed to fill rand values to iv:", err)
		return nil, err
	}
	stream, err := encryptStream(key, iv)
	if err != nil {
		log.Println("Encrypt stream failed:", err)
		return nil, err
	}
	n, err := w.Write(iv)
	if n != len(iv) || err != nil {
		log.Println("Unable to write full iv to writer: ", err)
		return nil, errors.New("Unable to write full iv to writer")
	}
	return &cipher.StreamWriter{S: stream, W: w}, nil
}

//decryptStream take encoding key and a byte slice.
//It return a cipher stream using the block and iv which is use for decryption.
//If cipher block failed to generate then it will throw error.
func decryptStream(key string, iv []byte) (cipher.Stream, error) {
	block, err := newCipherBlock(key)
	if err != nil {
		log.Println("Failed to generate cipher block:", err)
		return nil, err
	}
	return cipher.NewCFBDecrypter(block, iv), nil
}

//DecryptReader take two parameters encoding key and a io reader
//decrypt it using previously used random cryptographic number in EncryptWriter and return cipher streamReader and error as nil.
//If decryption failed then it will return nil and error.
func DecryptReader(key string, r io.Reader) (*cipher.StreamReader, error) {
	iv := make([]byte, aes.BlockSize)
	n, err := r.Read(iv)
	if n < len(iv) || err != nil {
		log.Println("Unable to read full iv: ", err)
		return nil, errors.New("Unable to read full iv")
	}
	stream, err := decryptStream(key, iv)
	if err != nil {
		log.Println("Decrypt stream failed:", err)
		return nil, err
	}
	return &cipher.StreamReader{S: stream, R: r}, nil
}

//newCipherBlock take encoding key as parameter
//and return a cipher block which is choosen by the md5 hash size.
//If cipher block creation failed due to invalid size i.e. except 16, 24, 32 bytes then it will return error.
func newCipherBlock(key string) (cipher.Block, error) {
	hasher := md5.New()
	fmt.Fprint(hasher, key)
	cipherKey := hasher.Sum(nil)
	return aesNewCipherVar(cipherKey)
}
