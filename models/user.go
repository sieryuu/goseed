package models

import "time"

// User represents the user model.
type User struct {
	ID            uint      `xorm:"autoincr"`
	Username      string    `validate:"required"`
	FirstName     string    `validate:"required"`
	LastName      string    `validate:"required"`
	PasswordHash  string    `validate:"required"`
	CreatedAt     time.Time `xorm:"created"`
	CreatedBy     string    `validate:"required"`
	LastUpdatedAt time.Time `xorm:"updated"`
	LastUpdatedBy string    `validate:"required"`
	IsActive      bool
}
