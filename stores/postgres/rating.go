package postgres

import (
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/ucladevx/BPool/models"
	"github.com/ucladevx/BPool/stores"
	"github.com/ucladevx/BPool/utils/id"
)

// RatingStore persists ratings in DB
type RatingStore struct {
	db    *sqlx.DB
	idGen IDgen
}

var (
	// ErrNoRatingFound error when no rating in db
	ErrNoRatingFound = errors.New("no rating found")
)

// NewRatingStore creates a new pg rating store
func NewRatingStore(db *sqlx.DB) *RatingStore {
	return &RatingStore{
		db:    db,
		idGen: id.New,
	}
}

// Insert persists a new rating into the DB
func (ratingStore *RatingStore) Insert(rating *models.Rating) error {
	rating.ID = ratingStore.idGen()

	row := ratingStore.db.QueryRow(
		ratingInsertSQL,
		rating.ID,
		rating.Rating,
		rating.RideID,
		rating.RaterID,
		rating.RateeID,
		rating.Comment,
	)

	return row.Scan(&rating.CreatedAt, &rating.UpdatedAt)
}

// GetByID finds a rating by ID if exits in the DB
func (ratingStore *RatingStore) GetByID(id string) (*models.Rating, error) {
	return ratingStore.getBy(ratingGetByIDSQL, id)
}

// Average determines average rating for provided rideID
func (ratingStore *RatingStore) Average(userID string) (float32, error) {
	row := ratingStore.db.QueryRow(ratingGetAverageRatingSQL, userID)
	average := float32(0.0)

	if err := row.Scan(&row); err != nil {
		return 0.0, err
	}

	return average, nil
}

// Delete deleted rating
func (ratingStore *RatingStore) Delete(ratingID string) error {
	_, err := ratingStore.db.Exec(ratingDeleteSQL, ratingID)
	return err
}

func (ratingStore *RatingStore) getBy(query string, arg interface{}) (*models.Rating, error) {
	rating := models.Rating{}

	if err := ratingStore.db.Get(&rating, query, arg); err != nil {
		if err == sql.ErrNoRows {
			err = ErrNoRatingFound
		}

		return nil, err
	}

	return &rating, nil
}

// WhereMany provides a generic query interface to get many ratings
func (ratingStore *RatingStore) WhereMany(clauses []stores.QueryModifier) ([]*models.Rating, error) {
	where, vals := generateWhereStatement(&clauses)
	query := "SELECT * FROM ratings " + where

	ratings := []*models.Rating{}

	if err := ratingStore.db.Select(&ratings, query, vals...); err != nil {
		return nil, err
	}

	return ratings, nil
}

// Migrate create ratings table in DB
func (ratingStore *RatingStore) migrate() {
	ratingStore.db.MustExec(ratingCreateTable)
}
