package dbb

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq" // allowed
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "charles"
	dbname   = "goweb1"
)

// Connect would connect to the postgres database
func Connect() (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		return nil, err
	}
	// defer db.Close()

	err = db.Ping()
	if err != nil {
		return nil, err

	}
	fmt.Println("Connection to database was successful")
	return db, err
}
