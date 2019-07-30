package api

import (
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Routes starts the server
func (a *API) Routes() *echo.Echo {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(session.Middleware(sessions.NewCookieStore([]byte(a.config.Server.Secret))))

	// indexFile, err := a.box.FindString("index.html")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// appJS, err := a.box.FindString("/public/js/app.js")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	e.GET("/*", func(c echo.Context) error {
		// filename := c.Request().RequestURI
		b, err := a.box.Find("index.html")
		if err != nil {
			return err
		}
		return c.Blob(http.StatusOK, "text/html", b)
	})

	e.GET("/js/app.js", func(c echo.Context) error {
		// filename := c.Request().RequestURI
		b, err := a.box.Find("js/app.js")
		if err != nil {
			return err
		}
		return c.Blob(http.StatusOK, "text/javascript", b)
	})

	// e.File("/", indexFile)
	// e.File("/js/app.js", appJS)

	e.POST("/auth/login", a.LoginHandler)
	e.POST("/auth/logout", a.LogoutHandler)
	e.GET("/auth/session", a.SessionHandler)

	api := e.Group("/api")
	api.Use(a.apiMiddleware)

	api.GET("/accounts", a.GetAccountsHandler)
	// api.GET("/account/:id/backup", a.BackupAccount)
	api.POST("/account", a.CreateAccountHandler)
	api.PATCH("/account/:id", a.UpdateAccountHandler)
	api.DELETE("/account/:id", a.DeleteAccountHandler)

	// e.File("/*", indexFile)
	return e
}
