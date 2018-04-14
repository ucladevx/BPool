package services

import (
	"github.com/ucladevx/BPool/interfaces"
	"github.com/ucladevx/BPool/models"
	"github.com/ucladevx/BPool/stores/postgres"
	"github.com/ucladevx/BPool/utils/auth"
)

const (
	// UserLevel is the auth level associated with standard users
	UserLevel = 0
	// AdminLevel is the auth level associated with admin users
	AdminLevel = 300
)

type (
	// UserService provides all use cases for users
	UserService struct {
		store      UserStore
		authorizer *auth.GoogleAuthorizer
		tokenizer  *auth.Tokenizer
		logger     interfaces.Logger
	}

	// UserStore any store that allows for users to be persisted
	UserStore interface {
		GetAll(lastID string, limit int) ([]*models.User, error)
		GetByID(id string) (*models.User, error)
		GetByEmail(email string) (*models.User, error)
		Insert(user *models.User) error
	}
)

// NewUserService creates a new user service
func NewUserService(store UserStore, t *auth.Tokenizer, l interfaces.Logger) *UserService {
	a := auth.NewGoogleAuthorizer(l)

	return &UserService{
		store:      store,
		logger:     l,
		authorizer: a,
		tokenizer:  t,
	}
}

// Login checks the token for a valid login, stores user info, and generates an auth token
func (u *UserService) Login(token string) (string, error) {
	googleUser, err := u.authorizer.UserLogin(token)

	if err != nil {
		return "", err
	}

	// check if user exists
	user, err := u.store.GetByEmail(googleUser.Email)
	if err != nil && err != postgres.ErrNoUserFound {
		return "", err
	}

	// first time we are seeing user
	if user == nil {
		user = &models.User{
			FirstName:    googleUser.FirstName,
			LastName:     googleUser.LastName,
			Email:        googleUser.Email,
			ProfileImage: googleUser.Picture,
		}

		if err := u.store.Insert(user); err != nil {
			u.logger.Error("UserService.Login - unable to create user", "error", err.Error())
			return "", err
		}
	}

	// generate our own access_token
	claims := map[string]interface{}{
		"id":         user.ID,
		"email":      user.Email,
		"auth_level": user.AuthLevel,
		"sub":        "access",
	}

	authToken, err := u.tokenizer.NewToken(claims)

	if err != nil {
		return "", err
	}

	return authToken, nil
}

// Get returns a user by ID
func (u *UserService) Get(id string) (*models.User, error) {
	return u.store.GetByID(id)
}

// GetAll returns a page of users
func (u *UserService) GetAll(lastID string, limit, userAuthLevel int) ([]*models.User, error) {
	if userAuthLevel < AdminLevel {
		return nil, ErrNotAllowed
	}

	if limit <= 0 || limit > 100 {
		limit = 15
	}

	return u.store.GetAll(lastID, limit)
}
