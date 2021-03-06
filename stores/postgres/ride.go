package postgres

import (
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"

	"github.com/ucladevx/BPool/models"
	"github.com/ucladevx/BPool/stores"
	"github.com/ucladevx/BPool/utils/id"
)

var (
	// ErrNoRideFound error when no ride in db
	ErrNoRideFound = errors.New("no ride found")
)

// RideStore persits rides in a pg DB
type RideStore struct {
	db    *sqlx.DB
	idGen IDgen
}

// NewRideStore creates a new pg ride store
func NewRideStore(db *sqlx.DB) *RideStore {
	return &RideStore{
		db:    db,
		idGen: id.New,
	}
}

// GetAll returns paginated list of rides, lastID = '' for no offset
func (r *RideStore) GetAll(lastID string, limit int) ([]*models.Ride, error) {
	rides := []*models.Ride{}

	if err := r.db.Select(&rides, rideGetAllSQL, lastID, limit); err != nil {
		return nil, err
	}

	return rides, nil
}

// GetAllWherePassenger returns all rides that the user passed is a passenger
func (r *RideStore) GetAllWherePassenger(passengerID string) ([]*models.Ride, error) {
	rides := []*models.Ride{}

	if err := r.db.Select(&rides, rideGetAllWherePassenger, passengerID); err != nil {
		return nil, err
	}

	return rides, nil
}

// WhereMany provides a generic query interface to get many rides
func (r *RideStore) WhereMany(clauses []stores.QueryModifier) ([]*models.Ride, error) {
	where, vals := generateWhereStatement(&clauses)
	query := "SELECT * FROM " + rideTableName + " " + where

	rides := []*models.Ride{}

	if err := r.db.Select(&rides, query, vals...); err != nil {
		return nil, err
	}

	return rides, nil
}

// GetByID finds a ride by ID if exits in the DB
func (r *RideStore) GetByID(id string) (*models.Ride, error) {
	return r.getBy(rideGetByIDSQL, id)
}

// Insert persists a ride to the DB
func (r *RideStore) Insert(ride *models.Ride) error {
	ride.ID = r.idGen()
	row := r.db.QueryRow(
		rideInsertSQL,
		ride.ID,
		ride.DriverID,
		ride.CarID,
		ride.Seats,
		ride.StartCity,
		ride.EndCity,
		ride.StartLat,
		ride.StartLon,
		ride.EndLat,
		ride.EndLon,
		ride.PricePerSeat,
		ride.Info,
		ride.StartDate,
	)

	if err := row.Scan(&ride.CreatedAt, &ride.UpdatedAt); err != nil {
		return err
	}

	return nil
}

// Update persists the updates for the given ride
func (r *RideStore) Update(ride *models.Ride) error {
	row := r.db.QueryRow(
		rideUpdateSQL,
		ride.CarID,
		ride.Seats,
		ride.StartCity,
		ride.EndCity,
		ride.StartLat,
		ride.StartLon,
		ride.EndLat,
		ride.EndLon,
		ride.PricePerSeat,
		ride.Info,
		ride.StartDate,
		ride.ID,
	)

	if err := row.Scan(&ride.UpdatedAt); err != nil {
		return err
	}

	return nil
}

// Delete deletes the ride, does no verification
func (r *RideStore) Delete(id string) error {
	_, err := r.db.Exec(rideDeleteSQL, id)
	return err
}

func (r *RideStore) getBy(query string, arg interface{}) (*models.Ride, error) {
	ride := models.Ride{}

	if err := r.db.Get(&ride, query, arg); err != nil {
		if err == sql.ErrNoRows {
			err = ErrNoRideFound
		}

		return nil, err
	}

	return &ride, nil
}

// Migrate creates ride table in DB
func (r *RideStore) migrate() {
	r.db.MustExec(rideCreateTable)
}
