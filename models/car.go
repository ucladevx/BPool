package models

import (
	"errors"
	"fmt"
	"time"
)

// Car model instance
type Car struct {
	ID        string    `json:"id"`
	Make      string    `json:"make"`
	Model     string    `json:"model"`
	Year      int       `json:"year"`
	Color     string    `json:"color"`
	UserID    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// Validate validates car before insertion and updates
func (c *Car) Validate() []error {
	var errs []error

	if c.Make == "" {
		errs = append(errs, errors.New("please provide car's make"))
	}

	if c.Year > time.Now().Year()+1 || c.Year < 1000 {
		errs = append(errs, errors.New("year must cannot be too far in the past or in the future"))
	}

	if c.Color == "" {
		errs = append(errs, errors.New("please provide car's color"))
	}

	if c.UserID == "" {
		errs = append(errs, errors.New("please provide car's associated user id"))
	}

	return errs
}

func (c *Car) String() string {
	return fmt.Sprintf("<Car id:%s, owner_id:%s, make:%s, year:%d, color:%s", c.ID, c.UserID, c.Make, c.Year, c.Color)
}
