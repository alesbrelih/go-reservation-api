package db

import (
	"database/sql"
	"os"

	"github.com/alesbrelih/go-reservation-api/pkg/myutil"
)

var dsn string

func init() {
	var found bool
	dsn, found = os.LookupEnv("POSTGRES_URL")
	if !found {
		panic(myutil.MissingEnvVariableMsg("DB_HOST"))
	}
}

func Connect() *sql.DB {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return db
}
