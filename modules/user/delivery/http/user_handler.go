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

// UserHandler represents httpHandler for user.
type UserHandler struct {
	UserUsecase user.Usecase
	I18nBundle  *i18n.Bundle
}

// NewUserHandler will initialize the user / resource endpoint.
func NewUserHandler(e *echo.Group, usecase user.Usecase, i18nBundle *i18n.Bundle) {
	handler := &UserHandler{
		UserUsecase: usecase,
		I18nBundle:  i18nBundle,
	}

	e.POST("/users", handler.Create)
}

// Create will create models.User into database.
func (a *UserHandler) Create(e echo.Context) error {
	l := echoutil.GetLocalizerFromEchoContext(a.I18nBundle, e)
	createFailTitle, _ := l.LocalizeMessage(&i18n.Message{ID: "UserCreateFailTitle"})

	userCreation := new(dto.UserCreation)

	if err := e.Bind(userCreation); err != nil {
		return e.JSON(http.StatusBadRequest, httputil.ResponseBody{
			Title: createFailTitle,
			Msg:   err.Error(),
		})
	}

	user := userCreation.GetUser()
	user.CreatedBy = echoutil.GetLoggedInUser(e)
	user.LastUpdatedBy = echoutil.GetLoggedInUser(e)

	if err := e.Validate(user); err != nil {
		return e.JSON(http.StatusBadRequest, httputil.ResponseBody{
			Title: createFailTitle,
			Msg:   err.Error(),
		})
	}

	err := a.UserUsecase.Create(user)
	if err != nil {
		statusCode := http.StatusInternalServerError
		switch err.Error() {
		case "UserCreateUsernameExistsMsg":
			statusCode = http.StatusConflict
		}

		return e.JSON(statusCode, httputil.ResponseBody{
			Title: createFailTitle,
			Msg: l.MustLocalize(&i18n.LocalizeConfig{
				DefaultMessage: &i18n.Message{
					ID:    err.Error(),
					Other: err.Error(), // system error
				},
			}),
		})
	}

	title, _ := l.LocalizeMessage(&i18n.Message{ID: "UserCreateSuccessTitle"})
	return e.JSON(http.StatusCreated, httputil.ResponseBody{
		Title: title,
		Msg: l.MustLocalize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID: "UserCreateSuccessMsg",
			},
			TemplateData: map[string]string{
				"username": user.Username,
			},
		}),
		Data: user,
	})
}
