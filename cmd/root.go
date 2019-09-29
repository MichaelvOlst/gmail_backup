package cmd

import (
	"fmt"
	"gmail_backup/pkg/config"
	"gmail_backup/pkg/database"
	"gmail_backup/pkg/gmail"
	"gmail_backup/pkg/models"
	"gmail_backup/pkg/storage"
	"os"

	"github.com/asdine/storm"
	"github.com/labstack/gommon/log"
	"github.com/robfig/cron/v3"
	"github.com/spf13/cobra"
)

// App start point of this app
type App struct {
	config  *config.Config
	db      *database.Store
	storage *storage.Storage
	cronjob *cron.Cron
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
	// absFile, _ := filepath.Abs("./test.zip")
	// fmt.Println(absFile)
	// if err != nil {
	// 	log.Fatalf("Could not find the file")
	// }
	// p, _ := app.storage.GetProvider("dropbox")
	// p.Mkdir("/test")
	// p.Put("dropbox.txt", "/test/")

	c := cron.New()
	app.cronjob = c
	app.cronjob.Start()

	var accounts []models.Account
	err = app.db.All(&accounts)
	if err != nil {
		log.Fatal(err)
	}
	for _, ac := range accounts {

		fmt.Println(ac.CronExpression + " - " + ac.Email)

		// job := func(a *App, account models.Account) {

		// }(app, ac)

		id, err := app.cronjob.AddFunc(ac.CronExpression, func() {
			g, err := gmail.New(app.config, app.db)
			if err != nil {
				log.Fatal(err)
				return
			}

			fmt.Println("Backing up account " + ac.Email)
			g.Backup(ac, app.storage)
		})

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(id)
	}

	fmt.Printf("%+v\n", app.cronjob.Entries())

	// app.cronjob.AddFunc("* * * * *", func() {
	// 	date := time.Now().Format("15:04:05 02-01-2006")
	// 	fmt.Println(date + " before")
	// })

	// app.cronjob.AddFunc("* * * * *", func() {
	// 	date := time.Now().Format("15:04:05 02-01-2006")
	// 	fmt.Println(date + " after")
	// })

	// fmt.Printf("%+v\n", app.cronjob.Entries())

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
