package http_test

import (
	"encoding/json"
	"errors"
	"goseed/models"
	"goseed/modules/user/mocks"
	"goseed/utils/echoutil"
	"goseed/utils/hashutil"
	"goseed/utils/httputil"
	"goseed/utils/i18nutil"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	_userhttp "goseed/modules/user/delivery/http"
)

func TestCreate(t *testing.T) {
	e := echoutil.MockEcho()
	bundle := i18nutil.GetTestingi18nBundle()
	l := i18n.NewLocalizer(bundle)

	mockUser := &models.User{
		Username:     "admin",
		FirstName:    "John",
		LastName:     "Snow",
		PasswordHash: hashutil.HashAndSalt([]byte("MyLongLongPassword")),
	}

	mockUserCreationJson, err := json.Marshal(mockUser)
	assert.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		mockUserUsecase := new(mocks.Usecase)
		mockUserUsecase.On("Create", mock.Anything).Return(nil).Once()

		req, err := http.NewRequest(echo.POST, "/v1/users", strings.NewReader(string(mockUserCreationJson)))
		assert.NoError(t, err)

		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", echoutil.GetLoginTkn(mockUser.Username))

		handler := _userhttp.UserHandler{
			UserUsecase: mockUserUsecase,
			I18nBundle:  bundle,
		}
		err = handler.Create(c)
		assert.NoError(t, err)

		body, err := ioutil.ReadAll(rec.Body)
		assert.NoError(t, err)

		responseBody := new(httputil.ResponseBody)
		err = json.Unmarshal(body, responseBody)
		assert.NoError(t, err)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)

		lUserCreateSuccessTitle, err := l.LocalizeMessage(&i18n.Message{ID: "UserCreateSuccessTitle"})
		assert.NoError(t, err)
		lUserCreateSuccessMsg, err := l.Localize(&i18n.LocalizeConfig{
			MessageID: "UserCreateSuccessMsg",
			TemplateData: map[string]string{
				"username": mockUser.Username,
			},
		})
		assert.NoError(t, err)
		assert.Contains(t, string(body), lUserCreateSuccessTitle)
		assert.Contains(t, string(body), lUserCreateSuccessMsg)

		mockUserUsecase.AssertExpectations(t)
	})

	t.Run("username-exists", func(t *testing.T) {
		mockUserUsecase := new(mocks.Usecase)
		mockUserUsecase.On("Create", mock.Anything).Return(errors.New("UserCreateUsernameExistsMsg")).Once()

		req, err := http.NewRequest(echo.POST, "/v1/users", strings.NewReader(string(mockUserCreationJson)))
		assert.NoError(t, err)

		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", echoutil.GetLoginTkn(mockUser.Username))

		handler := _userhttp.UserHandler{
			UserUsecase: mockUserUsecase,
			I18nBundle:  bundle,
		}
		err = handler.Create(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusConflict, rec.Code)

		body, err := ioutil.ReadAll(rec.Body)
		assert.NoError(t, err)

		lUserCreateFailTitle, err := l.LocalizeMessage(&i18n.Message{ID: "UserCreateFailTitle"})
		assert.NoError(t, err)
		lUserCreateUsernameExistsMsg, err := l.LocalizeMessage(&i18n.Message{ID: "UserCreateUsernameExistsMsg"})
		assert.NoError(t, err)
		assert.Contains(t, string(body), lUserCreateFailTitle)
		assert.Contains(t, string(body), lUserCreateUsernameExistsMsg)

		mockUserUsecase.AssertExpectations(t)
	})
}
