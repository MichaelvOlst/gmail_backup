package cmd

import (
	"gmail_backup/pkg/config"
	"os"

	"github.com/labstack/gommon/log"
	"github.com/spf13/cobra"
)

// App start point of this app
type App struct {
	config *config.Config
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

	app = &App{}
	app.config = config
}

var rootCmd = &cobra.Command{
	Use:   "Gmail backup",
	Short: "It will back up gmail accounts",
}

// Execute ..
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Errorf("Error executing command: %v", err)
		os.Exit(1)
	}
}
