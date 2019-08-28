package vault

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
	"sync"

	"github.com/gophercises/secret/crypt"
)

type Vault struct {
	encodingKey string
	keyValues   map[string]string
	filePath    string
	mutex       sync.Mutex
}

//File take encoding key and file path
//Set both the values to Vault type
func File(encodingKey, filePath string) *Vault {
	return &Vault{
		encodingKey: encodingKey,
		filePath:    filePath,
	}
}

var (
	DecryptReaderVar = crypt.DecryptReader
	EncryptWriterVar = crypt.EncryptWriter
)

//load method make a map for keyValues if it is not present earlier.
//It read the decoded json value from the file.
//If unable to get a cipher streamReader it will return error.
func (v *Vault) load() error {
	file, err := os.Open(v.filePath)
	if err != nil {
		log.Println("Failed to open file:", err)
		v.keyValues = make(map[string]string)
		return nil
	}
	defer file.Close()
	reader, err := DecryptReaderVar(v.encodingKey, file)
	if err != nil {
		return err
	}
	return v.readKeyValues(reader)
}

//save method write the encoded json to the file.
//If unable to get a cipher streamWriter it will return error.
func (v *Vault) save() error {
	file, err := os.OpenFile(v.filePath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Println("Failed to open file:", err)
		return err
	}
	defer file.Close()
	writer, err := EncryptWriterVar(v.encodingKey, file)
	if err != nil {
		return err
	}
	return v.writeKeyValues(writer)
}

//readKeyValues method read encoded json data from the provided reader parameter.
//It will store the decoded json to the vault's keyValues.
//If json decoding failed it will return error.
func (v *Vault) readKeyValues(r io.Reader) error {
	decoder := json.NewDecoder(r)
	return decoder.Decode(&v.keyValues)
}

//writeKeyValues method encode the vault's keyValues and
// write the encoded json to the provided writer parameter.
//If json encoding failed it will return error.
func (v *Vault) writeKeyValues(w io.Writer) error {
	encoder := json.NewEncoder(w)
	return encoder.Encode(v.keyValues)
}

//Get take a parameter key.
//If the key is present in the secret file
//it will return its value after decrypt it
//otherwise return error.
func (v *Vault) Get(key string) (string, error) {
	v.mutex.Lock()
	defer v.mutex.Unlock()

	err := v.load()
	if err != nil {
		log.Println("Failed to load key-values: ", err)
		return "", err
	}
	value, ok := v.keyValues[key]
	if !ok {
		log.Printf("Failed to get value of key %v : %v", key, err)
		return "", errors.New("Vault has no value for the key provided")
	}
	return value, nil
}

//Set take two parameter key and value.
//It will save the key and value to the secret file
//after encrypt it.
func (v *Vault) Set(key, value string) error {
	v.mutex.Lock()
	v.mutex.Unlock()

	err := v.load()
	if err != nil {
		log.Println("Failed to load key-values: ", err)
		return err
	}
	v.keyValues[key] = value
	err = v.save()
	if err != nil {
		log.Println("Failed to save key-values: ", err)
	}
	return err
}
