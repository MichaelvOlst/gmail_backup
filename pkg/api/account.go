package api

import (
	"encoding/json"
	"gmail_backup/pkg/models"
	"net/http"
	"strconv"
	"time"

	"github.com/asdine/storm"
	"github.com/gorilla/mux"

	"github.com/labstack/echo/v4"
)

// HandlerGetAllAccounts Gets all accounts
func (a *API) HandlerGetAllAccounts(w http.ResponseWriter, r *http.Request) error {
	var accounts []models.Account
	err := a.db.All(&accounts)
	if err != nil {
		return respond(w, http.StatusInternalServerError, envelope{Error: err})
	}
	return respond(w, http.StatusOK, envelope{Result: accounts})
}

// HandlerGetSingleAccount Gets all accounts
func (a *API) HandlerGetSingleAccount(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	aid, ok := vars["id"]
	if !ok {
		return respond(w, http.StatusUnprocessableEntity, envelope{Error: "Id not valid"})
	}
	id, _ := strconv.Atoi(aid)

	ac, err := a.db.GetAccountByID(id)
	if err != nil && err != storm.ErrNotFound {
		return respond(w, http.StatusInternalServerError, envelope{Error: err})
	}

	if err == storm.ErrNotFound {
		return respond(w, http.StatusUnprocessableEntity, envelope{Error: "Account not found"})
	}

	return respond(w, http.StatusOK, envelope{Result: ac})
}

// HandlerCreateAccount Creates an account
func (a *API) HandlerCreateAccount(w http.ResponseWriter, r *http.Request) error {
	var ac models.Account

	err := json.NewDecoder(r.Body).Decode(&ac)
	if err != nil {
		return err
	}

	na, err := a.db.CreateAccount(&ac)
	if err != nil && err != storm.ErrAlreadyExists {
		return respond(w, http.StatusInternalServerError, envelope{Error: err})
	}

	if err == storm.ErrAlreadyExists {
		return respond(w, http.StatusUnprocessableEntity, envelope{Error: "Account already exists"})
	}

	return respond(w, http.StatusOK, envelope{Result: na})
}

// HandlerUpdateAccount Updates an account
func (a *API) HandlerUpdateAccount(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	aid, ok := vars["id"]
	if !ok {
		return respond(w, http.StatusUnprocessableEntity, envelope{Error: "Id not valid"})
	}
	id, _ := strconv.Atoi(aid)

	ac, err := a.db.GetAccountByID(id)
	if err != nil && err != storm.ErrNotFound {
		return respond(w, http.StatusInternalServerError, envelope{Error: err})
	}

	if err == storm.ErrNotFound {
		return respond(w, http.StatusUnprocessableEntity, envelope{Error: "Account not found"})
	}

	var uc models.Account
	err = json.NewDecoder(r.Body).Decode(&uc)
	if err != nil {
		return respond(w, http.StatusInternalServerError, envelope{Error: err})
	}

	uc.ID = ac.ID
	nc, err := a.db.UpdateAccount(&uc)
	if err != nil {
		return respond(w, http.StatusInternalServerError, envelope{Error: err})
	}

	return respond(w, http.StatusOK, envelope{Result: nc})
}

// HandlerDeleteAccount Deletes an account
func (a *API) HandlerDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	aid, ok := vars["id"]
	if !ok {
		return respond(w, http.StatusUnprocessableEntity, envelope{Error: "Id not valid"})
	}
	id, _ := strconv.Atoi(aid)

	ac, err := a.db.GetAccountByID(id)
	if err != nil && err != storm.ErrNotFound {
		return respond(w, http.StatusInternalServerError, envelope{Error: err})
	}

	if err == storm.ErrNotFound {
		return respond(w, http.StatusUnprocessableEntity, envelope{Error: "Account not found"})
	}

	err = a.db.DeleteAccount(ac)
	if err != nil {
		return respond(w, http.StatusInternalServerError, envelope{Error: err})
	}

	return respond(w, http.StatusNoContent, nil)
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
