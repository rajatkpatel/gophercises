package dbconnect

import (
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"testing"

	homedir "github.com/mitchellh/go-homedir"

	"github.com/stretchr/testify/assert"
)

var (
	bucketErr = errors.New("Bucket not found")
	dbOpenErr = errors.New("database not open")
)

var DbConnections = []struct {
	input         string
	expectedError bool
}{
	{"test1.db", false},
	{"failCreateBucket", true},
	{"test2.db", false},
	{"/", true},
}

var DbCreateTask = []struct {
	input    string
	expected error
}{
	{"task1", nil},
	{"dbClosed", dbOpenErr},
	{"task2", nil},
	{"bucketUnavailable", bucketErr},
	{"task3", nil},
}

var DbAllTask = []struct {
	expected error
}{
	{nil},
	{dbOpenErr},
	{nil},
	{bucketErr},
	{nil},
}

var DbDeleteTask = []struct {
	input    int
	expected error
}{
	{1, nil},
	{2, dbOpenErr},
	{3, nil},
	{1, bucketErr},
	{2, nil},
}

var dir, _ = homedir.Dir()

func init() {
	log.SetPrefix("LOG: ")
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func setDbPath() {
	dbName := DbConnections[0].input
	dbPath := filepath.Join(dir, dbName)
	Init(dbPath)
}

func TestInit(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered: ", r)
		}
	}()
	for _, item := range DbConnections {
		dbPath := filepath.Join(dir, item.input)
		if item.input == "failCreateBucket" {
			taskBucket = []byte{}
		} else {
			taskBucket = []byte("tasks")
		}
		output := Init(dbPath)
		boolOutput := output != nil
		assert.Equalf(t, item.expectedError, boolOutput, "%s db connection provide error: %v", item.input, boolOutput)
		if !boolOutput {
			db.Close()
		}
	}
}

func TestCloseDB(t *testing.T) {
	setDbPath()
	CloseDB()
}

func TestCreateTask(t *testing.T) {
	for _, item := range DbCreateTask {
		if item.expected != dbOpenErr {
			setDbPath()
		}
		if item.expected == bucketErr {
			taskBucket = []byte{}
		} else {
			taskBucket = []byte("tasks")
		}
		_, err := CreateTask(item.input)
		db.Close()
		assert.Equalf(t, item.expected, err, "expected: %v, got: %v", item.expected, err)
	}

}

func TestAllTasks(t *testing.T) {
	for _, item := range DbAllTask {
		if item.expected != dbOpenErr {
			setDbPath()
		}
		if item.expected == bucketErr {
			taskBucket = []byte{}
		} else {
			taskBucket = []byte("tasks")
		}
		_, err := AllTasks()
		db.Close()
		assert.Equalf(t, item.expected, err, "expected: %v, got: %v", item.expected, err)
	}
}

func TestDeleteTask(t *testing.T) {
	for _, item := range DbDeleteTask {
		if item.expected != dbOpenErr {
			setDbPath()
		}
		if item.expected == bucketErr {
			taskBucket = []byte{}
		} else {
			taskBucket = []byte("tasks")
		}
		err := DeleteTask(item.input)
		db.Close()
		assert.Equalf(t, item.expected, err, "expected: %v, got: %v", item.expected, err)
	}

}
