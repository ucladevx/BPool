package postgres

import (
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	"github.com/ucladevx/BPool/models"
	"github.com/ucladevx/BPool/utils/id"
)

var (
	// ErrNoUserFound error when no user in db
	ErrNoUserFound = errors.New("no user found")

	// ErrUserAlreadyExists error when a duplicate user insertion is attempted
	ErrUserAlreadyExists = errors.New("user already exists")
)

// UserStore persits users in a pg DB
type UserStore struct {
	db    *sqlx.DB
	idGen IDgen
}

// NewUserStore creates a new pg user store
func NewUserStore(db *sqlx.DB) *UserStore {
	return &UserStore{
		db:    db,
		idGen: id.New,
	}
}

// GetAll returns paginated list of users, lastID = '' for no offset
func (u *UserStore) GetAll(lastID string, limit int) ([]*models.User, error) {
	users := []*models.User{}

	if err := u.db.Select(&users, userGetAllSQL, lastID, limit); err != nil {
		return nil, err
	}

	return users, nil
}

// GetByID finds a user by ID if exits in the DB
func (u *UserStore) GetByID(id string) (*models.User, error) {
	return u.getBy(userGetByIDSQL, id)
}

// GetByEmail finds a user by email if exits in the DB
func (u *UserStore) GetByEmail(email string) (*models.User, error) {
	return u.getBy(userGetByEmailSQL, email)
}

// Insert persists a user to the DB
func (u *UserStore) Insert(user *models.User) error {
	user.ID = u.idGen()
	row := u.db.QueryRow(userInsertSQL, user.ID, user.FirstName, user.LastName, user.Email, user.ProfileImage)

	if err := row.Scan(&user.AuthLevel, &user.CreatedAt, &user.UpdatedAt); err != nil {
		if pgErr, ok := err.(*pq.Error); ok {
			if pgErr.Code == "23505" {
				return ErrUserAlreadyExists
			}
		}
		return err
	}

	return nil
}

func (u *UserStore) getBy(query string, arg interface{}) (*models.User, error) {
	var user models.User

	if err := u.db.Get(&user, query, arg); err != nil {
		if err == sql.ErrNoRows {
			err = ErrNoUserFound
		}

		return nil, err
	}

	return &user, nil
}

// Migrate creates user table in DB
func (u *UserStore) migrate() {
	u.db.MustExec(userCreateTable)
}
