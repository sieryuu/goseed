package dto

import (
	"goseed/models"
	"goseed/utils/hashhlpr"
)

type (
	// UserLogin represents user login dto.
	UserLogin struct {
		Username string
		Password string
	}

	// UserCreation represents create user dto.
	UserCreation struct {
		Username  string
		Password  string
		FirstName string
		LastName  string
	}
)

// GetUser will map UserCreation to models.User.
func (u *UserCreation) GetUser() *models.User {
	return &models.User{
		Username:     u.Username,
		PasswordHash: hashhlpr.HashAndSalt([]byte(u.Password)),
		FirstName:    u.FirstName,
		LastName:     u.LastName,
	}
}
