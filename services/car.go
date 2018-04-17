package services

import (
	"errors"

	"github.com/ucladevx/BPool/stores/postgres"

	"github.com/ucladevx/BPool/interfaces"
	"github.com/ucladevx/BPool/models"
	"github.com/ucladevx/BPool/stores"
)

const (
	carLimit = 10
)

var (
	// ErrCarValidation error when user submits invalid car object
	ErrCarValidation = errors.New("car model validation failed")

	// ErrAtCarLimit error when no car is in db
	ErrAtCarLimit = errors.New("user currently has maximum number of cars")

	// ErrNotCarOwner error when no car is in db
	ErrNotCarOwner = errors.New("user does not own car")
)

type (
	// CarService provides all use cases for users
	CarService struct {
		store  CarStore
		logger interfaces.Logger
	}

	// CarStore any store that allows for users to be persisted
	CarStore interface {
		GetAll(lastID string, limit int) ([]*models.Car, error)
		GetByID(id string) (*models.Car, error)
		GetCount(queryModifiers []stores.QueryModifier) (int, error)
		GetByWhere(fields []string, queryModifiers []stores.QueryModifier) ([]postgres.CarRow, error)
		Insert(user *models.Car) error
		Remove(id string) error
	}
)

// NewCarService creates a new car
func NewCarService(store CarStore, l interfaces.Logger) *CarService {
	return &CarService{
		store:  store,
		logger: l,
	}
}

// GetAllCars returns all cars
func (c *CarService) GetAllCars(lastID string, limit int, authLevel int) ([]*models.Car, error) {
	if authLevel < AdminLevel {
		return nil, ErrNotAllowed
	}

	if limit <= 0 || limit > 100 {
		limit = 15
	}

	return c.store.GetAll(lastID, limit)
}

// GetCar returns a car by id
func (c *CarService) GetCar(id string) (*models.Car, error) {
	return c.store.GetByID(id)
}

// AddCar creates a new car
func (c *CarService) AddCar(body CarRequestBody, userID string) (*models.Car, error) {
	make, model, year, color := body.Make, body.Model, body.Year, body.Color

	queryModifiers := []stores.QueryModifier{
		stores.QueryMod("user_id", stores.EQ, userID),
	}

	count, err := c.store.GetCount(queryModifiers)

	if err != nil {
		return nil, err
	}

	if count >= carLimit {
		return nil, ErrAtCarLimit
	}

	car := &models.Car{
		Make:   make,
		Model:  model,
		Year:   year,
		Color:  color,
		UserID: userID,
	}

	// model validation
	if errs := car.Validate(); len(errs) > 0 {
		c.logger.Info("CarService.AddCar - validate", "error", errs)
		return nil, ErrCarValidation
	}

	// db insertion
	if err := c.store.Insert(car); err != nil {
		c.logger.Error("CarService.AddCar - unable to create car", "error", err.Error())
		return nil, err
	}

	return car, nil
}

// DeleteCar deletes based on id
func (c *CarService) DeleteCar(id, userID string) error {
	car, err := c.store.GetByID(id)

	if err != nil {
		return err
	}

	if car.UserID != userID {
		c.logger.Error("CarService.DeleteCar - unable to delete car", "error", ErrNotCarOwner)
		return ErrNotCarOwner
	}

	return c.store.Remove(id)
}
