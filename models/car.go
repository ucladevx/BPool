package models

import (
	"time"
)

// Car model instance
type Car struct {
	ID        string    `json:"id"`
	Make      string    `json:"make"`
	Model     string    `json:"model"`
	Year      string    `json:"year"`
	Color     string    `json:"color"`
	UserID    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TODO: Validator

func (c *Car) String() string {
	return "todo"
}
