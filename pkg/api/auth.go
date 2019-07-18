package api

import (
	"net/http"
	"strings"

	"github.com/asdine/storm"
	"github.com/labstack/echo/v4"

	"github.com/labstack/echo-contrib/session"
)

type login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (l *login) Sanitize() {
	l.Email = strings.ToLower(strings.TrimSpace(l.Email))
}

// LoginHandler starts the server
func (a *API) LoginHandler(c echo.Context) error {
	var l login
	if err := c.Bind(&l); err != nil {
		return err
	}
	l.Sanitize()

	u, err := a.db.GetUserByEmail(l.Email)
	if err != nil && err != storm.ErrNotFound {
		return err
	}

	if err == storm.ErrNotFound || u.ComparePassword(l.Password) != nil {
		return c.JSON(http.StatusUnauthorized, envelope{Error: "invalid_credentials"})
	}

	// ignore error here as we want a (new) session regardless
	session, _ := session.Get("auth", c)
	session.Values["user_id"] = u.ID
	err = session.Save(c.Request(), c.Response())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, envelope{Error: "Could not save session"})
	}
	return c.JSON(http.StatusOK, envelope{Result: true})
}

// LogoutHandler starts the server
func (a *API) LogoutHandler(c echo.Context) error {

	session, _ := session.Get("auth", c)
	if !session.IsNew {
		session.Options.MaxAge = -1
		err := session.Save(c.Request(), c.Response())
		if err != nil {
			return err
		}
	}
	return c.JSON(http.StatusOK, envelope{Result: true})
}
