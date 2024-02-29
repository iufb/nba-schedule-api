package main

import (
	"time"
)

type Account struct {
	Id        int       `json:"id" `
	Username  string    `json:"username" `
	Timezone  string    `json:"timezone" `
	CreatedAt time.Time `json:"createdAt" `
	// FavouriteTeams []Team
}
type CreateAccountRequest struct {
	Username string
	Timezone string
}

type CreateAccountResponse struct {
	Status string `json:"status"`
}

func NewAccount(username string, timezone string) *Account {
	return &Account{
		Username: username, Timezone: timezone,
	}
}

type Team struct {
	Name string
}

type Schedule struct {
	date time.Time
	ot   Team
}
