package services

import (
	"github.com/ucladevx/BPool/models"
)

type (
	// UserService provides all use cases for users
	UserService struct {
		store  UserStore
		logger Logger
	}

	// UserStore any store that allows for users to be persisted
	UserStore interface {
		GetAll(limit, offset int) ([]*models.User, error)
		GetByID(id int) (*models.User, error)
		GetByEmail(email string) (*models.User, error)
		Insert(user *models.User) error
	}
)

// NewUserService creates a new user
func NewUserService(store UserStore, l Logger) *UserService {
	return &UserService{
		store:  store,
		logger: l,
	}
}
