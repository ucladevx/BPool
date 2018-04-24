package postgres

import (
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"

	"github.com/ucladevx/BPool/models"
	"github.com/ucladevx/BPool/utils/id"
)

var (
	// ErrNoPassengerFound error when no passenger in db
	ErrNoPassengerFound = errors.New("no passenger found")
)

// PassengerStore persits rides in a pg DB
type PassengerStore struct {
	db    *sqlx.DB
	idGen IDgen
}

// NewPassengerStore creates a new pg passenger store
func NewPassengerStore(db *sqlx.DB) *PassengerStore {
	return &PassengerStore{
		db:    db,
		idGen: id.New,
	}
}

// GetAll returns paginated list of passengers, lastID = '' for no offset
func (r *PassengerStore) GetAll(lastID string, limit int) ([]*models.Passenger, error) {
	passengers := []*models.Passenger{}

	if err := r.db.Select(&passengers, passengerGetAllSQL, lastID, limit); err != nil {
		return nil, err
	}

	return passengers, nil
}

// GetByID finds a passenger by ID if exits in the DB
func (r *PassengerStore) GetByID(id string) (*models.Passenger, error) {
	return r.getBy(passengerGetByIDSQL, id)
}

// Insert persists a passenger to the DB
func (r *PassengerStore) Insert(passenger *models.Passenger) error {
	passenger.ID = r.idGen()
	row := r.db.QueryRow(
		passengerInsertSQL,
		passenger.ID,
		passenger.DriverID,
		passenger.PassengerID,
		passenger.RideID,
		passenger.PassengerID,
	)

	if err := row.Scan(&passenger.CreatedAt, &passenger.UpdatedAt); err != nil {
		return err
	}

	return nil
}

// Update persists the updates for the given ride
func (r *PassengerStore) Update(passenger *models.Passenger) error {
	row := r.db.QueryRow(
		passengerUpdateSQL,
		passenger.Status,
		passenger.ID,
	)

	if err := row.Scan(&passenger.UpdatedAt); err != nil {
		return err
	}

	return nil
}

// Delete deletes the passenger, does no verification
func (r *PassengerStore) Delete(id string) error {
	_, err := r.db.Exec(passengerDeleteSQL, id)
	return err
}

func (r *PassengerStore) getBy(query string, arg interface{}) (*models.Passenger, error) {
	passenger := models.Passenger{}

	if err := r.db.Get(&passenger, query, arg); err != nil {
		if err == sql.ErrNoRows {
			err = ErrNoPassengerFound
		}

		return nil, err
	}

	return &passenger, nil
}

// Migrate creates passenger table in DB
func (r *PassengerStore) migrate() {
	r.db.MustExec(passengerCreateTable)
}
