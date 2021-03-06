package http_test

import (
	"encoding/json"
	"errors"
	"goseed/models"
	"goseed/modules/tenant/mocks"
	"goseed/utils/echoutil"
	"goseed/utils/i18nutil"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bxcodec/faker"
	"github.com/labstack/echo/v4"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	_tenanthttp "goseed/modules/tenant/delivery/http"
)

func TestFind(t *testing.T) {
	var mockTenant models.Tenant
	err := faker.FakeData(&mockTenant)
	assert.NoError(t, err)

	mockListTenant := make([]models.Tenant, 0)
	mockListTenant = append(mockListTenant, mockTenant)

	mockUsecase := new(mocks.Usecase)
	mockUsecase.On("Find").Return(&mockListTenant, nil)

	e := echo.New()
	req, err := http.NewRequest(echo.GET, "/v1/tenants", strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := _tenanthttp.TenantHandler{
		TenantUsecase: mockUsecase,
	}
	err = handler.Find(c)
	// require.NoError will halt the test if found any error
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	mockUsecase.AssertExpectations(t)
}

func TestCreate(t *testing.T) {
	bundle := i18nutil.GetTestingi18nBundle()
	l := i18nutil.GetTestingLocalizer()

	mockTenant := &models.Tenant{
		Name:        "tenant",
		TenancyName: "tenant name",
	}

	json, err := json.Marshal(mockTenant)
	assert.NoError(t, err)

	e := echoutil.MockEcho()

	t.Run("success", func(t *testing.T) {
		mockUsecase := new(mocks.Usecase)
		mockUsecase.On("Create", mock.Anything).Return(nil).Once()

		req, err := http.NewRequest(echo.POST, "/v1/tenants", strings.NewReader(string(json)))
		assert.NoError(t, err)

		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		res := httptest.NewRecorder()
		c := e.NewContext(req, res)
		c.Set("user", echoutil.GetLoginTkn("admin"))

		handler := _tenanthttp.TenantHandler{
			TenantUsecase: mockUsecase,
			I18nBundle:    bundle,
		}
		err = handler.Create(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, res.Code)

		body, err := ioutil.ReadAll(res.Body)
		assert.NoError(t, err)

		title, err := l.LocalizeMessage(&i18n.Message{ID: "TenantCreateSuccessTitle"})
		assert.NoError(t, err)
		msg, err := l.Localize(&i18n.LocalizeConfig{
			MessageID: "TenantCreateSuccessMsg",
			TemplateData: map[string]string{
				"name": mockTenant.Name,
			},
		})
		assert.NoError(t, err)

		assert.Contains(t, string(body), title)
		assert.Contains(t, string(body), msg)

		mockUsecase.AssertExpectations(t)
	})

	t.Run("name-exists", func(t *testing.T) {
		mockUsecase := new(mocks.Usecase)
		mockUsecase.On("Create", mock.Anything).Return(errors.New("TenantCreateNameExistsMsg")).Once()

		req, err := http.NewRequest(echo.POST, "/v1/tenants", strings.NewReader(string(json)))
		assert.NoError(t, err)

		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		res := httptest.NewRecorder()
		c := e.NewContext(req, res)
		c.Set("user", echoutil.GetLoginTkn("admin"))

		handler := _tenanthttp.TenantHandler{
			TenantUsecase: mockUsecase,
			I18nBundle:    bundle,
		}
		err = handler.Create(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusConflict, res.Code)

		body, err := ioutil.ReadAll(res.Body)
		assert.NoError(t, err)

		title, err := l.LocalizeMessage(&i18n.Message{ID: "TenantCreateFail"})
		assert.NoError(t, err)
		msg, err := l.Localize(&i18n.LocalizeConfig{
			MessageID: "TenantCreateNameExistsMsg",
			TemplateData: map[string]string{
				"name": mockTenant.Name,
			},
		})
		assert.NoError(t, err)

		assert.Contains(t, string(body), title)
		assert.Contains(t, string(body), msg)

		mockUsecase.AssertExpectations(t)
	})
}
