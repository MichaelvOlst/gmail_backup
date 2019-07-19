package api

import (
	"encoding/json"
	"gmail_backup/pkg/models"
	"net/http"
	"strconv"
	"time"

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

	ac, err := a.db.GetAccountByID(id)
	if err != nil && err != storm.ErrNotFound {
		return c.JSON(http.StatusInternalServerError, envelope{Error: err})
	}

	if err == storm.ErrNotFound {
		return c.JSON(http.StatusUnprocessableEntity, envelope{Error: "Account not found"})
	}

	var uc models.Account
	if err := c.Bind(&uc); err != nil {
		return c.JSON(http.StatusInternalServerError, envelope{Error: err})
	}

	uc.ID = ac.ID
	nc, err := a.db.UpdateAccount(&uc)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, envelope{Error: err})
	}

	return c.JSON(http.StatusOK, envelope{Result: nc})
}

// DeleteAccountHandler Deletes an account
func (a *API) DeleteAccountHandler(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	ac, err := a.db.GetAccountByID(id)
	if err != nil && err != storm.ErrNotFound {
		return c.JSON(http.StatusInternalServerError, envelope{Error: err})
	}

	if err == storm.ErrNotFound {
		return c.JSON(http.StatusUnprocessableEntity, envelope{Error: "Account not found"})
	}

	err = a.db.DeleteAccount(ac)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, envelope{Error: err})
	}

	return c.JSON(http.StatusNoContent, nil)
}

// BackupAccount backing up an account
func (a *API) BackupAccount(c echo.Context) error {

	type (
		Geolocation struct {
			Altitude  float64
			Latitude  float64
			Longitude float64
		}
	)

	var (
		locations = []Geolocation{
			{-97, 37.819929, -122.478255},
			{1899, 39.096849, -120.032351},
			{2619, 37.865101, -119.538329},
			{42, 33.812092, -117.918974},
			{15, 37.77493, -122.419416},
		}
	)
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c.Response().WriteHeader(http.StatusOK)
	for _, l := range locations {
		if err := json.NewEncoder(c.Response()).Encode(l); err != nil {
			return err
		}
		c.Response().Flush()
		time.Sleep(1 * time.Second)
	}
	return nil
}
