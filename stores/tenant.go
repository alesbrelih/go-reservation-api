package stores

import (
	"context"
	"database/sql"

	"github.com/alesbrelih/go-reservation-api/db"
	"github.com/alesbrelih/go-reservation-api/models"
)

type TenantStoreSql struct{}

func (u *TenantStoreSql) GetAll(ctx context.Context) (models.Tenants, error) {

	myDb := db.Connect()
	defer myDb.Close()

	query := "SELECT id, title, email FROM tenant"
	rows, err := myDb.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items models.Tenants

	for rows.Next() {
		var item models.Tenant

		err = rows.Scan(&item.Id, &item.Title, &item.Email)
		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}
	return items, nil
}

func (u *TenantStoreSql) GetOne(ctx context.Context, id int64) (*models.Tenant, error) {

	myDb := db.Connect()
	defer myDb.Close()

	item := &models.Tenant{}

	stmt := "SELECT id, title, email FROM tenant WHERE id = $1"
	res := myDb.QueryRowContext(ctx, stmt, id)
	err := res.Scan(&item.Id, &item.Title, &item.Email)

	if err != nil {
		return nil, err
	}

	return item, nil
}

func (u *TenantStoreSql) Create(item *models.Tenant) (int64, error) {
	myDb := db.Connect()
	defer myDb.Close()

	stmt := "INSERT INTO tenant (title, email) VALUES ($1, $2, $3)"
	res, err := myDb.Exec(stmt, item.Title, item.Email)

	if err != nil {
		return 0, err
	}

	id, _ := res.LastInsertId()
	return id, nil
}

func (u *TenantStoreSql) Update(item *models.Tenant) error {
	myDb := db.Connect()
	defer myDb.Close()

	stmt := "UPDATE tenant SET title=$2, email=$3, show_to=$3 WHERE id = $1"
	res, err := myDb.Exec(stmt, item.Id, item.Title, item.Email)

	if err != nil {
		return err
	}

	num, _ := res.RowsAffected()
	if num == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (u *TenantStoreSql) Delete(id int64) error {
	myDb := db.Connect()
	defer myDb.Close()

	stmt := "DELETE FROM tenant WHERE id = $1"
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
