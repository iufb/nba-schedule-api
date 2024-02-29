package main

import (
	"time"
)

type Account struct {
	Username       string
	Timezone       string
	FavouriteTeams []Team
}

func NewAccount(username string, timezone string, ft []Team) *Account {
	return &Account{
		Username: username, Timezone: timezone, FavouriteTeams: ft,
	}
}

type Team struct {
	Name string
}

type Schedule struct {
	date time.Time
	ot   Team
}
