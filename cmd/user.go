package cmd

import (
	"github.com/spf13/cobra"
)

var name string
var email string
var password string

func init() {
	rootCmd.AddCommand(userCmd)
}

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Create or delete an user",
	Run: func(cmd *cobra.Command, args []string) {

	},
}
