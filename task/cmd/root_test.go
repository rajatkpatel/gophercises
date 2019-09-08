package cmd

import (
	"os"
	"path/filepath"

	"github.com/gophercises/task/dbconnect"
	homedir "github.com/mitchellh/go-homedir"
)

func removeDb() {
	dir, _ := homedir.Dir()
	dbPath := filepath.Join(dir, "test.db")
	os.Remove(dbPath)
}

func initDbPath() {
	dir, _ := homedir.Dir()
	dbPath := filepath.Join(dir, "test.db")
	dbconnect.Init(dbPath)
}
