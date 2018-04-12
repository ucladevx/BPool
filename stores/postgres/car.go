package postgres

import (
	// "database/sql"

	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	"github.com/ucladevx/BPool/models"
)

var (
	// ErrInvalidCarEntry error when user submits invalid car object
	ErrInvalidCarEntry = errors.New("Invalid car entry")

	// ErrNoCarFound error when no car is in db
	ErrNoCarFound = errors.New("No car found")
)

//CarStore persists cars in pgsql db
type CarStore struct {
	db *sqlx.DB
}

// NewCarStore creates a new pg car store
func NewCarStore(db *sqlx.DB) *CarStore {
	return &CarStore{
		db: db,
	}
}

func (c *CarStore) GetAll(limit, offset int) ([]*models.Car, error) {
	cars := []*models.Car{}

	if err := c.db.Get(&cars, carsGetAllSQL, offset, limit); err != nil {
		return nil, err
	}

	return cars, nil
}

// Insert persists a user to the DB
func (c *CarStore) Insert(car *models.Car) error {
	row := c.db.QueryRow(userInsertSQL, car.Make, car.Model, car.Year, car.Color, car.UserID, car.ID)

	if err := row.Scan(&car.ID, &car.CreatedAt, &car.UpdatedAt); err != nil {
		if pgErr, ok := err.(*pq.Error); ok {
			if pgErr.Code == "23505" {
				return ErrUserAlreadyExists
			}
		}
		return err
	}

	return nil
}

func (c *CarStore) Delete(car *models.Car) error {
	_, err := c.db.Exec(carsDeleteSQL, car.ID)

	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok {
			return ErrInvalidCarEntry
		}

		return err
	}

	return nil
}

// func (c *CarStore) GetByID(id int) (*models.Car, error) {
// 	var fields [1]string
// 	fields[0] = strconv.Itoa(id)
// 	return c.getByFields(carsGetByFieldsSQL, fields)
// }

// func (c *CarStore) getBy

// TODO: Figure out how to do migrations
