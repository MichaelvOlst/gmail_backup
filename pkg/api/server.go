package api

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// StartServer starts the server
func (a *API) StartServer() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.File("/", "public/index.html")

	e.POST("/auth/login", a.LoginHandler)
	e.POST("/auth/logout", a.LogoutHandler)

	address := a.cfg.Server.Host + ":" + a.cfg.Server.Port
	e.Logger.Fatal(e.Start(address))
}
