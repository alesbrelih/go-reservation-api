package stores

import (
	"context"
	"database/sql"

	"github.com/alesbrelih/go-reservation-api/db"
	"github.com/alesbrelih/go-reservation-api/models"
)

func NewTenantStore(db db.DbFactory) TenantStore {
	return &tenantStoreSql{db: db}
}

type TenantStore interface {
	GetAll(ctx context.Context) (models.Tenants, error)
	GetOne(ctx context.Context, id int64) (*models.Tenant, error)
	Create(*models.Tenant) (int64, error)
	Update(*models.Tenant) error
	Delete(id int64) error
}

type tenantStoreSql struct {
	db db.DbFactory
}

func (t *tenantStoreSql) GetAll(ctx context.Context) (models.Tenants, error) {

	myDb := t.db.Connect()
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

func (t *tenantStoreSql) GetOne(ctx context.Context, id int64) (*models.Tenant, error) {

	myDb := t.db.Connect()
	defer myDb.Close()

	stmt := `SELECT t.id, title, t.email, ru.id, ru.first_name, ru.last_name, ru.email
			FROM tenant t
				LEFT JOIN tenant_has_reservation_user thru ON (thru.tenant_id = t.id)
				LEFT JOIN reservation_user ru ON (ru.id = thru.reservation_user_id)
			WHERE t.id = $1`

	rows, err := myDb.QueryContext(ctx, stmt, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tenant *models.Tenant

	for rows.Next() {
		var id int64
		var title string
		var email string
		var userId int64
		var userFirstName string
		var userLastName string
		var userEmail string

		err := rows.Scan(&id, &title, &email, &userId, &userFirstName, &userLastName, &userEmail)
		if err != nil {
			return nil, err
		}

		if tenant == nil {
			tenant = &models.Tenant{
				Id:    id,
				Title: title,
				Email: email,
				Users: []models.User{},
			}
		}
		if userId != 0 {
			tenant.Users = append(tenant.Users, models.User{
				Id:        userId,
				FirstName: userFirstName,
				LastName:  userLastName,
				Email:     userEmail,
			})
		}

	}

	return tenant, nil
}

func (t *tenantStoreSql) Create(item *models.Tenant) (int64, error) {
	myDb := t.db.Connect()
	defer myDb.Close()

	stmt := "INSERT INTO tenant (title, email) VALUES ($1, $2, $3)"
	res, err := myDb.Exec(stmt, item.Title, item.Email)

	if err != nil {
		return 0, err
	}

	id, _ := res.LastInsertId()
	return id, nil
}

func (t *tenantStoreSql) Update(item *models.Tenant) error {
	myDb := t.db.Connect()
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

func (t *tenantStoreSql) Delete(id int64) error {
	myDb := t.db.Connect()
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
