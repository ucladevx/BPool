package postgres

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	"github.com/ucladevx/BPool/models"
	"github.com/ucladevx/BPool/stores"
)

var (
	// ErrInvalidCarEntry error when user submits invalid car object
	ErrInvalidCarEntry = errors.New("invalid car entry")

	// ErrNoCarFound error when no car is in db
	ErrNoCarFound = errors.New("no car found")
)

type (
	//CarStore persists cars in pgsql db
	CarStore struct {
		db *sqlx.DB
	}

	// CarRow is desired car database data
	CarRow struct {
		data map[string]interface{}
	}
)

// NewCarStore creates a new pg car store
func NewCarStore(db *sqlx.DB) *CarStore {
	return &CarStore{
		db: db,
	}
}

// GetAll finds all cars from db
func (c *CarStore) GetAll(lastID string, limit int) ([]*models.Car, error) {
	cars := []*models.Car{}

	if err := c.db.Get(&cars, carsGetAllSQL, lastID, limit); err != nil {
		return nil, err
	}

	return cars, nil
}

// GetByID finds car by id if it exists in db
func (c *CarStore) GetByID(id string) (*models.Car, error) {
	return c.getBy(carsGetByIDSQL, id)
}

// GetByWhere queries based on generated WHERE statement
func (c *CarStore) GetByWhere(fields []string, queryModifiers []stores.QueryModifier) ([]CarRow, error) {
	queryString := "SELECT "
	queryString += strings.Join(fields, ", ")

	query, vals := generateWhereStatement(&queryModifiers)

	queryString += query

	rows, err := c.db.Query(queryString, vals)

	defer rows.Close()

	if err != nil {
		return nil, err
	}

	columnNames, err := rows.Columns()

	if err != nil {
		return nil, err
	}
	// load car data with fields

	cars := []CarRow{}
	numFields := len(fields)

	for rows.Next() {
		cr := CarRow{}
		columns := make([]interface{}, numFields)
		columnPointers := make([]interface{}, numFields)

		for i := 0; i < numFields; i++ {
			columnPointers[i] = &columns[i]
		}

		if err := rows.Scan(columnPointers...); err != nil {
			return nil, err
		}

		for i, col := range columns {
			cr.data[columnNames[i]] = col
		}

		cars = append(cars, cr)
	}

	return cars, nil
}

// Insert persists a user to the DB
func (c *CarStore) Insert(car *models.Car) error {
	row := c.db.QueryRow(userInsertSQL, car.Make, car.Model, car.Year, car.Color, car.UserID, car.ID)

	if err := row.Scan(&car.ID, &car.CreatedAt, &car.UpdatedAt); err != nil {
		if pgErr, ok := err.(*pq.Error); ok {
			if pgErr.Code.Name() == "unique_violation" {
				return ErrUserAlreadyExists
			}
		}
		return err
	}

	return nil
}

// Remove car from DB
func (c *CarStore) Remove(id string) error {
	_, err := c.db.Exec(carsDeleteSQL, id)

	if err != nil {
		return ErrInvalidCarEntry
	}

	return nil
}

func (c *CarStore) getBy(query string, args interface{}) (*models.Car, error) {
	var car models.Car

	if err := c.db.Get(&car, query, args); err != nil {
		if err == sql.ErrNoRows {
			err = ErrNoCarFound
		}

		return nil, err
	}

	return &car, nil
}

// TODO: Figure out how to do migrations
