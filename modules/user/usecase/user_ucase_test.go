package usecase_test

import (
	"errors"
	"goseed/models"
	"goseed/modules/user/delivery/dto"
	"goseed/modules/user/mocks"
	"goseed/modules/user/usecase"
	"goseed/utils/hashhlpr"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLogin(t *testing.T) {
	var (
		username = "admin"
		password = []byte("MySecureL0ngPassw0rd")
	)

	mockUser := models.User{
		ID:           1,
		Username:     username,
		PasswordHash: "$2a$04$uedlNViwUDDo.oO2DfsUM.CFs/TRlQk2j5/WNOJbeD4KIs88iFKJy",
	}

	t.Run("success", func(t *testing.T) {
		mockUserRepo := new(mocks.Repository)
		mockUserRepo.On("GetByUsername", username).Return(&mockUser, nil).Once()
		usecase := usecase.NewUserUsecase(mockUserRepo)

		result, err := usecase.Login(username, password)

		assert.NoError(t, err)
		assert.True(t, result)
	})

	t.Run("invalid-password", func(t *testing.T) {
		mockUserRepo := new(mocks.Repository)
		// assign a wrong password hash
		mockUser.PasswordHash = "WrongPassword"

		mockUserRepo.On("GetByUsername", username).Return(&mockUser, nil).Once()
		usecase := usecase.NewUserUsecase(mockUserRepo)

		result, err := usecase.Login(username, password)

		assert.Error(t, err)
		assert.False(t, result)
	})

	t.Run("invalid-user", func(t *testing.T) {
		mockUserRepo := new(mocks.Repository)
		mockUserRepo.On("GetByUsername", mock.Anything).Return(&models.User{}, errors.New("")).Once()
		usecase := usecase.NewUserUsecase(mockUserRepo)

		result, err := usecase.Login(username, password)

		assert.Error(t, err)
		assert.False(t, result)
	})
}

func TestCreate(t *testing.T) {
	mockUserCreation := &dto.UserCreation{
		Username: "admin",
		Password: "MySecureL0ngPassw0rd",
	}

	mockUser := &models.User{
		Username:     "admin",
		PasswordHash: "$2a$04$uedlNViwUDDo.oO2DfsUM.CFs/TRlQk2j5/WNOJbeD4KIs88iFKJy",
		FirstName:    "John",
		LastName:     "Snow",
	}

	t.Run("success", func(t *testing.T) {
		mockUserRepo := new(mocks.Repository)
		mockUserRepo.On("GetByUsername", mock.Anything).Return(nil, nil).Once()
		mockUserRepo.On("Create", mock.AnythingOfType("*models.User")).Return(int64(1), nil).Once()

		usecase := usecase.NewUserUsecase(mockUserRepo)
		user, err := usecase.Create(mockUserCreation)

		assert.NoError(t, err)
		assert.Equal(t, mockUser.Username, user.Username)
		assert.True(t, hashhlpr.ComparePasswords(user.PasswordHash, []byte(mockUserCreation.Password)))
	})

	t.Run("user-exist", func(t *testing.T) {
		mockUserRepo := new(mocks.Repository)
		mockUserRepo.On("GetByUsername", mock.Anything).Return(mockUser, nil).Once()
		mockUserRepo.On("Create", mock.AnythingOfType("*models.User")).Return(int64(0), nil).Once()

		usecase := usecase.NewUserUsecase(mockUserRepo)
		user, err := usecase.Create(mockUserCreation)

		assert.Error(t, err)
		assert.Nil(t, user)
	})
}
