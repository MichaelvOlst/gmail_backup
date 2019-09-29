package cmd

import (
	"context"
	"fmt"
	"gmail_backup/pkg/api"
	"gmail_backup/pkg/gmail"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gobuffalo/packr"
	"github.com/labstack/gommon/log"
	"github.com/spf13/cobra"
)

var googleConfigFile string

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.Flags().StringVarP(&googleConfigFile, "google", "g", "", "Google config file")
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Starts the server",
	Run: func(cmd *cobra.Command, args []string) {

		addr := fmt.Sprintf("%s:%s", app.config.Server.Host, app.config.Server.Port)
		fmt.Printf("Running app on http://%s\n", addr)

		box := packr.NewBox("./../public")

		g, err := gmail.New(app.config, app.db)
		if err != nil {
			log.Fatalf("Could not init gmail. error %v", err)
		}

		api := api.New(app.config, app.db, &box, app.storage, g)

		server := &http.Server{
			Addr:         addr,
			Handler:      api.Routes(),
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		}

		go func() {
			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				fmt.Printf("listen: %s\n", err)
			}
		}()

		quit := make(chan os.Signal)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		fmt.Println("Shutting down..")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			fmt.Printf("Could not shutdown the server %s", err)
		}

		app.cronjob.Stop()
	},
}
