package user

import "goseed/models"

// Repository represents the user's repository contract.
type Repository interface {
	Find() (*[]models.User, error)
	Create(user *models.User) (int64, error)
	GetByUsername(username string) (*models.User, error)
}
