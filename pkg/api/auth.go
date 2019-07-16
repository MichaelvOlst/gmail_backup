package api

import (
	"net/http"

	"github.com/labstack/echo"
)

// User ..
type User struct {
	Name  string `json:"name" xml:"name"`
	Email string `json:"email" xml:"email"`
}

// LoginHandler starts the server
func (a *API) LoginHandler(c echo.Context) error {
	u := &User{
		Name:  "Michael",
		Email: "michael@cmspecialist.nl",
	}
	return c.JSON(http.StatusOK, u)
}

// LogoutHandler starts the server
func (a *API) LogoutHandler(c echo.Context) error {
	u := struct {
		Result string `json:"result"`
	}{
		Result: "Logged out",
	}
	return c.JSON(http.StatusOK, u)
}
