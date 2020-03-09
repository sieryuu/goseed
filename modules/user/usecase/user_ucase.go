package usecase

import (
	"errors"
	"goseed/models"
	"goseed/modules/user"
	"goseed/utils/hashutil"
)

type userUsecase struct {
	userRepo user.Repository
}

// NewUserUsecase will create a new userUsecase object representation of user.Usecase interface.
func NewUserUsecase(a user.Repository) user.Usecase {
	return &userUsecase{a}
}

func (a *userUsecase) Login(name string, password []byte) (bool, error) { // Login always return "invalid username or password" if fail
	user, err := a.userRepo.GetByUsername(name)
	if err != nil {
		return false, err
	}

	// wrong user name given
	if user == nil {
		return false, nil
	}

	isSuccess := hashutil.ComparePasswords(user.PasswordHash, password)

	if isSuccess == false {
		return false, nil
	}

	return isSuccess, nil
}

func (a *userUsecase) Create(user *models.User) error {
	user, err := a.userRepo.GetByUsername(user.Username)
	// something bad happens
	if err != nil {
		return err
	}
	// username already exists
	if user != nil {
		return errors.New("UserCreateUsernameExistsMsg")
	}

	_, err = a.userRepo.Create(user)
	return err
}
