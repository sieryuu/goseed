package http_test

import (
	"encoding/json"
	"errors"
	"goseed/modules/user/delivery/dto"
	"goseed/modules/user/mocks"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	_userhttp "goseed/modules/user/delivery/http"
)

type Token struct {
	token string
}

func TestLogin(t *testing.T) {
	mockUserUsecase := new(mocks.Usecase)

	mockuserLogin := &dto.UserLogin{
		Username: "admin",
		Password: "MyLongPassword",
	}

	mockUserLoginJson, err := json.Marshal(mockuserLogin)
	assert.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		mockUserUsecase.On("Login", mock.Anything, mock.Anything).Return(true, nil).Once()

		e := echo.New()
		req, err := http.NewRequest(echo.POST, "/v1/login", strings.NewReader(string(mockUserLoginJson)))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handler := _userhttp.UserHandler{
			UserUsecase: mockUserUsecase,
		}
		err = handler.Login(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		token := new(Token)
		body, _ := ioutil.ReadAll(rec.Result().Body)
		json.Unmarshal(body, token)
		assert.NotNil(t, token)

		mockUserUsecase.AssertExpectations(t)
	})

	t.Run("login-fail", func(t *testing.T) {
		mockUserUsecase.On("Login", mock.Anything, mock.Anything).Return(false, errors.New("")).Once()

		e := echo.New()
		req, err := http.NewRequest(echo.POST, "/v1/login", strings.NewReader(string(mockUserLoginJson)))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handler := _userhttp.UserHandler{
			UserUsecase: mockUserUsecase,
		}
		err = handler.Login(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		mockUserUsecase.AssertExpectations(t)
	})
}

func TestCreate(t *testing.T) {
	mockUserUsecase := new(mocks.Usecase)

	mockUserCreation := &dto.UserCreation{
		Username:  "admin",
		FirstName: "John",
		LastName:  "Snow",
		Password:  "LonglongPassword",
	}

	mockUserCreationJson, err := json.Marshal(mockUserCreation)
	assert.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		mockUserUsecase.On("Create", mock.Anything).Return(nil, nil).Once()

		e := echo.New()
		req, err := http.NewRequest(echo.POST, "/v1/users", strings.NewReader(string(mockUserCreationJson)))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handler := _userhttp.UserHandler{
			UserUsecase: mockUserUsecase,
		}
		err = handler.Create(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)
		mockUserUsecase.AssertExpectations(t)
	})

	t.Run("user-exist", func(t *testing.T) {
		mockUserUsecase.On("Create", mock.Anything).Return(nil, errors.New("")).Once()

		e := echo.New()
		req, err := http.NewRequest(echo.POST, "/v1/users", strings.NewReader(string(mockUserCreationJson)))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handler := _userhttp.UserHandler{
			UserUsecase: mockUserUsecase,
		}
		err = handler.Create(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		mockUserUsecase.AssertExpectations(t)
	})
}
