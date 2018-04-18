package models_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/ucladevx/BPool/models"
)

func TestRideApplyUpdates(t *testing.T) {
	assert := assert.New(t)

	ride := models.Ride{
		DriverID:     "123",
		CarID:        "abc",
		Seats:        4,
		StartCity:    "Los Angeles",
		EndCity:      "San Francisco",
		StartLat:     89.934523,
		StartLon:     -90.934523,
		EndLat:       80.934523,
		EndLon:       80.934523,
		PricePerSeat: 15,
		Info:         "",
		StartDate:    time.Now().Add(time.Hour),
	}

	newSeats := 3
	newEndCity := "San Jose"
	newPricePerSeat := 20

	changes := models.RideChangeSet{
		Seats:        &newSeats,
		EndCity:      &newEndCity,
		PricePerSeat: &newPricePerSeat,
	}

	ride1 := ride

	err := ride1.ApplyUpdates(&changes)

	assert.Nil(err)
	assert.Equal(newSeats, ride1.Seats)
	assert.Equal(newEndCity, ride1.EndCity)
	assert.Equal(newPricePerSeat, ride1.PricePerSeat)
}
