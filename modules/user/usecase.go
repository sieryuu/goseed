package user

import (
	"goseed/models"
)

// Usecase represents the user's usecase.
type Usecase interface {
	Login(name string, password []byte) (bool, error)
	Create(user *models.User) error
}
