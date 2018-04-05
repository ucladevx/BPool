package services

import (
	"github.com/ucladevx/BPool/models"
)

type (
	// CarService provides all use cases for users
	CarService struct {
		store  CarStore
		logger Logger
	}

	// CarStore any store that allows for users to be persisted
	CarStore interface {
		GetAll(limit, offset int) ([]*models.Car, error)
		GetByID(id int) (*models.Car, error)
		Search(str string) (*models.Car, error)
		Insert(user *models.Car) error
		Remove(id int) error
	}
)

// NewCarService creates a new car
func NewCarService(store CarStore, l Logger) *CarService {
	return &CarService{
		store:  store,
		logger: l,
	}
}
