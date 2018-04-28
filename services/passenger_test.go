package services_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/ucladevx/BPool/mocks"
	"github.com/ucladevx/BPool/models"
	"github.com/ucladevx/BPool/services"
	"github.com/ucladevx/BPool/utils/auth"
)

var (
	validPassenger = models.Passenger{
		ID:          "lit",
		DriverID:    "456",
		PassengerID: "xyz",
		RideID:      "abc",
		Status:      models.PassengerInterested,
	}
)

type mockedPassengerService struct {
	passengerStore   *mocks.PassengerStore
	rideStore        *mocks.RideStore
	carStore         *mocks.CarStore
	carService       *services.CarService
	rideService      *services.RideService
	passengerService *services.PassengerService
}

func newMockedPassengerService() *mockedPassengerService {
	logger := mocks.Logger{}
	passengerStore := new(mocks.PassengerStore)
	rideStore := new(mocks.RideStore)
	carStore := new(mocks.CarStore)
	carService := newCarService(carStore)
	rideService := newRideService(rideStore, carService)
	passengerService := services.NewPassengerService(passengerStore, rideService, logger)

	return &mockedPassengerService{
		passengerStore:   passengerStore,
		rideStore:        rideStore,
		carStore:         carStore,
		carService:       carService,
		rideService:      rideService,
		passengerService: passengerService,
	}
}

func TestCreatePassenger(t *testing.T) {
	service := newMockedPassengerService()
	assert := assert.New(t)

	user := auth.UserClaims{ID: "xyz"}
	service.rideStore.On("GetByID", "abc").Return(&ride1, nil)
	service.passengerStore.On("Insert", mock.AnythingOfType("*models.Passenger")).Return(nil)

	validPass := validPassenger
	noErr := service.passengerService.Create(&validPass, &user)
	assert.Nil(noErr)
	service.passengerStore.AssertExpectations(t)

	invalidPassenger := auth.UserClaims{ID: "tty"}
	forbiddenErr := service.passengerService.Create(&validPass, &invalidPassenger)
	assert.Equal(services.ErrForbidden, forbiddenErr, "should not insert a passenger that is not the current user")

	driverPassenger := auth.UserClaims{ID: "123"}
	service.rideStore.On("GetByID", "abc").Return(&ride1, nil)
	newPass := validPass
	newPass.PassengerID = "123"
	passengerErr := service.passengerService.Create(&newPass, &driverPassenger)
	assert.Equal(services.ErrPassengerIsDriver, passengerErr, "driver should not be able to be a passenger")
}

func TestUpdatePassenger(t *testing.T) {
	svc := newMockedPassengerService()
	assert := assert.New(t)

	newStatus := models.PassengerAccepted
	update := models.PassengerChangeSet{Status: &newStatus}
	driver := auth.UserClaims{ID: "456"}
	pass := validPassenger

	svc.passengerStore.On("GetByID", pass.ID).Return(&pass, nil).Once()
	svc.passengerStore.On("Count", mock.Anything).Return(2, nil).Once()
	svc.rideStore.On("GetByID", pass.RideID).Return(&ride1, nil).Once()
	svc.passengerStore.On("Update", mock.AnythingOfType("*models.Passenger")).Return(nil).Once()

	newPass, noErr := svc.passengerService.Update(&update, pass.ID, &driver)
	assert.Nil(noErr, "there should be no error")
	assert.NotNil(newPass, "there should be a new passenger")
	assert.Equal(models.PassengerAccepted, newPass.Status, "the passenger should be accepted")

	svc.passengerStore.On("GetByID", pass.ID).Return(&pass, nil).Once()
	notDriver := auth.UserClaims{ID: "xyz"}

	noPass, forb := svc.passengerService.Update(&update, pass.ID, &notDriver)
	assert.NotNil(forb, "there should be an error")
	assert.Equal(services.ErrForbidden, forb, "the user should be forbidden")
	assert.Nil(noPass, "there should be no passenger returned")

	badStatus := "bad status"
	badUpdates := models.PassengerChangeSet{Status: &badStatus}
	svc.passengerStore.On("GetByID", pass.ID).Return(&pass, nil).Once()

	noPass, validErr := svc.passengerService.Update(&badUpdates, pass.ID, &driver)
	assert.Nil(noPass, "there should be no passenger")
	assert.NotNil(validErr, "there should have been an error")

	svc.passengerStore.On("GetByID", pass.ID).Return(&pass, nil).Once()
	svc.passengerStore.On("Count", mock.Anything).Return(4, nil).Once()
	svc.rideStore.On("GetByID", pass.RideID).Return(&ride1, nil).Once()
	noPass, noMoreSeats := svc.passengerService.Update(&update, pass.ID, &driver)

	assert.NotNil(noMoreSeats, "there should have been an error")
	assert.Equal(services.ErrNoMoreSeats, noMoreSeats, "there should have been a at capacity error")
	assert.Nil(noPass, "there should be no passengers")
}

func TestGetPassenger(t *testing.T) {
	service := newMockedPassengerService()
	assert := assert.New(t)

	user := auth.UserClaims{ID: "xyz"}
	service.passengerStore.On("GetByID", validPassenger.ID).Return(&validPassenger, nil)

	pass, noErr := service.passengerService.Get(validPassenger.ID, &user)
	assert.Nil(noErr, "there should be no error")
	assert.NotNil(pass, "there should have been a passenger")

	badUser := auth.UserClaims{ID: "random", AuthLevel: services.UserLevel}
	noPass, err := service.passengerService.Get(validPassenger.ID, &badUser)
	assert.Equal(services.ErrForbidden, err, "the user should be forbidden")
	assert.Nil(noPass, "when a bad user there should be no passenger returned")
}

func TestGetAllPassengers(t *testing.T) {
	service := newMockedPassengerService()
	assert := assert.New(t)

	noPasses, err := service.passengerService.GetAll("", 15, services.UserLevel)
	assert.Equal(services.ErrNotAllowed, err, "should not be allowed")
	assert.Nil(noPasses, "there should be no passengers on error")

	lowLimit := -1
	service.passengerStore.On("GetAll", "", 15).Return([]*models.Passenger{&validPassenger}, nil)

	passengers, noErr := service.passengerService.GetAll("", lowLimit, services.AdminLevel)
	assert.Nil(noErr, "there should be no error")
	assert.NotNil(passengers, "there should be passengers")
	assert.Equal(1, len(passengers), "there should be a slice of 1 passenger")
}

func TestDeletePassenger(t *testing.T) {
	m := newMockedPassengerService()
	assert := assert.New(t)

	p := validPassenger
	m.passengerStore.On("GetByID", p.ID).Return(&p, nil)
	m.passengerStore.On("Delete", p.ID).Return(nil)
	user := auth.UserClaims{ID: p.DriverID}

	noErr := m.passengerService.Delete(p.ID, &user)
	assert.Nil(noErr, "there should be no error")

	notDriver := auth.UserClaims{ID: "adsfasdfa"}
	m.passengerStore.On("GetByID", p.ID).Return(&p, nil)
	err := m.passengerService.Delete(p.ID, &notDriver)
	assert.NotNil(err, "there should be an error")
	assert.Equal(services.ErrForbidden, err, "the user should have been forbidden")
}
