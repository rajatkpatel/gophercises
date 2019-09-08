package main

import (
	"log"
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"

	"github.com/gophercises/task/cmd"
	"github.com/gophercises/task/dbconnect"
)

//DbName variable store the database name of task command prompt.
var DbName = "tasksCmd.db"

func main() {
	initDB()
}

func initDB() error {
	dir, _ := homedir.Dir()
	dbPath := filepath.Join(dir, DbName)
	err := dbconnect.Init(dbPath)
	if err != nil {
		return err
	}
	defer dbconnect.CloseDB()
	cmd.RootCmd.Execute()
	return nil
}

func init() {
	log.SetPrefix("LOG: ")
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
