package cmd

import (
	"fmt"
	"gmail_backup/pkg/models"

	"github.com/labstack/gommon/log"
	"github.com/spf13/cobra"
)

// newUserCmd represents the new command
var newUserCmd = &cobra.Command{
	Use:   "new",
	Short: "Creating a new user",
	Run: func(cmd *cobra.Command, args []string) {

		u := models.NewUser(email, name, password)
		u, err := app.db.CreateUser(u)
		if err != nil {
			log.Fatalf("Could not create user: %s\n", err)
		}

		fmt.Println("Created user")
	},
}

func init() {

	userCmd.AddCommand(newUserCmd)

	newUserCmd.PersistentFlags().StringVarP(&name, "name", "n", "", "Name of the user")
	newUserCmd.PersistentFlags().StringVarP(&email, "email", "e", "", "Email of the user")
	newUserCmd.PersistentFlags().StringVarP(&password, "password", "p", "", "Password for this user")
}
