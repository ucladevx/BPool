package postgres

import (
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"

	"github.com/ucladevx/BPool/models"
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

// GetByID finds a ride by ID if exits in the DB
func (r *RideStore) GetByID(id string) (*models.Ride, error) {
	return r.getBy(rideGetByIDSQL, id)
}

// Insert persists a ride to the DB
func (r *RideStore) Insert(ride *models.Ride) error {
	ride.ID = r.idGen()
	row := r.db.QueryRow(
		rideInsertSQL,
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
	)

	if err := row.Scan(&ride.CreatedAt, &ride.UpdatedAt); err != nil {
		return err
	}

	return nil
}

func (r *RideStore) getBy(query string, arg interface{}) (*models.Ride, error) {
	var ride models.Ride

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
