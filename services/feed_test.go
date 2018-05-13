package services_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/ucladevx/BPool/mocks"
	"github.com/ucladevx/BPool/models"
	"github.com/ucladevx/BPool/services"
)

type mockedFeedService struct {
	rideStore   *mocks.RideStore
	feedService *services.FeedService
}

func newMockedFeedService() *mockedFeedService {
	logger := mocks.Logger{}
	rideStore := new(mocks.RideStore)
	feedService := services.NewFeedService(rideStore, logger)

	return &mockedFeedService{
		rideStore:   rideStore,
		feedService: feedService,
	}
}

func TestGetUserFeedRides(t *testing.T) {
	feed := newMockedFeedService()
	assert := assert.New(t)

	userID := ride1.DriverID

	feed.rideStore.On("WhereMany", mock.Anything).Return([]*models.Ride{&ride2}, nil)
	feed.rideStore.On("GetAllWherePassenger", userID).Return([]*models.Ride{&ride1}, nil)

	rides, err := feed.feedService.GetUserRides(userID)
	assert.Nil(err, "there should be no error")
	assert.NotNil(rides, "there should be some rides")
	assert.Equal(2, len(rides), "there should be two rides for the user")
	assert.Equal(&ride1, rides[0], "the first ride should be the one with the later start date")
	assert.Equal(&ride2, rides[1], "the second ride should have a start date before the first ride")
}
