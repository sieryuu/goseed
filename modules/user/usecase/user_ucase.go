package usecase

import (
	"errors"
	"goseed/models"
	"goseed/modules/user"
	"goseed/modules/user/delivery/dto"
	"goseed/utils/encryption"
)

type userUsecase struct {
	userRepo user.Repository
}

// NewUserUsecase will create a new userUsecase object representation of user.Usecase interface.
func NewUserUsecase(a user.Repository) user.Usecase {
	return &userUsecase{a}
}

func (a *userUsecase) Login(name string, password []byte) (bool, error) {
	user, err := a.userRepo.GetByUsername(name)
	// something bad happens
	if err != nil {
		return false, err
	}
	// wrong user name given
	if user == nil {
		return false, errors.New("invalid user or password")
	}

	isSuccess := encryption.ComparePasswords(user.PasswordHash, password)

	if isSuccess == false {
		return false, errors.New("invalid user or password")
	}

	return isSuccess, nil
}

func (a *userUsecase) Create(userCreation *dto.UserCreation) (*models.User, error) {
	user, err := a.userRepo.GetByUsername(userCreation.Username)
	// something bad happens
	if err != nil {
		return nil, err
	}
	// wrong user name given
	if user != nil {
		return nil, errors.New("user already exists")
	}

	user = userCreation.GetUser()
	_, err = a.userRepo.Create(user)
	return user, err
}
