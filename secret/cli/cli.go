package main

import (
	"log"

	"github.com/gophercises/secret/cli/cmd"
)

func main() {
	cmd.RootCmd.Execute()
}

func init() {
	log.SetPrefix("LOG: ")
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
