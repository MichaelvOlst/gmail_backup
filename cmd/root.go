package cmd

import (
	"gmail_backup/pkg/config"
	"gmail_backup/pkg/database"
	"gmail_backup/pkg/models"
	"gmail_backup/pkg/storage"
	"os"
	"path/filepath"

	"github.com/asdine/storm"
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
	// fmt.Printf("%+v\n", &s)
	app.storage.RegisterAll(&s)
	// fmt.Printf("%+v\n", app.storage.GetProviders())

	// fmt.Println(app.storage.GetProviders())
	absFile, err := filepath.Abs("./dropbox.txt")
	// fmt.Println(absFile)
	if err != nil {
		log.Fatalf("Could not find the file")
	}
	app.storage.GetProvider("dropbox").Put(absFile)

	// app.storage.GetProvider("dropbox").ListFolder()

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
