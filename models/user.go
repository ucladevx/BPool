package models

import (
	"fmt"
	"time"
)

// User model instance
type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Image     string    `json:"image"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u *User) String() string {
	return fmt.Sprintf("<User name:%s id:%d email:%s>", u.Name, u.ID, u.Email)
}
