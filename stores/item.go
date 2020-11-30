package stores

import (
	"context"
	"database/sql"

	"github.com/alesbrelih/go-reservation-api/db"
	"github.com/alesbrelih/go-reservation-api/models"
)

type ItemStoreSql struct{}

func (u *ItemStoreSql) GetAll(ctx context.Context) (models.Items, error) {

	myDb := db.Connect()
	defer myDb.Close()

	query := "SELECT id, title, show_from, show_to, price FROM item"
	rows, err := myDb.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items models.Items

	for rows.Next() {
		var item models.Item

		err = rows.Scan(&item.Id, &item.Title, &item.ShowFrom, &item.ShowTo, &item.Price)
		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}
	return items, nil
}

func (u *ItemStoreSql) GetOne(ctx context.Context, id int64) (*models.Item, error) {

	myDb := db.Connect()
	defer myDb.Close()

	item := &models.Item{}

	stmt := "SELECT * FROM item WHERE id = $1"
	res := myDb.QueryRowContext(ctx, stmt, id)
	err := res.Scan(&item.Id, &item.Title, &item.ShowFrom, &item.ShowTo, &item.Price)

	if err != nil {
		return nil, err
	}

	return item, nil
}

func (u *ItemStoreSql) Create(item *models.Item) (int64, error) {
	myDb := db.Connect()
	defer myDb.Close()

	stmt := "INSERT INTO item (title, show_from, show_to, price) VALUES ($1, $2, $3, $4)"
	res, err := myDb.Exec(stmt, item.Title, item.ShowFrom, item.ShowTo, item.Price)

	if err != nil {
		return 0, err
	}

	id, _ := res.LastInsertId()
	return id, nil
}

func (u *ItemStoreSql) Update(item *models.Item) error {
	myDb := db.Connect()
	defer myDb.Close()

	stmt := "UPDATE item SET title=$2, show_from=$3, show_to=$4, price=$5 WHERE id = $1"
	res, err := myDb.Exec(stmt, item.Id, item.Title, item.ShowFrom, item.ShowTo, item.Price)

	if err != nil {
		return err
	}

	num, _ := res.RowsAffected()
	if num == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (u *ItemStoreSql) Delete(id int64) error {
	myDb := db.Connect()
	defer myDb.Close()

	stmt := "DELETE FROM item WHERE id = $1"
	res, err := myDb.Exec(stmt, int64(id))
	if err != nil {
		return err
	}

	num, _ := res.RowsAffected()
	if num == 0 {
		return sql.ErrNoRows
	}

	return nil
}
