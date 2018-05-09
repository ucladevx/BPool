package models

import (
	"fmt"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
)

const (
	// PassengerAccepted means the passenger is in the ride
	PassengerAccepted = "accepted"
	// PassengerInterested means the passenger would like to join the ride
	PassengerInterested = "interested"
	// PassengerRejected means the passenger cannot join the ride
	PassengerRejected = "rejected"
)

type (
	// Passenger is the entity for the ride passenger relation
	Passenger struct {
		ID          string    `json:"id" db:"id"`
		DriverID    string    `json:"driver_id" db:"driver_id"`
		PassengerID string    `json:"passenger_id" db:"passenger_id"`
		RideID      string    `json:"ride_id" db:"ride_id"`
		Status      string    `json:"status" db:"status"`
		CreatedAt   time.Time `json:"created_at" db:"created_at"`
		UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	}

	// PassengerChangeSet is what is allowed to be changed
	PassengerChangeSet struct {
		PassengerID *string `json:"passenger_id"`
		RideID      *string `json:"ride_id"`
		Status      *string `json:"status"`
	}
)

// NewPassenger creates and validates a passenger from a changeset
func NewPassenger(p *PassengerChangeSet) (*Passenger, error) {
	passenger := Passenger{}

	if err := passenger.ApplyUpdates(p); err != nil {
		return nil, err
	}

	return &passenger, nil
}

func (p *Passenger) String() string {
	return fmt.Sprintf("<Passenger id:%s ride:%s driver:%s passenger_id:%s>", p.ID, p.RideID, p.DriverID, p.PassengerID)
}

// Validate validates the passenger
func (p *Passenger) Validate() error {
	return validation.ValidateStruct(p,
		validation.Field(&p.RideID, validation.Required),
		validation.Field(&p.PassengerID, validation.Required),
		validation.Field(&p.Status, validation.Required, validation.In(PassengerAccepted, PassengerInterested, PassengerRejected)),
	)
}

// ApplyUpdates attempts to update the passenger with the given changes
func (p *Passenger) ApplyUpdates(o *PassengerChangeSet) error {
	newPassenger := *p

	if o.RideID != nil {
		newPassenger.RideID = *o.RideID
	}

	if o.PassengerID != nil {
		newPassenger.PassengerID = *o.PassengerID
	}

	if o.Status != nil {
		newPassenger.Status = *o.Status
	}

	if err := newPassenger.Validate(); err != nil {
		return err
	}

	*p = newPassenger

	return nil
}
