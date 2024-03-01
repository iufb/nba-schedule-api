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

func (s *PostgresStore) Init() error {
	return s.CreateAccountTable()
}

func (s *PostgresStore) CreateAccountTable() error {
	query := `create table if not exists accounts  (
       id SERIAL PRIMARY KEY,
       username varchar(50),
       timezone varchar(6),
       createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP
 )`
	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStore) CreateAccount(acc *Account) error {
	query := `
INSERT INTO accounts ( username,timezone)
VALUES ($1,$2);
    `
	_, err := s.db.Exec(query, acc.Username, acc.Timezone)
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

func scanIntoAccount(r *sql.Row) (*Account, error) {
	acc := &Account{}
	err := r.Scan(&acc.Id, &acc.Username, &acc.Timezone, &acc.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("No account found.")
		}
		fmt.Printf("Error while get account %s", err)
	}
	return acc, nil
}
