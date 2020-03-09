package http_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"goseed/modules/user/delivery/dto"
	_userhttp "goseed/modules/user/delivery/http"
	"goseed/modules/user/mocks"
	"goseed/utils/i18nutil"
)

func TestLogin(t *testing.T) {
	e := echo.New()
	bundle := i18nutil.GetTestingi18nBundle()
	l := i18n.NewLocalizer(bundle)

	mockuserLogin := &dto.UserLogin{
		Username: "admin",
		Password: "MyLongPassword",
	}
	mockUserLoginJson, err := json.Marshal(mockuserLogin)
	assert.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		mockUserUsecase := new(mocks.Usecase)
		mockUserUsecase.On("Login", mock.Anything, mock.Anything).Return(true, nil).Once()

		req, err := http.NewRequest(echo.POST, "/me/v1/login", strings.NewReader(string(mockUserLoginJson)))
		assert.NoError(t, err)

		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		res := httptest.NewRecorder()
		c := e.NewContext(req, res)

		handler := _userhttp.MeHandler{
			UserUsecase: mockUserUsecase,
			I18nBundle:  bundle,
		}
		err = handler.Login(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, res.Code)

		body, err := ioutil.ReadAll(res.Body)
		assert.NoError(t, err)

		lMeLoginSuccessTitle, err := l.LocalizeMessage(&i18n.Message{ID: "MeLoginSuccessTitle"})
		assert.NoError(t, err)
		lMeLoginSuccessMsg, err := l.Localize(&i18n.LocalizeConfig{
			MessageID: "MeLoginSuccessMsg",
			TemplateData: map[string]string{
				"username": mockuserLogin.Username,
			},
		})
		assert.NoError(t, err)
		assert.Contains(t, string(body), lMeLoginSuccessTitle) // "login successful"
		assert.Contains(t, string(body), lMeLoginSuccessMsg)   // "Welcome User"

		mockUserUsecase.AssertExpectations(t)
	})

	t.Run("login-fail", func(t *testing.T) {
		mockUserUsecase := new(mocks.Usecase)
		mockUserUsecase.On("Login", mock.Anything, mock.Anything).Return(false, nil).Once()

		req, err := http.NewRequest(echo.POST, "/me/v1/login", strings.NewReader(string(mockUserLoginJson)))
		assert.NoError(t, err)

		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		res := httptest.NewRecorder()
		c := e.NewContext(req, res)

		handler := _userhttp.MeHandler{
			UserUsecase: mockUserUsecase,
			I18nBundle:  bundle,
		}
		err = handler.Login(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, res.Code)

		body, err := ioutil.ReadAll(res.Body)
		assert.NoError(t, err)

		lMeLoginInvalidTitle, err := l.LocalizeMessage(&i18n.Message{ID: "MeLoginInvalidTitle"})
		assert.NoError(t, err)
		lMeLoginInvalidMsg, err := l.LocalizeMessage(&i18n.Message{ID: "MeLoginInvalidMsg"})
		assert.NoError(t, err)
		assert.Contains(t, string(body), lMeLoginInvalidTitle)
		assert.Contains(t, string(body), lMeLoginInvalidMsg) // "invalid username or password"

		mockUserUsecase.AssertExpectations(t)
	})
}
