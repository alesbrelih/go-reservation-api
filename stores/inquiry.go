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
	Create(ctx context.Context, inquiry *models.InquiryCreate) error
}

type inquiryStoreSql struct {
	dbFactory db.DbFactory
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
