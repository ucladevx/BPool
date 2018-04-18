package models

import (
	"errors"
	"time"

	"github.com/imdario/mergo"
)

type (
	// Ride is a ride entity
	Ride struct {
		ID           string
		DriverID     string    `json:"driver_id" db:"driver_id"`
		CarID        string    `json:"car_id" db:"car_id"`
		Seats        int       `json:"seats" db:"seats"`
		StartCity    string    `json:"start_city" db:"start_city"`
		EndCity      string    `json:"end_city" db:"end_city"`
		StartLat     float64   `json:"start_dest_lat" db:"start_dest_lat"`
		StartLon     float64   `json:"start_dest_lon" db:"start_dest_lon"`
		EndLat       float64   `json:"end_dest_lat" db:"end_dest_lat"`
		EndLon       float64   `json:"end_dest_lon" db:"end_dest_lon"`
		PricePerSeat string    `json:"price_per_seat" db:"price_per_seat"`
		Info         string    `json:"info" db:"info"`
		StartDate    time.Time `json:"start_date" db:"start_date"`
		CreatedAt    time.Time `json:"created_at" db:"created_at"`
		UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
	}

	// RideChangeSet is the fields that are modifiable in the ride
	RideChangeSet struct {
		DriverID     *string    `json:"-"`
		CarID        *string    `json:"car_id"`
		Seats        *int       `json:"seats"`
		StartCity    *string    `json:"start_city"`
		EndCity      *string    `json:"end_city"`
		StartLat     *float64   `json:"start_dest_lat"`
		StartLon     *float64   `json:"start_dest_lon"`
		EndLat       *float64   `json:"end_dest_lat"`
		EndLon       *float64   `json:"end_dest_lon"`
		PricePerSeat *string    `json:"price_per_seat"`
		Info         *string    `json:"info"`
		StartDate    *time.Time `json:"start_date"`
	}
)

// Validate validates a ride
func (r *Ride) Validate() error {
	return nil
}

// ApplyUpdates attempts to update ride and validates the updates
func (r *Ride) ApplyUpdates(o *RideChangeSet) error {
	newRide := Ride{}

	if err := mergo.Merge(&newRide, *r); err != nil {
		return errors.New("There was a problem applying updates to ride")
	}

	if o.DriverID != nil {
		newRide.DriverID = *o.DriverID
	}

	if o.CarID != nil {
		newRide.CarID = *o.CarID
	}

	if o.Seats != nil {
		newRide.Seats = *o.Seats
	}

	if o.StartCity != nil {
		newRide.StartCity = *o.StartCity
	}

	if o.EndCity != nil {
		newRide.EndCity = *o.EndCity
	}

	if o.StartLat != nil {
		newRide.StartLat = *o.StartLat
	}

	if o.StartLon != nil {
		newRide.StartLon = *o.StartLon
	}

	if o.EndLat != nil {
		newRide.EndLat = *o.EndLat
	}

	if o.EndLon != nil {
		newRide.EndLon = *o.EndLon
	}

	if o.PricePerSeat != nil {
		newRide.PricePerSeat = *o.PricePerSeat
	}

	if o.Info != nil {
		newRide.Info = *o.Info
	}

	if o.StartDate != nil {
		newRide.StartDate = *o.StartDate
	}

	err := newRide.Validate()
	if err != nil {
		return err
	}

	if err = mergo.Merge(&r, newRide, mergo.WithOverride); err != nil {
		return errors.New("There was a problem applying updates to ride")
	}

	return nil
}
