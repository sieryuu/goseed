package models

import (
	"time"
)

// Article represents the article model.
type Article struct {
	ID            uint      `xorm:"autoincr"`
	TenantID      uint      `validate:"required" json:"tenant_id"`
	Title         string    `validate:"required" json:"title"`
	CreatedAt     time.Time `xorm:"created"`
	CreatedBy     string    `validate:"required"`
	LastUpdatedAt time.Time `xorm:"updated"`
	LastUpdatedBy string    `validate:"required"`
}
