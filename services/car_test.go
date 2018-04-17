package services_test

import (
	"testing"

	"github.com/ucladevx/BPool/stores"
	"github.com/ucladevx/BPool/stores/postgres"

	"github.com/stretchr/testify/assert"
	"github.com/ucladevx/BPool/mocks"
	"github.com/ucladevx/BPool/models"
	"github.com/ucladevx/BPool/services"
)

var (
	testCar = models.Car{
		ID:     "testCar1",
		Make:   "Toyata",
		Model:  "Prius",
		Year:   2015,
		Color:  "White",
		UserID: "user1",
	}

	testCarBody = models.Car{
		Make:   "Toyata",
		Model:  "Prius",
		Year:   2015,
		Color:  "White",
		UserID: "user1",
	}

	testUser = models.User{
		ID:           "user1",
		FirstName:    "JOHN",
		LastName:     "DOE",
		Email:        "johndoe@g.ucla.edu",
		ProfileImage: "ucladevx.com",
		AuthLevel:    services.UserLevel,
	}
)

func newCarService(store *mocks.CarStore) *services.CarService {
	logger := mocks.Logger{}

	return services.NewCarService(store, logger)
}

func TestCarAdd(t *testing.T) {
	store := new(mocks.CarStore)
	service := newCarService(store)
	assert := assert.New(t)

	queryModifiers := []stores.QueryModifier{
		stores.QueryMod("user_id", stores.EQ, testUser.ID),
	}

	// testing model validation
	store.On("GetCount", queryModifiers).Return(0, nil)

	invalidRequestBody := services.CarRequestBody{
		Make:  "",
		Year:  3000,
		Color: "",
	}

	noCar, err := service.AddCar(invalidRequestBody, testUser.ID)

	assert.Nil(noCar, "user cannot insert a car with invalid parameters")
	assert.EqualError(err, "car model validation failed")

	validRequestBody := services.CarRequestBody{
		Make:  "Toyata",
		Model: "Prius",
		Year:  2015,
		Color: "White",
	}

	store.On("Insert", &testCarBody).Return(nil)

	car, err := service.AddCar(validRequestBody, testUser.ID)

	assert.Nil(err, "no error is returned when passing in valid request body")
	assert.Equal(car, &testCarBody)

	store.AssertExpectations(t)
}

func TestCarDelete(t *testing.T) {
	store := new(mocks.CarStore)
	service := newCarService(store)
	assert := assert.New(t)

	// testing unauthorized deletion
	store.On("GetByID", testCar.ID).Return(&testCar, nil)

	err := service.DeleteCar(testCar.ID, "notUserID")

	assert.EqualError(err, "user does not own car")

	store.On("Remove", testCar.ID).Return(nil)

	err = service.DeleteCar(testCar.ID, testUser.ID)

	assert.Nil(err, "no error is returned when passing in valid request body")

	store.AssertExpectations(t)
}

func TestCarGetByID(t *testing.T) {
	store := new(mocks.CarStore)
	service := newCarService(store)
	assert := assert.New(t)

	store.On("GetByID", "invalidID").Return(nil, postgres.ErrNoCarFound)

	noCar, err := service.GetCar("invalidID")

	assert.Nil(noCar, "querying with an invalid ID should return no car")
	assert.EqualError(err, "no car found")

	c := testCar
	validID := "testCar1"

	store.On("GetByID", validID).Return(&c, nil)

	car, err := service.GetCar(validID)
	assert.Nil(err, "for valid car id, err should be nil")
	assert.Equal(car.ID, validID)

	store.AssertExpectations(t)
}

func TestCarGetAll(t *testing.T) {
	store := new(mocks.CarStore)
	service := newCarService(store)
	assert := assert.New(t)

	noCars, err := service.GetAllCars("", 10, services.UserLevel)
	assert.Nil(noCars, "returns no cars when user does not have correct auth level")
	assert.EqualError(err, "user is not allowed")

	store.On("GetAll", "", 15).Return([]*models.Car{&testCar}, nil)
	badLimit := -1

	cars, err := service.GetAllCars("", badLimit, services.AdminLevel)

	assert.Nil(err, "no error should be returned for a bad limit")
	assert.Equal(1, len(cars), "returned cars should have length 1")

	store.AssertExpectations(t)
}
