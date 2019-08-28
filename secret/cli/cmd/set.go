package cmd

import (
	"fmt"

	"github.com/gophercises/secret/vault"

	"github.com/spf13/cobra"
)

//setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Set a secret in your secret storage",
	Run: func(cmd *cobra.Command, args []string) {
		v := vault.File(encodingKey, secretsPath())
		key, value := args[0], args[1]
		err := v.Set(key, value)
		if err != nil {
			fmt.Println("value not able to set: ", err)
			return
		}
		fmt.Println("value set successfully")
	},
}

func init() {
	RootCmd.AddCommand(setCmd)
}
