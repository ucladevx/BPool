package models

import (
	"fmt"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
)

// Rating model instance
type Rating struct {
	ID        string    `json:"id"`
	Rating    int       `json:"rating"`
	RideID    string    `json:"ride_id" db:"ride_id"`
	RaterID   string    `json:"rater_id"`
	RateeID   string    `json:"ratee_id"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// Validate validates rating before insertion
func (r *Rating) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.Rating, validation.Required, validation.Min(1)),
		validation.Field(&r.RideID, validation.Required),
		validation.Field(&r.RaterID, validation.Required),
		validation.Field(&r.RateeID, validation.Required),
	)
}

func (r *Rating) String() string {
	return fmt.Sprintf("<Rating id: %s, rating: %d, rideID: %s, rater: %s, ratee: %s",
		r.ID, r.Rating, r.RideID, r.RaterID, r.RateeID,
	)
}
