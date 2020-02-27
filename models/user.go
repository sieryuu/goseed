package models

import "time"

// User represents the user model.
type User struct {
	ID           uint   `xorm:"'id' pk autoincr"`
	Username     string `xorm:"varchar(32) not null unique"`
	FirstName    string `xorm:"varchar(32) not null "`
	LastName     string `xorm:"varchar(32) not null "`
	PasswordHash string
	CreatedAt    time.Time `xorm:"created"`
	UpdatedAt    time.Time `xorm:"updated"`
}
