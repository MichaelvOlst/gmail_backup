package cmd

import (
	"gmail_backup/pkg/config"
	"gmail_backup/pkg/database"
	"gmail_backup/pkg/models"
	"gmail_backup/pkg/storage"
	"os"

	"github.com/labstack/gommon/log"
	"github.com/spf13/cobra"
)

// App start point of this app
type App struct {
	config  *config.Config
	db      *database.Store
	storage *storage.Storage
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

	s := models.Settings{
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

	app.storage = storage.New()

	app.storage.RegisterAll(&s)

	// for key, val := range s.StorageOptions {
	// 	fmt.Println(key)
	// 	fmt.Println(val)
	// 	// 	if val.Active {
	// 	// 		app.storage.Register(val.Option, val.Config)
	// 	// 	}
	// }

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
