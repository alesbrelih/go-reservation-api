package stores

import (
	"context"
	"database/sql"

	"github.com/alesbrelih/go-reservation-api/db"
	"github.com/alesbrelih/go-reservation-api/models"
	"github.com/pkg/errors"
)

func NewInquiryStore(dbFactory db.DbFactory) InquiryStore {
	return &inquiryStoreSql{
		dbFactory: dbFactory,
	}
}

type InquiryStore interface {
	GetAll(ctx context.Context) (models.Inquiries, error)
	Create(ctx context.Context, inquiry *models.InquiryCreate) error
	Delete(ctx context.Context, id int64) error
}

type inquiryStoreSql struct {
	dbFactory db.DbFactory
}

func (i *inquiryStoreSql) GetAll(ctx context.Context) (models.Inquiries, error) {
	db := i.dbFactory.Connect()
	defer db.Close()

	q := `SELECT inq.id, inq.inquirer, inq.email, inq.phone,
				inq.date_reservation, inq.date_created, inq.comment,
				i.id, i.title, i.price
			FROM inquiry inq 
				LEFT JOIN item i ON (i.id = inq.item_id)`

	rows, err := db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	inquiries := []models.Inquiry{}
	for rows.Next() {
		inquiry := models.Inquiry{Item: models.Item{}}
		err = rows.Scan(&inquiry.Id, &inquiry.Inquirer,
			&inquiry.Email, &inquiry.Phone, &inquiry.DateReservation,
			&inquiry.DateCreated, &inquiry.Comment, &inquiry.Item.Id,
			&inquiry.Item.Title, &inquiry.Item.Price,
		)

		if err != nil {
			return nil, err
		}

		inquiries = append(inquiries, inquiry)
	}

	return inquiries, nil
}

func (i *inquiryStoreSql) Create(ctx context.Context, inquiry *models.InquiryCreate) error {
	db := i.dbFactory.Connect()
	defer db.Close()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "Could not initialize Create inqiry transaction")
	}

	//TODO: validate when inserting item that dateprices dont overlap!!
	// AND THEY NEED TO BE REQUIRE

	item := &models.Item{}
	q := `SELECT i.id, i.title, COALESCE(ip.price, i.price) 
			FROM item i
				LEFT JOIN item_date_range_price ip ON (ip.item_id = i.id AND 
					ip.date_from <= now() at time zone 'utc' AND ip.date_to >= now() at time zone 'utc')
			WHERE i.id = $1`

	// SELECT CORRECT PRICE
	if err := tx.QueryRowContext(ctx, q, inquiry.ItemId).Scan(&item.Id, &item.Title, &item.Price); err != nil {
		return errors.Wrap(err, "Error retrieving item on inquiry create")
	}

	// get item that belongs to

	q = `INSERT INTO inquiry 
		(inquirer,email,phone,item_id, item_title, item_price, date_reservation,date_created)
		VALUES
		($1, $2, $3, $4, $5, $6, $7, now() at time zone 'utc')`

	_, err = tx.ExecContext(ctx, q, inquiry.Inquirer, inquiry.Email,
		inquiry.Phone, item.Id, item.Title, item.Price, inquiry.Date)
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "Error creating new inquiry")
	}

	tx.Commit()
	return nil
}

func (i *inquiryStoreSql) Delete(ctx context.Context, id int64) error {
	db := i.dbFactory.Connect()
	defer db.Close()

	q := "DELETE FROM inquiry WHERE id = $1"
	res, err := db.ExecContext(ctx, q, id)
	if err != nil {
		return errors.Wrap(err, "Error deleting inquiry from DB")
	}

	if num, err := res.RowsAffected(); num == 0 || err != nil {
		if num == 0 {
			return sql.ErrNoRows
		}
		return err

	}

	return nil
}
