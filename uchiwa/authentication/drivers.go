package authentication

import (
	"errors"
	"fmt"
	"strings"

	"github.com/palourde/crypt"

	// Supported schemas for hashed passwords
	_ "github.com/palourde/crypt/apr1_crypt"
	_ "github.com/palourde/crypt/md5_crypt"
	_ "github.com/palourde/crypt/sha256_crypt"
	_ "github.com/palourde/crypt/sha512_crypt"
)

// Advanced function allows a third party Identification driver
func (a *Config) Advanced(driver loginFn, driverName string) {
	a.DriverFn = driver
	a.DriverName = driverName

	initToken(a.Auth)
}

// None function sets the Config struct in order to disable authentication
func (a *Config) None() {
	a.DriverFn = none
	a.DriverName = "none"
}

// Simple function sets the Config struct in order to enable simple authentication based on provided user and pass
func (a *Config) Simple(u []User) {
	a.DriverFn = simple
	a.DriverName = "simple"

	users = u

	initToken(a.Auth)
}

// none represents the authentication driver when auth is disabled
func none(u, p string) (*User, error) {
	return &User{}, nil
}

// simple represents the simple authentication driver
func simple(u, p string) (*User, error) {
	for _, user := range users {
		if u != user.Username {
			continue
		}

		if strings.HasPrefix(user.Password, "{crypt}") {
			password := user.Password
			password = strings.Replace(password, "{crypt}", "", 1)
			crypter, err := crypt.NewFromHash(password)
			if err != nil {
				return &user, errors.New("Invalid hashed password")
			}

			return &user, crypter.Verify(password, []byte(p))
		}

		if p == user.Password {
			return &user, nil
		}
	}
	return &User{}, fmt.Errorf("invalid user '%s' or invalid password", u)
}
