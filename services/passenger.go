package services

import (
	"errors"

	"github.com/ucladevx/BPool/interfaces"
	"github.com/ucladevx/BPool/models"
	"github.com/ucladevx/BPool/stores"
	"github.com/ucladevx/BPool/utils/auth"
)

var (
	// ErrNoMoreSeats occurs when a ride is already full
	ErrNoMoreSeats = errors.New("There are no more seats available")

	// ErrPassengerIsDriver is returned when a driver tries to be a passenger in own ride
	ErrPassengerIsDriver = errors.New("You cannot join your own ride")
)

type (
	// PassengerService provides all use cases for ride passengers
	PassengerService struct {
		store       PassengerStore
		rideService *RideService
		logger      interfaces.Logger
	}

	// PassengerStore any store that allows for ride passengers to be persisted
	PassengerStore interface {
		GetAll(lastID string, limit int) ([]*models.Passenger, error)
		GetByID(id string) (*models.Passenger, error)
		Insert(ride *models.Passenger) error
		Delete(id string) error
		Update(ride *models.Passenger) error
		Count(clauses []stores.QueryModifier) (int, error)
		WhereMany(clauses []stores.QueryModifier) ([]*models.Passenger, error)
	}
)

// NewPassengerService creates a new ride service
func NewPassengerService(store PassengerStore, r *RideService, l interfaces.Logger) *PassengerService {
	return &PassengerService{
		store:       store,
		rideService: r,
		logger:      l,
	}
}

// Create persists a ride passenger
func (p *PassengerService) Create(passenger *models.Passenger, user *auth.UserClaims) error {
	passenger.Status = models.PassengerInterested

	if err := passenger.Validate(); err != nil {
		return err
	}

	// only a passenger can initiate becoming a passenger
	if passenger.PassengerID != user.ID && user.AuthLevel != AdminLevel {
		return ErrForbidden
	}

	// check if ride exists
	ride, err := p.rideService.Get(passenger.RideID)
	if err != nil {
		return err
	}

	passenger.DriverID = ride.DriverID

	if passenger.DriverID == user.ID {
		return ErrPassengerIsDriver
	}

	if err := p.store.Insert(passenger); err != nil {
		p.logger.Error("PassengerService.Create - unable to create passenger", "error", err.Error())
		return err
	}

	return nil
}

// Update attempts to apply updates to a ride passenger
func (p *PassengerService) Update(updates *models.PassengerChangeSet, passengerID string, user *auth.UserClaims) (*models.Passenger, error) {
	passenger, err := p.store.GetByID(passengerID)
	if err != nil {
		return nil, err
	}

	// only a driver can update a status
	if user.AuthLevel != AdminLevel && passenger.DriverID != user.ID {
		return nil, ErrForbidden
	}

	err = passenger.ApplyUpdates(updates)
	if err != nil {
		p.logger.Error("PassengerService.Update - apply updates", "error", err.Error())
		return nil, err
	}

	if passenger.Status == models.PassengerAccepted {
		clauses := []stores.QueryModifier{
			stores.QueryMod("ride_id", stores.EQ, passenger.RideID),
			stores.And,
			stores.QueryMod("status", stores.EQ, models.PassengerAccepted),
		}

		// NOTE: probably try to get count and ride in parallel
		alreadyAccepted, err := p.store.Count(clauses)
		if err != nil {
			return nil, err
		}

		ride, err := p.rideService.Get(passenger.RideID)
		if err != nil {
			return nil, err
		}

		if alreadyAccepted >= ride.Seats {
			return nil, ErrNoMoreSeats
		}
	}

	if err = p.store.Update(passenger); err != nil {
		p.logger.Error("PassengerService.Update - store update", "error", err.Error())
		return nil, err
	}

	return passenger, nil
}

// Get returns a ride by ID
func (p *PassengerService) Get(id string, user *auth.UserClaims) (*models.Passenger, error) {
	passenger, err := p.store.GetByID(id)
	if err != nil {
		return nil, err
	}

	// only allow driver, passenger, or admin to view the ride
	if passenger.DriverID != user.ID && passenger.PassengerID != user.ID && user.AuthLevel != AdminLevel {
		return nil, ErrForbidden
	}

	return passenger, nil
}

// GetAll returns a page of rides
func (p *PassengerService) GetAll(lastID string, limit, userAuthLevel int) ([]*models.Passenger, error) {
	if userAuthLevel < AdminLevel {
		return nil, ErrNotAllowed
	}

	if limit <= 0 || limit > 100 {
		limit = 15
	}

	return p.store.GetAll(lastID, limit)
}

// GetAllByCarID returns all passengers in all statuses for a given car
func (p *PassengerService) GetAllByCarID(carID string) ([]*models.Passenger, error) {
	return p.store.WhereMany([]stores.QueryModifier{stores.QueryMod("car_id", stores.EQ, carID)})
}

// Delete removes a ride from the store if the user is allowed to
func (p *PassengerService) Delete(id string, user *auth.UserClaims) error {
	passenger, err := p.store.GetByID(id)
	if err != nil {
		p.logger.Error("PassengerService.Delete", "error", err.Error())
		return err
	}

	if user.AuthLevel != AdminLevel && passenger.DriverID != user.ID {
		return ErrForbidden
	}

	return p.store.Delete(passenger.ID)
}
