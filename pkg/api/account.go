package api

import (
	"gmail_backup/pkg/models"
	"net/http"
	"strconv"

	"github.com/asdine/storm"
	"github.com/labstack/echo/v4"
)

// GetAccountsHandler Gets all accounts
func (a *API) GetAccountsHandler(c echo.Context) error {
	var accounts []models.Account
	err := a.db.All(&accounts)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, envelope{Error: err})
	}
	return c.JSON(http.StatusOK, envelope{Result: accounts})
}

// CreateAccountHandler Creates an account
func (a *API) CreateAccountHandler(c echo.Context) error {
	var ac models.Account
	if err := c.Bind(&ac); err != nil {
		return c.JSON(http.StatusInternalServerError, envelope{Error: err})
	}

	na, err := a.db.CreateAccount(&ac)
	if err != nil && err != storm.ErrAlreadyExists {
		return c.JSON(http.StatusInternalServerError, envelope{Error: err})
	}

	if err == storm.ErrAlreadyExists {
		return c.JSON(http.StatusUnprocessableEntity, envelope{Error: "Account already exists"})
	}

	return c.JSON(http.StatusOK, envelope{Result: na})
}

// UpdateAccountHandler Updates an account
func (a *API) UpdateAccountHandler(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	return c.JSON(http.StatusOK, envelope{Result: "Update account " + id})
}

// DeleteAccountHandler Deletes an account
func (a *API) DeleteAccountHandler(c echo.Context) error {

	return c.JSON(http.StatusOK, envelope{Result: "Delete account"})
}
