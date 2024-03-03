package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*Account) error
	DeleteAccount(int) error
	UpdateAccount(*Account) error
	GetAccountById(int) (*Account, error)
	GetAccountByUsername(string) (*Account, error)
	AddTeam(*Team) error
	AddTeamToFavourite(int, string) error
	GetAccountFavouriteTeams(int) ([]*Team, error)
}
type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	conSrt := "user=iufb dbname=go-nba port=5433 password=1243 sslmode=disable"
	db, err := sql.Open("postgres", conSrt)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &PostgresStore{db: db}, nil
}

// var teams = [30]Team{
// 	{Name: "Atlanta Hawks", Abbr: "ATL"},
// 	{Name: "Boston Celtics", Abbr: "BOS"},
// 	{Name: "Brooklyn Nets", Abbr: "BKN"},
// 	{Name: "Charlotte Hornets", Abbr: "CHA"},
// 	{Name: "Chicago Bulls", Abbr: "CHI"},
// 	{Name: "Cleveland Cavaliers", Abbr: "CLE"},
// 	{Name: "Dallas Mavericks", Abbr: "DAL"},
// 	{Name: "Denver Nuggets", Abbr: "DEN"},
// 	{Name: "Detroit Pistons", Abbr: "DET"},
// 	{Name: "Golden State Warriors", Abbr: "GSW"},
// 	{Name: "Houston Rockets", Abbr: "HOU"},
// 	{Name: "Indiana Pacers", Abbr: "IND"},
// 	{Name: "LA Clippers", Abbr: "LAC"},
// 	{Name: "Los Angeles Lakers", Abbr: "LAL"},
// 	{Name: "Memphis Grizzlies", Abbr: "MEM"},
// 	{Name: "Miami Heat", Abbr: "MIA"},
// 	{Name: "Milwaukee Bucks", Abbr: "MIL"},
// 	{Name: "Minnesota Timberwolves", Abbr: "MIN"},
// 	{Name: "New Orleans Pelicans", Abbr: "NOP"},
// 	{Name: "New York Knicks", Abbr: "NYK"},
// 	{Name: "Oklahoma City Thunder", Abbr: "OKC"},
// 	{Name: "Orlando Magic", Abbr: "ORL"},
// 	{Name: "Philadelphia 76ers", Abbr: "PHI"},
// 	{Name: "Phoenix Suns", Abbr: "PHX"},
// 	{Name: "Portland Trail Blazers", Abbr: "POR"},
// 	{Name: "Sacramento Kings", Abbr: "SAC"},
// 	{Name: "San Antonio Spurs", Abbr: "SAS"},
// 	{Name: "Toronto Raptors", Abbr: "TOR"},
// 	{Name: "Utah Jazz", Abbr: "UTA"},
// 	{Name: "Washington Wizards", Abbr: "WAS"},
// }

func (s *PostgresStore) Init() error {
	err := s.CreateAccountTable()
	if err != nil {
		return err
	}
	err = s.CreateTeamTable()
	if err != nil {
		return err
	}
	err = s.CreateAccountTeamsTable()
	if err != nil {
		return err
	}
	return nil
}

func (s *PostgresStore) CreateAccountTable() error {
	query := `create table if not exists accounts  (
       id SERIAL PRIMARY KEY,
       username varchar(50) UNIQUE,
       encrypted_password varchar(100),
       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
 )`
	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStore) CreateTeamTable() error {
	query := ` create table if not exists teams (
       id SERIAL PRIMARY KEY,
       name varchar(50),
       abbr varchar(3) UNIQUE
    )
    `
	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStore) CreateAccountTeamsTable() error {
	query := ` create table if not exists account_teams(
       account_id INT,
       team_abbr varchar(3),
       PRIMARY KEY(account_id, team_abbr),
       FOREIGN KEY(account_id) REFERENCES accounts(id),
       FOREIGN KEY(team_abbr) REFERENCES teams(abbr)
    )
    `
	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStore) AddTeam(team *Team) error {
	query := `
 INSERT INTO teams (name,abbr)
VALUES ($1,$2);
      `
	_, err := s.db.Exec(query, team.Name, team.Abbr)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostgresStore) CreateAccount(acc *Account) error {
	query := `
INSERT INTO accounts ( username,encrypted_password)
VALUES ($1,$2);
    `
	_, err := s.db.Exec(query, acc.Username, acc.EncryptedPassword)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostgresStore) GetAccountById(id int) (*Account, error) {
	query := `
    select * from accounts where id = $1
    `
	row := s.db.QueryRow(query, id)
	return scanIntoAccount(row)
}

func (s *PostgresStore) GetAccountByUsername(username string) (*Account, error) {
	query := `
    select * from accounts where username = $1
    `
	row := s.db.QueryRow(query, username)
	return scanIntoAccount(row)
}

func (s *PostgresStore) UpdateAccount(acc *Account) error {
	return nil
}

func (s *PostgresStore) DeleteAccount(id int) error {
	query := `
    delete  from accounts where id=$1
    `
	_, err := s.db.Exec(query, id)
	return err
}

func (s *PostgresStore) AddTeamToFavourite(accountId int, abbr string) error {
	query := `
INSERT INTO account_teams(account_id,team_abbr)
VALUES ($1,$2);
    `
	_, err := s.db.Exec(query, accountId, abbr)
	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresStore) GetAccountFavouriteTeams(accountId int) ([]*Team, error) {
	query := `
    select name, abbr from account_teams join teams on account_teams.team_abbr = teams.abbr where account_id =$1    `
	rows, err := s.db.Query(query, accountId)
	if err != nil {
		return nil, err
	}
	teams := []*Team{}
	for rows.Next() {
		team := &Team{}
		err := rows.Scan(&team.Name, &team.Abbr)
		if err != nil {
			return nil, err
		}
		teams = append(teams, team)
	}
	return teams, nil
}

func scanIntoAccount(r *sql.Row) (*Account, error) {
	acc := &Account{}
	err := r.Scan(&acc.Id, &acc.Username, &acc.EncryptedPassword, &acc.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("No account found.")
		}
		fmt.Printf("Error while get account %s", err)
	}
	return acc, nil
}
