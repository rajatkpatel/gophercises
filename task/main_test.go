package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestDBConnection struct {
	input         string
	expectedError bool
}

func TestM(t *testing.T) {
	main()
}

func TestInitDB(t *testing.T) {
	dbConnections := []TestDBConnection{
		{"test.db", false},
		{"/", true},
	}
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered: ", r)
		}
	}()
	for _, item := range dbConnections {
		DbName = item.input
		output := initDB()
		boolOutput := output != nil
		assert.Equalf(t, item.expectedError, boolOutput, "%s db connection provide error: %v", DbName, boolOutput)
	}
}
