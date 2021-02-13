package data

import (
	"github.com/upper/db/v4"
	"github.com/upper/db/v4/adapter/postgresql"
)

// NewPostgresqlSession creates a upper db.Session to perform operations on a postgresql database
func NewPostgresqlSession() (db.Session, error) {
	var settings = postgresql.ConnectionURL{
		User:     "postgres",
		Password: "mysecretpassword",
		Host:     "127.0.0.1:5432",
		Database: "postgres",
		Options:  map[string]string{"connect_timeout": "10"}, // https://www.postgresql.org/docs/10/libpq-connect.html#LIBPQ-PARAMKEYWORDS
	}

	return postgresql.Open(settings)
}
