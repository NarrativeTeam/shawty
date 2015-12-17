package storages

import (
	"database/sql"
	"fmt"
	"math/rand"

	_ "github.com/lib/pq"
)

var runes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789")

func GetRandomToken() string {
	b := make([]rune, 6)
	for i := range b {
		b[i] = runes[rand.Intn(len(runes))]
	}
	return string(b)
}

func NewPostgres(host, dbName, user, password string, useSSL bool) (*Postgres, error) {
	pgConn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s", host, user, password, dbName)
	if !useSSL {
		pgConn += " sslmode=disable"
	}
	db, err := sql.Open("postgres", pgConn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return &Postgres{db}, err
}

type DB interface {
	Exec(string, ...interface{}) (sql.Result, error)
	QueryRow(string, ...interface{}) *sql.Row
	Close() error
}

type Postgres struct {
	db DB
}

func (ps *Postgres) Save(url string) (token string, err error) {
	token = GetRandomToken()
	_, err = ps.db.Exec(INSERT_URL_SQL, url, token)
	return
}

func (ps *Postgres) Load(token string) (url string, err error) {
	row := ps.db.QueryRow(SELECT_URL_SQL, token)
	err = row.Scan(&url)
	return
}

func (ps *Postgres) Close() error {
	return ps.db.Close()
}

// Setup the database, create tables etc.
func (ps *Postgres) Setup() {
	_, err := ps.db.Exec(CREATE_TABLES_SQL)
	if err != nil {
		panic(err)
	}
}

// Remove all table/data from database.
func (ps *Postgres) Teardown() {
	_, err := ps.db.Exec(DROP_TABLES_SQL)
	if err != nil {
		panic(err)
	}
}