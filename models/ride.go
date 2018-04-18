package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/imdario/mergo"
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
		PricePerSeat int       `json:"price_per_seat" db:"price_per_seat"`
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
		PricePerSeat *int      `json:"price_per_seat"`
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
	errs := ""

	if r.DriverID == "" {
		errs += "Driver ID: the ride needs a driver;"
	}

	if r.DriverID == "" {
		errs += "Driver ID: the ride needs a driver;"
	}

	if r.CarID == "" {
		errs += "Car ID: the ride needs a car;"
	}

	if r.Seats < 0 {
		errs += "seats: there needs to be a positive number of seats;"
	}

	if r.StartCity == "" {
		errs += "start city: there must be a start city;"
	}

	if r.EndCity == "" {
		errs += "end city: there must be a end city;"
	}

	if r.StartLat < -90 || r.StartLat > 90 {
		errs += "start lat: latitude is between -90 and 90 degrees;"
	}

	if r.EndLat < -90 || r.EndLat > 90 {
		errs += "end lat: latitude is between -90 and 90 degrees;"
	}

	if r.StartLon < -180 || r.StartLon > 180 {
		errs += "start lon: longitude is between -180 and 180 degrees;"
	}

	if r.EndLon < -180 || r.EndLon > 180 {
		errs += "end lon: longitude is between -180 and 180 degrees;"
	}

	if r.PricePerSeat < 0 {
		errs += "price per seat: you must supply a price for the ride;"
	}

	now := time.Now()
	if r.StartDate.Before(now) {
		errs += "StartDate: the ride start date and time must be some time in the future;"
	}

	if errs != "" {
		return errors.New(errs)
	}

	return nil
}

// String returns a string representation of a ride
func (r *Ride) String() string {
	return fmt.Sprintf("<Ride id:%s driver:%s car:%s>", r.ID, r.DriverID, r.CarID)
}

// ApplyUpdates attempts to update ride and validates the updates
func (r *Ride) ApplyUpdates(o *RideChangeSet) error {
	newRide := Ride{}

	if err := mergo.Merge(&newRide, *r, mergo.WithOverride); err != nil {
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

	if !o.StartDate.IsZero() {
		newRide.StartDate = o.StartDate
	}

	err := newRide.Validate()
	if err != nil {
		return err
	}

	if err = mergo.Merge(r, newRide, mergo.WithOverride); err != nil {
		return errors.New("There was a problem applying updates to ride " + err.Error())
	}

	return nil
}
