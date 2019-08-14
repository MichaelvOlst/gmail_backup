package cmd

import (
	"gmail_backup/pkg/config"
	"gmail_backup/pkg/database"
	"gmail_backup/pkg/models"
	"gmail_backup/pkg/storage"
	"gmail_backup/pkg/storage/drive"
	"gmail_backup/pkg/storage/dropbox"
	"gmail_backup/pkg/storage/ftp"
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

	s := models.Settings{StorageOptions: []models.StorageOptions{
		models.StorageOptions{
			Name:     "ftp",
			Provider: ftp.New(),
			Active:   true,
			Config:   ftp.Config{},
		},
		models.StorageOptions{
			Name:     "dropbox",
			Provider: dropbox.New(),
			Active:   true,
			Config:   dropbox.Config{},
		},
		models.StorageOptions{
			Name:     "google_drive",
			Provider: drive.New(),
			Active:   true,
			Config:   drive.Config{},
		},
	}}
	app.db.Set("settings", "settings", &s)

	app.storage = storage.New()

	for _, val := range s.StorageOptions {
		if val.Active {
			app.storage.Register(val.Provider)
		}
	}

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
