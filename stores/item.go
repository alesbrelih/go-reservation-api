package stores

import (
	"context"
	"database/sql"
	"time"

	"github.com/alesbrelih/go-reservation-api/db"
	"github.com/alesbrelih/go-reservation-api/models"
	"github.com/lib/pq"
)

type ItemStoreSql struct {
	db *sql.DB
}

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

	q := `SELECT i.id, i.title, i.show_from, i.show_to, i.price,
					idrp.id, idrp.date_from, idrp.date_to, idrp.price
				FROM item i
					LEFT JOIN item_date_range_price idrp ON (idrp.item_id = i.id) WHERE i.id = $1`

	rows, err := myDb.QueryContext(ctx, q, id)
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	var item *models.Item
	for rows.Next() {
		var itemId int64
		var title string
		var showFrom time.Time
		var showTo time.Time
		var price int64
		var pId sql.NullInt64
		var pDateFrom sql.NullTime
		var pDateTo sql.NullTime
		var pPrice sql.NullInt64

		err = rows.Scan(&itemId, &title, &showFrom, &showTo, &price, &pId, &pDateFrom, &pDateTo, &pPrice)
		if err != nil {
			return nil, err
		}

		if item == nil {
			item = &models.Item{
				Id:         itemId,
				Title:      &title,
				ShowFrom:   &showFrom,
				ShowTo:     &showTo,
				Price:      price,
				DatePrices: []models.ItemDatePrice{},
			}
		}
		if pId.Valid {
			item.DatePrices = append(item.DatePrices, models.ItemDatePrice{
				Id:       pId.Int64,
				DateFrom: pDateFrom.Time,
				DateTo:   pDateTo.Time,
				Price:    pPrice.Int64,
			})
		}
	}

	if item == nil {
		return nil, sql.ErrNoRows
	}

	return item, nil
}

func (u *ItemStoreSql) Create(ctx context.Context, item *models.Item) (int64, error) {
	myDb := db.Connect()
	defer myDb.Close()

	tx, err := myDb.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}

	var id int64
	q := "INSERT INTO item (title, show_from, show_to, price) VALUES ($1, $2, $3, $4) RETURNING id"
	err = tx.QueryRowContext(ctx, q, item.Title, item.ShowFrom, item.ShowTo, item.Price).Scan(&id)

	if err != nil {
		tx.Rollback()
		return 0, err
	}

	if len(item.DatePrices) != 0 {
		stmt, err := tx.PrepareContext(ctx, `INSERT INTO item_date_range_price (item_id, date_from, date_to, price)
				VALUES ($1, $2, $3, $4)`)
		defer stmt.Close()
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		for _, price := range item.DatePrices {
			_, err = stmt.ExecContext(ctx, id, price.DateFrom, price.DateTo, price.Price)
			if err != nil {
				tx.Rollback()
				return 0, err
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (u *ItemStoreSql) Update(ctx context.Context, item *models.Item) error {
	myDb := db.Connect()
	defer myDb.Close()

	tx, err := myDb.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	stmt := "UPDATE item SET title=$2, show_from=$3, show_to=$4, price=$5 WHERE id = $1"
	res, err := tx.ExecContext(ctx, stmt, item.Id, item.Title, item.ShowFrom, item.ShowTo, item.Price)
	if err != nil {
		tx.Rollback()
		return err
	}

	num, _ := res.RowsAffected()
	if num == 0 {
		return sql.ErrNoRows
	}

	// delete those that were removed
	ids := []interface{}{}
	for _, i := range item.DatePrices {
		if i.Id != 0 {
			ids = append(ids, i.Id)
		}
	}

	stmt = "DELETE FROM item_date_range_price WHERE item_id = $1 AND (id != ANY ($2) OR $3)"
	res, err = tx.ExecContext(ctx, stmt, item.Id, pq.Array(ids), len(ids) == 0)
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, i := range item.DatePrices {
		if i.Id != 0 {
			// update
			stmt = `UPDATE item_date_range_price
				 SET date_from = $2, date_to = $3, price = $4
				 WHERE id = $1`
			res, err = tx.ExecContext(ctx, stmt, i.Id, i.DateFrom, i.DateTo, i.Price)
		} else {
			// create
			stmt = "INSERT INTO item_date_range_price (item_id, date_from, date_to, price) VALUES ($1, $2, $3, $4)"
			res, err = tx.ExecContext(ctx, stmt, item.Id, i.DateFrom, i.DateTo, i.Price)
		}

		if err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
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
