package stores

import (
	"context"
	"database/sql"

	"github.com/alesbrelih/go-reservation-api/db"
	"github.com/alesbrelih/go-reservation-api/pkg/myutil"
	"github.com/pkg/errors"
)

type AuthStore interface {
	Authenticate(ctx context.Context, username string, password string) (int64, error)
	HasAccess(ctx context.Context, id int64) error
}

func NewAuthStoreSql(db db.DbFactory) AuthStore {
	return &authStoreSql{
		dbFactory: db,
	}
}

type authStoreSql struct {
	dbFactory db.DbFactory
}

func (a *authStoreSql) Authenticate(ctx context.Context, username string, password string) (int64, error) {

	db := a.dbFactory.Connect()
	defer db.Close()

	q := "SELECT id, pass FROM reservation_user WHERE username = $1"
	rows := db.QueryRowContext(ctx, q, username)

	var id int64
	var passDb string
	err := rows.Scan(&id, &passDb)
	if err != nil {
		return 0, errors.Wrap(err, "Error scanning id to variable")
	}

	if !myutil.CheckPasswordHash(passDb, password) {
		return 0, errors.Wrap(sql.ErrNoRows, "Passwords do not match")
	}

	return id, nil
}

// Validates if user with given id is still active / exists
func (a *authStoreSql) HasAccess(ctx context.Context, id int64) error {

	db := a.dbFactory.Connect()
	defer db.Close()

	q := "SELECT id FROM reservation_user WHERE id = $1 AND active = true"
	rows := db.QueryRowContext(ctx, q, id)

	var _id int64
	err := rows.Scan(&_id)
	if err != nil {
		return errors.Wrap(err, "Error scaning id")
	}

	return nil
}
