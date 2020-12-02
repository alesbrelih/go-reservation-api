package db

import (
	"database/sql"
)

func NewDbFactory(postgresDsn string) DbFactory {
	return &dbFactory{
		dsn: postgresDsn,
	}
}

type DbFactory interface {
	Connect() *sql.DB
}

type dbFactory struct {
	dsn string
}

func (d *dbFactory) Connect() *sql.DB {

	db, err := sql.Open("postgres", d.dsn)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return db
}
