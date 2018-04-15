package services

import (
	"errors"

	"github.com/ucladevx/BPool/interfaces"
	"github.com/ucladevx/BPool/models"
	"github.com/ucladevx/BPool/stores"
)

const (
	carLimit = 10
)

var (
	// ErrInvalidCarEntry error when user submits invalid car object
	ErrCarValidation = errors.New("car model validation failed")

	// ErrNoCarFound error when no car is in db
	ErrAtCarLimit = errors.New("user currently has maximum number of cars")
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
		GetById(id string) (*models.Car, error)
		GetByWhere(fields []string, queryModifiers []stores.QueryModifier) ([]CarRow, error)
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
	return c.store.GetById(id)
}

// AddCar creates a new car
func (c *CarService) AddCar(body map[interface{}]interface{}, userID string) (*models.Car, error) {
	make, year, color := body["make"].(string), body["year"].(int), body["color"].(string)

	queryModifiers := []stores.QueryModifier{
		stores.QueryMod("user_id", stores.EQ, userID),
	}

	fields := []string{"id"}

	carRows, err := c.store.GetByWhere(fields, queryModifiers)

	if len(carRows) > carLimit {
		return nil, ErrAtCarLimit
	}

	car := &models.Car{
		Make:   make,
		Year:   year,
		Color:  color,
		UserID: userID,
	}

	// model validation
	if errs := car.Validate(); len(errs) > 0 {
		c.logger.Info("Validation", errs)
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
func (c *CarService) DeleteCar(id string) error {
	return c.store.Remove(id)
}
