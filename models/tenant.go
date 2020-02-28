package models

import "time"

// Tenant represents the tenant model.
type Tenant struct {
	ID            uint      `xorm:"'id' pk autoincr"`
	Name          string    `xorm:"varchar(255) not null unique"`
	TenancyName   string    `xorm:"varchar(255) not null"`
	CreatedAt     time.Time `xorm:"created"`
	CreatedBy     string    `xorm:"varchar(32) not null"`
	LastUpdatedAt time.Time `xorm:"updated"`
	LastUpdatedBy string    `xorm:"varchar(32) not null"`
	IsActive      bool
}
