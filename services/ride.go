package services

import (
	"github.com/ucladevx/BPool/interfaces"
	"github.com/ucladevx/BPool/models"
)

type (
	// RideService provides all use cases for rides
	RideService struct {
		store  RideStore
		logger interfaces.Logger
	}

	// RideStore any store that allows for users to be persisted
	RideStore interface {
		GetAll(lastID string, limit int) ([]*models.Ride, error)
		GetByID(id string) (*models.Ride, error)
		Insert(ride *models.Ride) error
	}
)

// NewRideService creates a new user service
func NewRideService(store RideStore, l interfaces.Logger) *RideService {
	return &RideService{
		store:  store,
		logger: l,
	}
}

// Create persists a user
func (r *RideService) Create(ride *models.Ride) error {
	if err := ride.Validate(); err != nil {
		return err
	}

	if err := r.store.Insert(ride); err != nil {
		r.logger.Error("RideService.Create - unable to create ride", "error", err.Error())
		return err
	}

	return nil
}

// Get returns a ride by ID
func (r *RideService) Get(id string) (*models.Ride, error) {
	return r.store.GetByID(id)
}

// GetAll returns a page of users
func (r *RideService) GetAll(lastID string, limit, userAuthLevel int) ([]*models.Ride, error) {
	if userAuthLevel < AdminLevel {
		return nil, ErrNotAllowed
	}

	if limit <= 0 || limit > 100 {
		limit = 15
	}

	return r.store.GetAll(lastID, limit)
}
