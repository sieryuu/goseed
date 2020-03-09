package models

import "time"

// Tenant represents the tenant model.
type Tenant struct {
	ID            uint      `xorm:"autoincr"`
	Name          string    `validate:"required" json:"name"`
	TenancyName   string    `validate:"required" json:"tenancy_name"`
	CreatedAt     time.Time `xorm:"created"`
	CreatedBy     string    `validate:"required"`
	LastUpdatedAt time.Time `xorm:"updated"`
	LastUpdatedBy string    `validate:"required"`
	IsActive      bool
}
