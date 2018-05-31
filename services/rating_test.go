package services_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ucladevx/BPool/mocks"
	"github.com/ucladevx/BPool/models"
	"github.com/ucladevx/BPool/services"
	"github.com/ucladevx/BPool/stores"
	"github.com/ucladevx/BPool/stores/postgres"
	"github.com/ucladevx/BPool/utils/auth"
)

var (
	rating1 = models.Rating{
		ID:      "goodRating",
		Rating:  5,
		RideID:  ride1.ID,
		RaterID: johnDoe.ID,
		RateeID: janeSmith.ID,
		Comment: "good ride",
	}

	rating2 = models.Rating{
		ID:      "poorRating",
		Rating:  2,
		RideID:  ride1.ID,
		RaterID: "dawg",
		RateeID: janeSmith.ID,
		Comment: "",
	}
)

func newRatingService(store *mocks.RatingStore, p *services.PassengerService) *services.RatingService {
	logger := mocks.Logger{}
	return services.NewRatingService(store, p, logger)
}

func TestRatingGetByID(t *testing.T) {
	store := new(mocks.RatingStore)
	p := newMockedPassengerService()
	service := newRatingService(store, p.passengerService)
	assert := assert.New(t)

	// Return error if user does not have proper auth level
	store.On("GetByID", "forbidden").Return(nil, services.ErrForbidden)

	user := auth.UserClaims{
		ID:        "testUser",
		Email:     "test@gmail.com",
		AuthLevel: services.UserLevel,
	}

	noRating, err := service.GetByID("forbidden", &user)

	assert.Nil(noRating, "for unauthorized user, there should be no rating")
	assert.Equal(services.ErrForbidden, err, "if user unauthorized, should return error")

	// Return error if rating does not exist
	store.On("GetByID", "badID").Return(nil, postgres.ErrNoRatingFound)

	user = auth.UserClaims{
		ID:        "testUser",
		Email:     "test@gmail.com",
		AuthLevel: services.AdminLevel,
	}

	noRating, err = service.GetByID("badID", &user)

	assert.Nil(noRating, "for a bad id there should be no rating")
	assert.Equal(postgres.ErrNoRatingFound, err, "if no rating found should return no rating found error")

	// Return reating if rating does exist and user is authorized
	store.On("GetByID", rating1.ID).Return(&rating1, nil)

	user = auth.UserClaims{
		ID:        "testUser",
		Email:     "test@gmail.com",
		AuthLevel: services.AdminLevel,
	}

	rating, noErr := service.GetByID(rating1.ID, &user)

	assert.Nil(noErr, "no error if rating exists")
	assert.Equal(rating.ID, rating1.ID, "should return rating if no problem with query")
}

func TestRatingCreate(t *testing.T) {
	store := new(mocks.RatingStore)
	p := newMockedPassengerService()
	service := newRatingService(store, p.passengerService)
	assert := assert.New(t)

	user := auth.UserClaims{
		ID:        "testUser",
		Email:     "test@gmail.com",
		AuthLevel: services.AdminLevel,
	}

	passengers := []*models.Passenger{
		&validPassenger,
	}

	// Return err on duplicate ratings
	queryModifiers := []stores.QueryModifier{
		stores.QueryMod("ride_id", stores.EQ, rating1.RideID),
	}

	duplicates := []*models.Rating{
		&rating1, &rating2,
	}

	p.rideStore.On("GetByID", rating1.RideID).Return(&ride1, nil)
	p.passengerStore.On("WhereMany", queryModifiers).Return(passengers)
	store.On("WhereMany", queryModifiers).Return(duplicates)

	toAdd := rating1

	noRating, err := service.Create(toAdd, &user)

	assert.Nil(noRating, "fail rating creation if duplicate")
	assert.Equal(services.ErrDuplciateRating, err, "return error on duplicate rating")
}

func TestRatingGetRatingByUserID(t *testing.T) {
	store := new(mocks.RatingStore)
	p := newMockedPassengerService()
	service := newRatingService(store, p.passengerService)
	assert := assert.New(t)

	// Returns average of all ratings for ride
	store.On("Average", janeSmith.ID).Return(float32(3.5), nil)

	rating, err := service.GetRatingByUserID(janeSmith.ID)

	assert.Nil(err, "should not return error if ratee has any ratings")
	assert.Equal(rating, float32(3.5), "successfully returns average rating")
}

func TestRatingDelete(t *testing.T) {
	store := new(mocks.RatingStore)
	p := newMockedPassengerService()
	service := newRatingService(store, p.passengerService)
	assert := assert.New(t)

	// Successfully delete is ride exists and user is admin level
	store.On("GetByID", rating1.ID).Return(&rating1, nil)
	store.On("Delete", rating1.ID).Return(nil)
	user := auth.UserClaims{
		ID:        "testUser",
		Email:     "test@gmail.com",
		AuthLevel: services.AdminLevel,
	}

	err := service.Delete(rating1.ID, &user)
	assert.Nil(err, "successful delete should not return an error")

	// Return err if rating does not exist
	store.On("GetByID", "badID").Return(nil, postgres.ErrNoRatingFound)
	err = service.Delete("badID", &user)

	assert.NotNil(err, "for a bad id there should be no rating")
	assert.Equal(postgres.ErrNoRatingFound, err, "if no rating found should return no rating found error")

	// Return error if user is not admin level
	user = auth.UserClaims{
		ID:        "testUser",
		Email:     "test@gmail.com",
		AuthLevel: services.UserLevel,
	}

	err = service.Delete(rating1.ID, &user)

	assert.NotNil(err, "for unauthorized user, user cannot delete rating")
	assert.Equal(services.ErrForbidden, err, "if user unauthorized, should return error")
}
