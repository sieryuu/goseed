package http

import (
	"goseed/models"
	"goseed/modules/tenant"
	"goseed/utils/echoutil"
	"goseed/utils/httputil"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

// TenantHandler represents httpHandler for Tenant.
type TenantHandler struct {
	TenantUsecase tenant.Usecase
	I18nBundle    *i18n.Bundle
}

// NewTenantHandler will initialize the tenants / resources endpoint.
func NewTenantHandler(e *echo.Group, usecase tenant.Usecase, bundle *i18n.Bundle) {
	handler := &TenantHandler{
		TenantUsecase: usecase,
		I18nBundle:    bundle,
	}
	e.GET("/tenants", handler.Find)
	e.POST("/tenants", handler.Create)
}

// Find will find all tenants.
func (a *TenantHandler) Find(e echo.Context) error {
	tenants, _ := a.TenantUsecase.Find()
	return e.JSON(http.StatusOK, httputil.ResponseBody{
		Data: tenants,
	})
}

// Create will insert a new tenant.
func (a *TenantHandler) Create(e echo.Context) error {
	l := echoutil.GetLocalizerFromEchoContext(a.I18nBundle, e)
	tenant := new(models.Tenant)

	createFailTitle, _ := l.LocalizeMessage(&i18n.Message{ID: "TenantCreateFail"})
	if err := e.Bind(tenant); err != nil {
		return e.JSON(http.StatusBadRequest, httputil.ResponseBody{
			Title: createFailTitle,
			Msg:   err.Error(),
		})
	}

	tenant.CreatedBy = echoutil.GetLoggedInUser(e)
	tenant.LastUpdatedBy = echoutil.GetLoggedInUser(e)
	if err := e.Validate(tenant); err != nil {
		return e.JSON(http.StatusBadRequest, httputil.ResponseBody{
			Title: createFailTitle,
			Msg:   err.Error(),
		})
	}

	err := a.TenantUsecase.Create(tenant)
	if err != nil {
		statusCode := http.StatusInternalServerError
		switch err.Error() {
		case "TenantCreateNameExistsMsg":
			statusCode = http.StatusConflict
		}

		return e.JSON(statusCode, httputil.ResponseBody{
			Title: createFailTitle,
			Msg: l.MustLocalize(&i18n.LocalizeConfig{
				DefaultMessage: &i18n.Message{
					ID:    err.Error(),
					Other: err.Error(),
				},
			}),
		})
	}

	title, _ := l.LocalizeMessage(&i18n.Message{ID: "TenantCreateSuccessTitle"})
	msg, _ := l.Localize(&i18n.LocalizeConfig{
		MessageID: "TenantCreateSuccessMsg",
		TemplateData: map[string]string{
			"name": tenant.Name,
		},
	})
	return e.JSON(http.StatusCreated, httputil.ResponseBody{
		Title: title,
		Msg:   msg,
		Data:  tenant,
	})
}
