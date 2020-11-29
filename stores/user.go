package stores

import (
	"context"
	"database/sql"
	"errors"

	"github.com/alesbrelih/go-reservation-api/db"
	"github.com/alesbrelih/go-reservation-api/models"
	"github.com/alesbrelih/go-reservation-api/pkg/myutil"
)

var PasswordMissmatch = errors.New("Passwords missmatch")

type UserStoreSql struct{}

func (u *UserStoreSql) GetAll(ctx context.Context) (models.Users, error) {

	myDb := db.Connect()
	defer myDb.Close()

	query := "SELECT id, first_name, last_name, username, email FROM reservations.public.user"
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

func (u *UserStoreSql) GetOne(ctx context.Context, id int64) (*models.User, error) {
	myDb := db.Connect()
	defer myDb.Close()

	user := &models.User{}

	stmt := "SELECT id, first_name, last_name, username, email FROM reservations.public.user WHERE id = $1"
	res := myDb.QueryRowContext(ctx, stmt, id)
	err := res.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Username, &user.Email)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserStoreSql) Create(user *models.UserReqBody) (int64, error) {
	myDb := db.Connect()
	defer myDb.Close()

	stmt := `INSERT INTO reservations.public.user (first_name, last_name, username, email, pass) 
		VALUES ($1, $2, $3, $4, $5)`

	password, err := u.setPassword(user)
	if err != nil {
		return 0, err
	}
	res, err := myDb.Exec(stmt, user.FirstName, user.LastName, user.Username, user.Email, password)
	if err != nil {
		return 0, err
	}

	id, _ := res.LastInsertId()
	return id, nil
}

func (u *UserStoreSql) Update(user *models.UserReqBody) error {
	myDb := db.Connect()
	defer myDb.Close()

	stmt := `UPDATE reservations.public.user 
			SET first_name = $2,
			last_name = $3,
			email = $4,
			pass = COALESCE($5, passw) WHERE id = $1`

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

func (u *UserStoreSql) setPassword(user *models.UserReqBody) (*string, error) {
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

func (u *UserStoreSql) Delete(id int64) error {
	myDb := db.Connect()
	defer myDb.Close()

	stmt := "DELETE FROM reservations.public.user WHERE id = $1"

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
