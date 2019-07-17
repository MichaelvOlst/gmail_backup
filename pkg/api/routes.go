package api

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Routes starts the server
func (a *API) Routes() *echo.Echo {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.File("/", "public/index.html")

	e.POST("/auth/login", a.LoginHandler)
	e.POST("/auth/logout", a.LogoutHandler)

	return e
}
