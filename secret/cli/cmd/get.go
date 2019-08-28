package cmd

import (
	"fmt"

	"github.com/gophercises/secret/vault"
	"github.com/spf13/cobra"
)

//getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a secret from your secret storage",
	Run: func(cmd *cobra.Command, args []string) {
		v := vault.File(encodingKey, secretsPath())
		key := args[0]
		value, err := v.Get(key)
		if err != nil {
			fmt.Println("no value for the key: ", err)
			return
		}
		fmt.Printf("%s = %s\n", key, value)
	},
}

func init() {
	RootCmd.AddCommand(getCmd)
}
