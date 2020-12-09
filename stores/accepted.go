package stores

import (
	"context"
	"database/sql"

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
	GetAll(ctx context.Context) (models.AcceptedList, error)
	ProcessInquiry(ctx context.Context, accepted *models.Accepted) (int64, error)
	Delete(ctx context.Context, id int64) error
}

type acceptedStoreSql struct {
	dbFactory db.DbFactory
}

func (a *acceptedStoreSql) GetAll(ctx context.Context) (models.AcceptedList, error) {
	db := a.dbFactory.Connect()
	defer db.Close()

	// TODO: add index to date_accepted
	q := `SELECT id, inquirer, inquirer_email, inquirer_phone,
				item_id, item_title, item_price, date_reservation 
			FROM accepted a
			ORDER BY a.date_accepted DESC`

	rows, err := db.QueryContext(ctx, q)
	if err != nil {
		return nil, errors.Wrap(err, "Error querying all accepted from db")
	}

	var acceptedList []*models.Accepted
	for rows.Next() {
		accepted := &models.Accepted{}
		if err := rows.Scan(&accepted.Id, &accepted.Inquirer, &accepted.InquirerEmail,
			&accepted.InquirerPhone, &accepted.ItemId, &accepted.ItemTitle, &accepted.ItemPrice, &accepted.DateReservation); err != nil {
			return nil, errors.Wrap(err, "Error scaning accepted info to model")
		}
		acceptedList = append(acceptedList, accepted)
	}
	return acceptedList, nil
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

func (a *acceptedStoreSql) Delete(ctx context.Context, id int64) error {
	db := a.dbFactory.Connect()
	defer db.Close()

	tx, err := db.BeginTx(ctx, nil)
	defer tx.Rollback()

	if err != nil {
		return errors.Wrap(err, "Error initializing transaction for Delete in accepted store")
	}

	q := "DELETE FROM accepted WHERE id = $1"
	res, err := tx.ExecContext(ctx, q, id)
	if err != nil {
		return errors.Wrapf(err, "Error deleting accepted inside store. Id: %v. Error: %v", id, err)
	}

	if num, err := res.RowsAffected(); err != nil || num == 0 {
		if num == 0 {
			return sql.ErrNoRows
		}
		return errors.Wrap(err, "Error retrieving rows affected")
	}

	tx.Commit()
	return nil
}
