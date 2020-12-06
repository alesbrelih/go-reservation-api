package stores

import (
	"context"

	"github.com/alesbrelih/go-reservation-api/db"
	"github.com/alesbrelih/go-reservation-api/models"
	"github.com/pkg/errors"
)

func NewAcceptedStoreSql(dbFactory db.DbFactory) AcceptedStore {
	return &acceptedStoreSql{
		dbFactory: dbFactory,
	}
}

type AcceptedStore interface {
	ProcessInquiry(ctx context.Context, accepted *models.Accepted) (int64, error)
}

type acceptedStoreSql struct {
	dbFactory db.DbFactory
}

func (a *acceptedStoreSql) ProcessInquiry(ctx context.Context, accepted *models.Accepted) (int64, error) {
	db := a.dbFactory.Connect()
	defer db.Close()

	q := `INSERT INTO accepted 
				(inquirer, inquirer_email, inquirer_phone, 
					inquirer_comment, item_id, item_title, item_price,
					notes, date_reservation, date_inquiry_created,
					date_accepted)
			VALUES 
				($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, now() at time zone 'utc')
			RETURNING id`

	var id int64
	err := db.QueryRowContext(ctx, q, accepted.Inquirer, accepted.InquirerEmail, accepted.InquirerPhone,
		accepted.InquirerComment, accepted.ItemId, accepted.ItemTitle, accepted.ItemPrice,
		accepted.Notes, accepted.DateReservation, accepted.DateInquiryCreated).Scan(&id)

	if err != nil {
		return 0, errors.Wrap(err, "Error processing inquiry to accepted inside DB")
	}

	return id, nil
}
