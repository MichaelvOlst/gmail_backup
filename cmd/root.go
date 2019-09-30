package cmd

import (
	"fmt"
	"gmail_backup/pkg/account"
	"gmail_backup/pkg/config"
	"gmail_backup/pkg/database"
	"gmail_backup/pkg/models"
	"gmail_backup/pkg/storage"
	"os"

	"github.com/asdine/storm"
	"github.com/labstack/gommon/log"
	"github.com/spf13/cobra"
)

// App start point of this app
type App struct {
	config  *config.Config
	db      *database.Store
	storage *storage.Storage
	account *account.Account
}

var configFile string

func init() {
	cobra.OnInitialize(initApp)
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "config file")
}

var app *App

func initApp() {
	err := config.Load(configFile)
	if err != nil {
		log.Errorf("Cannot load config : %v", err)
		os.Exit(1)
	}

	config, err := config.Parse()
	if err != nil {
		log.Errorf("Cannot load config: %v", err)
		os.Exit(1)
	}

	db, err := database.New(config)
	if err != nil {
		log.Fatalf("Could not open db: %v\n", err)
	}

	app = &App{}
	app.config = config
	app.db = db

	app.storage = storage.New()

	var s = models.Settings{}
	_, err = app.db.GetBytes("settings", "settings")

	if err == storm.ErrNotFound {
		s = models.Settings{
			StorageOptions: models.StorageOptions{
				Dropbox: models.Dropbox{
					StorageOption: models.StorageOption{
						Option: "dropbox",
						Name:   "Dropbox",
						Active: true,
					},
				},
				Ftp: models.Ftp{
					StorageOption: models.StorageOption{
						Option: "ftp",
						Name:   "Ftp",
						Active: true,
					},
				},
				GoogleDrive: models.GoogleDrive{
					StorageOption: models.StorageOption{
						Option: "google_drive",
						Name:   "Google Drive",
						Active: true,
					},
				},
			},
		}

		app.db.Set("settings", "settings", &s)
		app.storage.RegisterAll(&s)
	} else {
		err := app.db.Get("settings", "settings", &s)
		if err != nil {
			log.Fatalf("Could not load the settings: %v\n", err)
		}

	}
	app.storage.RegisterAll(&s)

	acc := account.New(app.db, app.storage, app.config)
	app.account = acc
	acc.Start()

	err = app.db.Drop(&models.Message{})
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println("done dropping")

}

var rootCmd = &cobra.Command{
	Use:   "Gmail backup",
	Short: "It will back up gmail accounts",
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		app.db.Close()
	},
}

// Execute ..
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Errorf("Error executing command: %v", err)
		os.Exit(1)
	}
}
