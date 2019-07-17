package cmd

import (
	"gmail_backup/pkg/config"
	"os"

	"github.com/asdine/storm"
	"github.com/labstack/gommon/log"
	"github.com/spf13/cobra"
)

// App start point of this app
type App struct {
	config *config.Config
	db     *storm.DB
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
		log.Errorf("Cannot load config: %v", err)
		os.Exit(1)
	}

	config, err := config.Parse()
	if err != nil {
		log.Errorf("Cannot load config: %v", err)
		os.Exit(1)
	}

	db, err := storm.Open(config.Database.Filename)
	if err != nil {
		log.Fatalf("Could not open db: %v\n", err)
	}

	app = &App{}
	app.config = config
	app.db = db
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
