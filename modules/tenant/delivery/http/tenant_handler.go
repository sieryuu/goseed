package http

import (
	"goseed/models"
	"goseed/modules/tenant"
	"net/http"

	"github.com/labstack/echo/v4"
)

// TenantHandler represents httpHandler for Tenant.
type TenantHandler struct {
	TenantUsecase tenant.Usecase
}

// NewTenantHandler will initialize the tenants / resources endpoint.
func NewTenantHandler(e *echo.Group, usecase tenant.Usecase) {
	handler := &TenantHandler{
		TenantUsecase: usecase,
	}
	e.GET("/tenants", handler.Find)
	e.POST("/tenants", handler.Create)
}

// Find will find all tenants.
func (a *TenantHandler) Find(e echo.Context) error {
	tenants, _ := a.TenantUsecase.Find()
	return e.JSON(http.StatusOK, tenants)
}

// Create will insert a new tenant.
func (a *TenantHandler) Create(e echo.Context) error {
	tenant := new(models.Tenant)

	if err := e.Bind(tenant); err != nil {
		return e.JSON(http.StatusBadRequest, err)
	}

	err := a.TenantUsecase.Create(tenant)
	if err != nil {
		return e.JSON(http.StatusBadRequest, err.Error)
	}
	return e.JSON(http.StatusCreated, tenant)
}
