package services

import (
	"github.com/ucladevx/BPool/interfaces"
	"github.com/ucladevx/BPool/models"
	"github.com/ucladevx/BPool/utils/auth"
)

type (
	// RideService provides all use cases for rides
	RideService struct {
		store      RideStore
		carService *CarService
		logger     interfaces.Logger
	}

	// RideStore any store that allows for rides to be persisted
	RideStore interface {
		GetAll(lastID string, limit int) ([]*models.Ride, error)
		GetByID(id string) (*models.Ride, error)
		Insert(ride *models.Ride) error
		Delete(id string) error
		Update(ride *models.Ride) error
	}
)

// NewRideService creates a new ride service
func NewRideService(store RideStore, c *CarService, l interfaces.Logger) *RideService {
	return &RideService{
		store:      store,
		carService: c,
		logger:     l,
	}
}

// Create persists a ride
func (r *RideService) Create(ride *models.Ride, user *auth.UserClaims) error {
	if err := ride.Validate(); err != nil {
		return err
	}

	isOwner, err := r.carService.IsOwnerOrAdmin(ride.CarID, user)
	if err != nil {
		return err
	}

	if !isOwner {
		return ErrNotCarOwner
	}

	if err := r.store.Insert(ride); err != nil {
		r.logger.Error("RideService.Create - unable to create ride", "error", err.Error())
		return err
	}

	return nil
}

// Update attempts to apply updates to a ride
func (r *RideService) Update(updates *models.RideChangeSet, rideID string, user *auth.UserClaims) (*models.Ride, error) {
	ride, err := r.store.GetByID(rideID)
	if err != nil {
		return nil, err
	}

	if user.AuthLevel != AdminLevel && ride.DriverID != user.ID {
		return nil, ErrForbidden
	}

	originalCarID := ride.CarID

	err = ride.ApplyUpdates(updates)
	if err != nil {
		r.logger.Error("RideService.Update - apply updates", "error", err.Error())
		return nil, err
	}

	// Only check if the car is owned by the user if that field gets updated
	if originalCarID != ride.CarID {
		isOwner, err := r.carService.IsOwnerOrAdmin(ride.CarID, user)
		if err != nil {
			return nil, err
		}

		if !isOwner {
			return nil, ErrNotCarOwner
		}
	}

	if err = r.store.Update(ride); err != nil {
		r.logger.Error("RideService.Update - store update", "error", err.Error())
		return nil, err
	}

	return ride, nil
}

// Get returns a ride by ID
func (r *RideService) Get(id string) (*models.Ride, error) {
	return r.store.GetByID(id)
}

// GetAll returns a page of rides
func (r *RideService) GetAll(lastID string, limit, userAuthLevel int) ([]*models.Ride, error) {
	if userAuthLevel < AdminLevel {
		return nil, ErrNotAllowed
	}

	if limit <= 0 || limit > 100 {
		limit = 15
	}

	return r.store.GetAll(lastID, limit)
}

// Delete removes a ride from the store if the user is allowed to
func (r *RideService) Delete(id string, user *auth.UserClaims) error {
	ride, err := r.store.GetByID(id)
	if err != nil {
		r.logger.Error("RideService.Delete - GetRide", "error", err.Error())
		return err
	}

	if user.AuthLevel != AdminLevel && ride.DriverID != user.ID {
		return ErrForbidden
	}

	return r.store.Delete(ride.ID)
}
