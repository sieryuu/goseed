package http

import (
	"goseed/modules/user"
	"goseed/modules/user/delivery/dto"
	"goseed/utils/echoutil"
	"goseed/utils/httputil"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

// MeHandler represents httpHandler for Me
type MeHandler struct {
	UserUsecase user.Usecase
	I18nBundle  *i18n.Bundle
}

// NewMeHandler will initialize the user(me) / resource endpoint
func NewMeHandler(e *echo.Group, usecase user.Usecase, i18nBundle *i18n.Bundle) {
	handler := &MeHandler{
		UserUsecase: usecase,
		I18nBundle:  i18nBundle,
	}

	e.POST("/login", handler.Login)
}

// Login will authenticate user & password given.
func (a *MeHandler) Login(e echo.Context) error {
	l := echoutil.GetLocalizerFromEchoContext(a.I18nBundle, e)
	user := new(dto.UserLogin)

	loginFailTitle, _ := l.LocalizeMessage(&i18n.Message{ID: "MeLoginInvalidTitle"})

	if err := e.Bind(user); err != nil {
		return e.JSON(http.StatusBadRequest, httputil.ResponseBody{
			Title: loginFailTitle,
			Msg:   err.Error(),
		})
	}

	res, err := a.UserUsecase.Login(user.Username, []byte(user.Password))
	if err != nil {
		return e.JSON(http.StatusInternalServerError, httputil.ResponseBody{
			Title: loginFailTitle,
			Msg:   err.Error(),
		})
	}
	if res == false {
		msg, _ := l.LocalizeMessage(&i18n.Message{ID: "MeLoginInvalidMsg"})
		return e.JSON(http.StatusBadRequest, httputil.ResponseBody{
			Title: loginFailTitle,
			Msg:   msg,
		})
	}

	// generate encoded token and send it as response
	t, err := echoutil.GenLoginTknString(user.Username)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, httputil.ResponseBody{
			Title: loginFailTitle,
			Msg:   err.Error(),
		})
	}

	title, _ := l.LocalizeMessage(&i18n.Message{ID: "MeLoginSuccessTitle"})
	msg, _ := l.Localize(&i18n.LocalizeConfig{
		MessageID: "MeLoginSuccessMsg",
		TemplateData: map[string]string{
			"username": user.Username,
		},
	})
	return e.JSON(http.StatusOK, httputil.ResponseBody{
		Title: title,
		Msg:   msg,
		Data:  t, // token
	})
}
