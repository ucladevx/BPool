package postgres

import (
	"fmt"

	"github.com/ucladevx/BPool/interfaces"

	"github.com/jmoiron/sqlx"
)

// NewConnection creates a new connection pool to a Postgres DB
func NewConnection(user, password, name, port, host string, l interfaces.Logger) *sqlx.DB {
	dbConnection := fmt.Sprintf("user=%s password=%s dbname=%s port=%s host=%s sslmode=disable",
		user, password, name, port, host)

	l.Info("DB CONN", "host", host, "port", port)

	var db *sqlx.DB
	var err error
	if db, err = sqlx.Connect("postgres", dbConnection); err != nil {
		l.Panic("DB CONN FAILED", "error", err.Error())
	}

	return db
}
