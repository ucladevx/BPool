package services_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/ucladevx/BPool/mocks"
	"github.com/ucladevx/BPool/models"
	"github.com/ucladevx/BPool/services"
	"github.com/ucladevx/BPool/stores/postgres"
	"github.com/ucladevx/BPool/utils/auth"
)

var (
	ride1 = models.Ride{
		ID:           "abc",
		DriverID:     "123",
		CarID:        "xyz",
		Seats:        3,
		StartCity:    "Los Angeles",
		EndCity:      "San Francisco",
		StartLat:     11.00,
		StartLon:     12.00,
		EndLat:       -1.00,
		EndLon:       -23.00,
		PricePerSeat: 15,
		Info:         "test",
		StartDate:    time.Now().Add(time.Hour),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	ride2 = models.Ride{
		ID:           "lit",
		DriverID:     "456",
		CarID:        "iop",
		Seats:        0,
		StartCity:    "Las Vegas",
		EndCity:      "New York",
		StartLat:     1.00,
		StartLon:     12.00,
		EndLat:       -1.00,
		EndLon:       -23.00,
		PricePerSeat: 15,
		Info:         "",
		StartDate:    time.Now(),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
)

func newRideService(store *mocks.RideStore) *services.RideService {
	logger := mocks.Logger{}
	return services.NewRideService(store, logger)
}

func TestRideGet(t *testing.T) {
	store := new(mocks.RideStore)
	service := newRideService(store)
	assert := assert.New(t)

	store.On("GetByID", "abc").Return(nil, postgres.ErrNoRideFound)

	noRide, err := service.Get("abc")

	assert.Nil(noRide, "for a bad id there should be no ride")
	assert.Equal(postgres.ErrNoRideFound, err, "if no ride found should return no ride found error")

	r := ride1
	validID := "test1234"
	r.ID = validID

	store.On("GetByID", validID).Return(&r, nil)

	ride, err := service.Get(validID)
	assert.NotNil(ride, "ride should not be nil for a valid id")
	assert.Nil(err, "for a valid ride, error should be nil")
}

func TestRideCreate(t *testing.T) {
	store := new(mocks.RideStore)
	service := newRideService(store)
	assert := assert.New(t)

	badRide := ride1
	badRide.Seats = -4

	validationErr := service.Create(&badRide)
	assert.NotNil(validationErr, "There should be a validation error when an invalid car is inserted")

	validRide := ride1

	store.On("Insert", mock.AnythingOfType("*models.Ride")).Return(nil)

	noErr := service.Create(&validRide)
	assert.Nil(noErr)
}

func TestRideUpdate(t *testing.T) {
	store := new(mocks.RideStore)
	service := newRideService(store)
	assert := assert.New(t)

	validRide := ride1
	newSeats := -1
	newPrice := 100.0
	badRideChanges := models.RideChangeSet{
		Seats:        &newSeats,
		PricePerSeat: &newPrice,
	}

	notDriver := auth.UserClaims{
		ID:        "ghy",
		Email:     "test@gmail.com",
		AuthLevel: services.UserLevel,
	}

	store.On("GetByID", validRide.ID).Return(&validRide, nil)

	noRide, forbiddenErr := service.Update(&badRideChanges, validRide.ID, &notDriver)

	assert.Nil(noRide, "there should be no ride when user is not admin or driver")
	assert.NotNil(forbiddenErr, "should have a forbidden error when not admin or driver")
	assert.Equal(services.ErrForbidden, forbiddenErr, "err should be a forbidden when not driver or admin")

	driver := notDriver
	driver.ID = validRide.DriverID

	noRide, validationErr := service.Update(&badRideChanges, validRide.ID, &driver)
	assert.Nil(noRide, "there should be no ride when invalid updates")
	assert.NotNil(validationErr, "should have errored when updates are invalid")

	store.On("GetByID", validRide.ID).Return(&validRide, nil)
	store.On("Update", mock.AnythingOfType("*models.Ride")).Return(nil)

	rideChanges := badRideChanges
	newSeats = 0
	rideChanges.Seats = &newSeats
	updatedRide, noErr := service.Update(&rideChanges, validRide.ID, &driver)

	assert.Nil(noErr, "there should be no error on valid ride update")
	assert.NotNil(updatedRide, "there should have been an updated ride on valid update")
	assert.Equal(newSeats, validRide.Seats, "the new ride should have the updated seats")
	assert.Equal(newPrice, validRide.PricePerSeat, "the new ride should have the updated price per seat")
}

func TestRideGetAll(t *testing.T) {
	store := new(mocks.RideStore)
	service := newRideService(store)
	assert := assert.New(t)

	noRides, err := service.GetAll("", 15, services.UserLevel)
	assert.Nil(noRides, "when a user does not have the right auth level there should be no rides")
	assert.Equal(services.ErrNotAllowed, err, "when a user does not have the right auth level there should be a not allowed error")

	badLimit := -1
	store.On("GetAll", "", 15).Return([]*models.Ride{&ride1, &ride2}, nil)

	rides, err := service.GetAll("", badLimit, services.AdminLevel)

	assert.Nil(err, "for a bad limit, should still return no error")
	assert.Equal(2, len(rides), "the returned rides should have length 2")

	store.AssertExpectations(t)
}

func TestRideDelete(t *testing.T) {
	store := new(mocks.RideStore)
	service := newRideService(store)
	assert := assert.New(t)

	user := auth.UserClaims{
		ID:        "456",
		Email:     "idontcare@gmail.com",
		AuthLevel: 0,
	}

	store.On("GetByID", "abc").Return(&ride1, nil)
	store.On("Delete", "lit").Return(nil)
	err := service.Delete("abc", &user)
	assert.NotNil(err, "there should be an error when the user does not own the ride")

	store.On("GetByID", "lit").Return(&ride2, nil)
	store.On("Delete", "lit").Return(nil)
	err = service.Delete("lit", &user)

	assert.Nil(err, "should have successfully delete ride if user is owner")
}
