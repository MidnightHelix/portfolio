package models

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	host     = "127.0.0.1"
	port     = 5432
	user     = "hex"
	password = "toor"
	dbname   = "urlShortener"
)

func InitDB() (*sql.DB, error) {
	var connectionString = fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var error error
	db, error := sql.Open("postgres", connectionString)

	if error != nil {
		return nil, error
	}

	stmt, error := db.Prepare("CREATE TABLE IF NOT EXISTS urlShortener(ID SERIAL PRIMARY KEY, URL TEXT NOT NULL);")

	if error != nil {
		return nil, error
	}

	_, error = stmt.Exec()

	if error != nil {
		return nil, error
	}

	return db, nil
}
