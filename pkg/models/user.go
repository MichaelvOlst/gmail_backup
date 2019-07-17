package models

import (
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// User ...
type User struct {
	ID       int    `json:"id" storm:"id,increment"`
	Name     string `json:"name"`
	Email    string `json:"email" storm:"unique"`
	Password string `json:"-"`
}

// NewUser creates a new User with the given email and password
func NewUser(e, n string, pwd string) *User {
	u := &User{
		Email: strings.ToLower(strings.TrimSpace(e)),
		Name:  n,
	}
	u.SetPassword(pwd)
	return u
}

// SetPassword sets a brcrypt encrypted password from the given plaintext pwd
func (u *User) SetPassword(pwd string) {
	hash, _ := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	u.Password = string(hash)
}

// ComparePassword returns true when the given plaintext password matches the encrypted pwd
func (u *User) ComparePassword(pwd string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(pwd))
}
