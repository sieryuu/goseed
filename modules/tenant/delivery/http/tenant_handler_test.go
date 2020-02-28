package http_test

import (
	"encoding/json"
	"goseed/models"
	"goseed/modules/tenant/mocks"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bxcodec/faker"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

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
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)
	mockUsecase.AssertExpectations(t)
}

func TestCreate(t *testing.T) {
	var mockTenant models.Tenant
	err := faker.FakeData(&mockTenant)
	assert.NoError(t, err)

	mockUsecase := new(mocks.Usecase)
	mockUsecase.On("Create", mock.Anything).Return(nil)

	json, err := json.Marshal(mockTenant)
	assert.NoError(t, err)

	e := echo.New()
	req, err := http.NewRequest(echo.POST, "/v1/tenants", strings.NewReader(string(json)))
	assert.NoError(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := _tenanthttp.TenantHandler{
		TenantUsecase: mockUsecase,
	}
	err = handler.Create(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)
	mockUsecase.AssertExpectations(t)
}
