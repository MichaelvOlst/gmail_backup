package cmd

import (
	"gmail_backup/pkg/api"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(serveCmd)
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Starts the server",
	Run: func(cmd *cobra.Command, args []string) {
		a := api.New(app.config)
		a.StartServer()

	},
}

// func main() {
// 	// Echo instance
// 	e := echo.New()

// 	// Middleware
// 	e.Use(middleware.Logger())
// 	e.Use(middleware.Recover())

// 	// Routes
// 	e.GET("/", hello)

// 	// Start server
// 	e.Logger.Fatal(e.Start(":1323"))
// }

// // Handler
// func hello(c echo.Context) error {
// 	return c.String(http.StatusOK, "Hello, World!  ")
// }
