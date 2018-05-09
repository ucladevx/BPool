package postgres

import (
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	"github.com/ucladevx/BPool/models"
	"github.com/ucladevx/BPool/stores"
	"github.com/ucladevx/BPool/utils/id"
)

var (
	// ErrNoPassengerFound error when no passenger in db
	ErrNoPassengerFound = errors.New("no passenger found")

	// ErrAlreadyPassenger occurs when a user has already shown interest
	ErrAlreadyPassenger = errors.New("user is already a passenger")
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
		passenger.Status,
	)

	if err := row.Scan(&passenger.CreatedAt, &passenger.UpdatedAt); err != nil {
		if pgErr, ok := err.(*pq.Error); ok {
			if pgErr.Code.Name() == "unique_violation" {
				return ErrAlreadyPassenger
			}
		}
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

// Count determines the number of records in the db that fit the where clauses
func (r *PassengerStore) Count(clauses []stores.QueryModifier) (int, error) {
	where, vals := generateWhereStatement(&clauses)

	query := "SELECT COUNT(*) FROM" + passengerTableName + " " + where

	row := r.db.QueryRow(query, vals...)
	count := 0

	if err := row.Scan(&row); err != nil {
		return 0, err
	}

	return count, nil
}

// Where provides a generic query interface to get a single passenger
func (r *PassengerStore) Where(clauses []stores.QueryModifier) (*models.Passenger, error) {
	where, vals := generateWhereStatement(&clauses)
	query := "SELECT * FROM" + passengerTableName + " " + where + " LIMIT 1"

	return r.getBy(query, vals...)
}

// WhereMany provides a generic query interface to get many passengers
func (r *PassengerStore) WhereMany(clauses []stores.QueryModifier) ([]*models.Passenger, error) {
	where, vals := generateWhereStatement(&clauses)
	query := "SELECT * FROM" + passengerTableName + " " + where

	passengers := []*models.Passenger{}

	if err := r.db.Select(&passengers, query, vals...); err != nil {
		return nil, err
	}

	return passengers, nil
}

func (r *PassengerStore) getBy(query string, arg ...interface{}) (*models.Passenger, error) {
	passenger := models.Passenger{}

	if err := r.db.Get(&passenger, query, arg...); err != nil {
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
