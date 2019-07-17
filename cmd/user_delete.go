package cmd

import (
	"fmt"

	"github.com/labstack/gommon/log"
	"github.com/spf13/cobra"
)

// deleteUserCmd represents the new command
var deleteUserCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a user",
	Run: func(cmd *cobra.Command, args []string) {

		err := app.db.DeleteUser(email)

		if err != nil {
			log.Fatalf("Could not delete user: %s\n", err)
		}

		fmt.Println("User deleted")
	},
}

func init() {

	userCmd.AddCommand(deleteUserCmd)

	deleteUserCmd.PersistentFlags().StringVarP(&email, "email", "e", "", "Email of the user")
}
