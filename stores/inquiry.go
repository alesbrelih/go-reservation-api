package stores

import (
	"context"

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

	q := `INSERT INTO inquiry 
		(inquirer,email,phone,item_id,date_reservation,date_created)
		VALUES
		($1, $2, $3, $4, $5, NOW())`

	_, err := db.ExecContext(ctx, q, inquiry.Inquirer, inquiry.Email, inquiry.Phone, inquiry.ItemId, inquiry.Date)
	if err != nil {
		return errors.Wrap(err, "Error creating new inquiry")
	}

	return nil
}
