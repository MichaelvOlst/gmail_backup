package api

import (
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

	e.File("/", "public/index.html")

	e.POST("/auth/login", a.LoginHandler)
	e.POST("/auth/logout", a.LogoutHandler)

	api := e.Group("/api")
	api.Use(a.apiMiddleware)

	api.GET("/accounts", a.GetAccountsHandler)
	api.POST("/account", a.CreateAccountHandler)
	api.PATCH("/account/:id", a.UpdateAccountHandler)
	api.DELETE("/account/:id", a.DeleteAccountHandler)

	return e
}
