package main

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type (
	RegisterRequest struct{ Auth }
	LoginRequest    struct {
		Auth
	}
)

type AddToFavourite struct {
	Abbr string `json:"abbr"`
}
type Account struct {
	Id                int       `json:"id" `
	Username          string    `json:"username" `
	EncryptedPassword string    `json:"-" `
	CreatedAt         time.Time `json:"createdAt" `
	// FavouriteTeams []Team
}
type CreateAccountRequest struct {
	Username string
	Timezone string
}

type WithStatusResponse struct {
	Status string `json:"status"`
}

func (acc *Account) ValidateAccount(pw string) bool {
	return bcrypt.CompareHashAndPassword([]byte(acc.EncryptedPassword), []byte(pw)) == nil
}

func NewAccount(username string, password string) (*Account, error) {
	encpw, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return &Account{
		Username:          username,
		EncryptedPassword: string(encpw),
	}, nil
}

type Team struct {
	Name string `json:"name"`
	Abbr string `json:"abbr"`
}

type Schedule struct {
	date time.Time
	ot   Team
}
