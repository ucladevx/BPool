package services

import (
	"sort"

	"github.com/ucladevx/BPool/interfaces"
	"github.com/ucladevx/BPool/models"
	"github.com/ucladevx/BPool/stores"
)

// FeedService provides all data for the feed
type FeedService struct {
	rideStore RideStore
	logger    interfaces.Logger
}

// NewFeedService creates a new feed service
func NewFeedService(store RideStore, l interfaces.Logger) *FeedService {
	return &FeedService{
		rideStore: store,
		logger:    l,
	}
}

// GetUserRides gets all user's associated rides
func (f *FeedService) GetUserRides(userID string) ([]*models.Ride, error) {
	clauses := []stores.QueryModifier{
		stores.QueryMod("driver_id", stores.EQ, userID),
	}

	rides, err := f.rideStore.WhereMany(clauses)
	if err != nil {
		return nil, err
	}

	passengerRides, err := f.rideStore.GetAllWherePassenger(userID)
	if err != nil {
		return nil, err
	}

	allRides := append(rides, passengerRides...)

	// Sort rides by descending start time
	sort.Slice(allRides, func(i, j int) bool { return allRides[j].StartDate.Before(allRides[i].StartDate) })
	return allRides, nil
}
