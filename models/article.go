package models

import (
	"time"
)

// Article represents the article model.
type Article struct {
	ID        uint `xorm:"'id' pk autoincr"`
	Title     string
	CreatedAt time.Time `xorm:"created"`
	UpdatedAt time.Time `xorm:"updated"`
	// DeletedAt *time.Time `sql:"index"`
}
