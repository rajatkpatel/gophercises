package cmd

import (
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

//SecretFile variable declared globally.
var SecretFile = ".secrets"

//RootCmd is the initial command.
var RootCmd = &cobra.Command{
	Use:   "secret",
	Short: "Secret is an API key and other secrets manager",
}

var encodingKey string

func init() {
	RootCmd.PersistentFlags().StringVarP(&encodingKey, "key", "k", "", "use for encoding and decoding secrets")
}

//Set the secret file in the home directory
func secretsPath() string {
	dir, _ := homedir.Dir()
	return filepath.Join(dir, SecretFile)
}
