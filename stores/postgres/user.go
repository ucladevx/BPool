package postgres

import (
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	"github.com/ucladevx/BPool/models"
)

var (
	// ErrNoUserFound error when no user in db
	ErrNoUserFound = errors.New("no user found")

	// ErrUserAlreadyExists error when a duplicate user insertion is attempted
	ErrUserAlreadyExists = errors.New("user already exists")
)

// UserStore persits users in a pg DB
type UserStore struct {
	db *sqlx.DB
}

// NewUserStore creates a new pg user store
func NewUserStore(db *sqlx.DB) *UserStore {
	return &UserStore{
		db: db,
	}
}

// GetAll returns paginated list of users
func (u *UserStore) GetAll(limit, offset int) ([]*models.User, error) {
	users := []*models.User{}

	if err := u.db.Get(&users, userGetAllSQL, offset, limit); err != nil {
		return nil, err
	}

	return users, nil
}

// GetByID finds a user by ID if exits in the DB
func (u *UserStore) GetByID(id int) (*models.User, error) {
	return u.getBy(userGetByIDSQL, id)
}

// GetByEmail finds a user by email if exits in the DB
func (u *UserStore) GetByEmail(email string) (*models.User, error) {
	return u.getBy(userGetByEmailSQL, email)
}

// Insert persists a user to the DB
func (u *UserStore) Insert(user *models.User) error {
	row := u.db.QueryRow(userInsertSQL, user.Name, user.Email, user.ID)

	if err := row.Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt); err != nil {
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
