package stores

import (
	"context"
	"database/sql"
	"errors"

	"github.com/alesbrelih/go-reservation-api/db"
	"github.com/alesbrelih/go-reservation-api/models"
	"github.com/alesbrelih/go-reservation-api/pkg/myutil"
)

func NewUserStore(db db.DbFactory) UserStore {
	return &userStoreSql{db: db}
}

var PasswordMissmatch = errors.New("Passwords missmatch")

type UserStore interface {
	GetAll(ctx context.Context) (models.Users, error)
	GetOne(ctx context.Context, id int64) (*models.User, error)
	Create(*models.UserReqBody) (int64, error)
	Update(*models.UserReqBody) error
	Delete(id int64) error
}

type userStoreSql struct {
	db db.DbFactory
}

func (u *userStoreSql) GetAll(ctx context.Context) (models.Users, error) {

	myDb := u.db.Connect()
	defer myDb.Close()

	query := "SELECT id, first_name, last_name, username, email FROM reservation_user"
	rows, err := myDb.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users models.Users

	for rows.Next() {
		var user models.User

		err := rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Username, &user.Email)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil

}

func (u *userStoreSql) GetOne(ctx context.Context, id int64) (*models.User, error) {
	myDb := u.db.Connect()
	defer myDb.Close()

	user := &models.User{}

	stmt := "SELECT id, first_name, last_name, username, email FROM reservation_user WHERE id = $1"
	res := myDb.QueryRowContext(ctx, stmt, id)
	err := res.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Username, &user.Email)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *userStoreSql) Create(user *models.UserReqBody) (int64, error) {
	myDb := u.db.Connect()
	defer myDb.Close()

	stmt := `INSERT INTO reservation_user (first_name, last_name, username, email, pass) 
		VALUES ($1, $2, $3, $4, $5)`

	password, err := u.setPassword(user)
	if err != nil {
		return 0, err
	}
	res, err := myDb.Exec(stmt, user.FirstName, user.LastName, user.Username, user.Email, password)
	if err != nil {
		return 0, err
	}

	// FIX ME: (postgre) returns 0 alwys, need to use queryrow and scan
	id, _ := res.LastInsertId()
	return id, nil
}

func (u *userStoreSql) Update(user *models.UserReqBody) error {
	myDb := u.db.Connect()
	defer myDb.Close()

	stmt := `UPDATE reservation_user 
			SET first_name = $2,
			last_name = $3,
			email = $4,
			pass = COALESCE($5, pass) WHERE id = $1`

	password, err := u.setPassword(user)
	if err != nil {
		return err
	}

	res, err := myDb.Exec(stmt, user.Id, user.FirstName, user.LastName, user.Email, password)
	if err != nil {
		return err
	}

	num, _ := res.RowsAffected()
	if num == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (u *userStoreSql) setPassword(user *models.UserReqBody) (*string, error) {
	if myutil.Empty(user.Password) && myutil.Empty(user.Confirm) {
		return nil, nil
	}
	if user.Password != user.Confirm {
		return nil, PasswordMissmatch
	}

	hash, err := myutil.HashPassword(user.Password, 14)
	if err != nil {
		return nil, err
	}
	return &hash, nil
}

func (u *userStoreSql) Delete(id int64) error {
	myDb := u.db.Connect()
	defer myDb.Close()

	stmt := "DELETE FROM reservation_user WHERE id = $1"

	res, err := myDb.Exec(stmt, id)
	if err != nil {
		return err
	}

	num, _ := res.RowsAffected()
	if num == 0 {
		return sql.ErrNoRows
	}

	return nil
}
