package models

import (
	"fmt"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
)

type (
	// Ride is a ride entity
	Ride struct {
		ID           string    `json:"id" db:"id"`
		DriverID     string    `json:"driver_id" db:"driver_id"`
		CarID        string    `json:"car_id" db:"car_id"`
		Seats        int       `json:"seats" db:"seats"`
		StartCity    string    `json:"start_city" db:"start_city"`
		EndCity      string    `json:"end_city" db:"end_city"`
		StartLat     float64   `json:"start_dest_lat" db:"start_dest_lat"`
		StartLon     float64   `json:"start_dest_lon" db:"start_dest_lon"`
		EndLat       float64   `json:"end_dest_lat" db:"end_dest_lat"`
		EndLon       float64   `json:"end_dest_lon" db:"end_dest_lon"`
		PricePerSeat float64   `json:"price_per_seat" db:"price_per_seat"`
		Info         string    `json:"info" db:"info"`
		StartDate    time.Time `json:"start_date" db:"start_date"`
		CreatedAt    time.Time `json:"created_at" db:"created_at"`
		UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
	}

	// RideChangeSet is the fields that are modifiable in the ride
	RideChangeSet struct {
		DriverID     *string   `json:"-"`
		CarID        *string   `json:"car_id"`
		Seats        *int      `json:"seats"`
		StartCity    *string   `json:"start_city"`
		EndCity      *string   `json:"end_city"`
		StartLat     *float64  `json:"start_dest_lat"`
		StartLon     *float64  `json:"start_dest_lon"`
		EndLat       *float64  `json:"end_dest_lat"`
		EndLon       *float64  `json:"end_dest_lon"`
		PricePerSeat *float64  `json:"price_per_seat"`
		Info         *string   `json:"info"`
		StartDate    time.Time `json:"start_date"`
	}
)

// NewRide returns a Ride with the change set fields applied
func NewRide(r *RideChangeSet) *Ride {
	ride := Ride{}

	// ignore errors because we will jsut return a zero valued ride
	ride.ApplyUpdates(r)

	return &ride
}

// Validate validates a ride
func (r *Ride) Validate() error {
	now := time.Now()

	return validation.ValidateStruct(r,
		validation.Field(&r.DriverID, validation.Required),
		validation.Field(&r.CarID, validation.Required),
		validation.Field(&r.Seats, validation.Min(0)),
		validation.Field(&r.StartCity, validation.Required),
		validation.Field(&r.EndCity, validation.Required),
		validation.Field(&r.StartLat, validation.Required, validation.Min(-90.0), validation.Max(90.0)),
		validation.Field(&r.EndLat, validation.Required, validation.Min(-90.0), validation.Max(90.0)),
		validation.Field(&r.StartLon, validation.Required, validation.Min(-180.0), validation.Max(180.0)),
		validation.Field(&r.EndLon, validation.Required, validation.Min(-180.0), validation.Max(180.0)),
		validation.Field(&r.PricePerSeat, validation.Min(0.0), validation.Max(100000000.00)),
		validation.Field(&r.StartDate, validation.Required, validation.Min(now)),
	)
}

// String returns a string representation of a ride
func (r *Ride) String() string {
	return fmt.Sprintf("<Ride id:%s driver:%s car:%s>", r.ID, r.DriverID, r.CarID)
}

// ApplyUpdates attempts to update ride and validates the updates
func (r *Ride) ApplyUpdates(o *RideChangeSet) error {
	newRide := *r

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

	if !o.StartDate.IsZero() {
		newRide.StartDate = o.StartDate
	}

	err := newRide.Validate()
	if err != nil {
		return err
	}

	*r = newRide

	return nil
}
