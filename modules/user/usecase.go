package user

import (
	"goseed/models"
	"goseed/modules/user/delivery/dto"
)

// Usecase represents the user's usecase.
type Usecase interface {
	Login(name string, password []byte) (bool, error)
	Create(user *dto.UserCreation) (*models.User, error)
}
