package postgres

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

func NewConnection(user, password, name, port, host string) *sqlx.DB {
	dbConnection := fmt.Sprintf("user=%s password=%s dbname=%s port=%s host=%s sslmode=disable",
		user, password, name, port, host)

	var db *sqlx.DB
	var err error
	if db, err = sqlx.Connect("postgres", dbConnection); err != nil {
		log.Fatal("❌  DB CONN: " + err.Error())
	}

	return db
}
