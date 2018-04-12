package models

import (
	"fmt"
	"time"
)

// User model instance
type User struct {
	ID           string    `json:"id"`
	FirstName    string    `json:"first_name" db:"first_name"`
	LastName     string    `json:"last_name" db:"last_name"`
	Email        string    `json:"email"`
	ProfileImage string    `json:"profile_image" db:"profile_image"`
	AuthLevel    int       `json:"auth_level" db:"auth_level"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

func (u *User) String() string {
	return fmt.Sprintf("<User name:%s id:%s email:%s auth_level:%d>", u.FirstName+" "+u.LastName, u.ID, u.Email, u.AuthLevel)
}
