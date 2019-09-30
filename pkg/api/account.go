package api

import (
	"encoding/json"
	"gmail_backup/pkg/models"
	"net/http"
	"strconv"

	"github.com/asdine/storm"
	"github.com/gorilla/mux"
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

	validateErrors := ac.Validate()
	if len(validateErrors) != 0 {
		return respond(w, http.StatusUnprocessableEntity, envelope{Error: validateErrors})
	}

	na, err := a.db.CreateAccount(&ac)
	if err != nil && err != storm.ErrAlreadyExists {
		return respond(w, http.StatusInternalServerError, envelope{Error: err})
	}

	if err == storm.ErrAlreadyExists {
		return respond(w, http.StatusUnprocessableEntity, envelope{Error: "Account already exists"})
	}

	err = a.account.Reset()
	if err != nil {
		return respond(w, http.StatusUnprocessableEntity, envelope{Error: err})
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

	err = a.account.Reset()
	if err != nil {
		return respond(w, http.StatusUnprocessableEntity, envelope{Error: err})
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

	err = a.account.Reset()
	if err != nil {
		return respond(w, http.StatusUnprocessableEntity, envelope{Error: err})
	}

	return respond(w, http.StatusNoContent, nil)
}
