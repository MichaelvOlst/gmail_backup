package cmd

import (
	"context"
	"gmail_backup/pkg/api"
	"os"
	"os/signal"
	"time"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(serveCmd)
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Starts the server",
	Run: func(cmd *cobra.Command, args []string) {
		a := api.New(app.config, app.db)
		e := a.Routes()

		address := app.config.Server.Host + ":" + app.config.Server.Port

		go func(address string) {
			if err := e.Start(address); err != nil {
				e.Logger.Info("shutting down the server")
			}
		}(address)

		quit := make(chan os.Signal)
		signal.Notify(quit, os.Interrupt)
		<-quit
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := e.Shutdown(ctx); err != nil {
			e.Logger.Fatalf("Shutting down server: %v\n", err)
		}

	},
}
