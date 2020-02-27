package http

import (
	"goseed/modules/user"
	"goseed/modules/user/delivery/dto"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

// UserHandler represents httpHandler for user.
type UserHandler struct {
	UserUsecase user.Usecase
}

// NewUserHandler will initialize the user / resource endpoint.
func NewUserHandler(e *echo.Group, usecase user.Usecase) {
	handler := &UserHandler{
		UserUsecase: usecase,
	}

	e.POST("/login", handler.Login)
}

// Login will authenticate user & password given.
func (a *UserHandler) Login(e echo.Context) error {
	user := new(dto.UserLogin)

	if err := e.Bind(user); err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}

	res, err := a.UserUsecase.Login(user.Username, []byte(user.Password))
	if err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}
	if res == false {
		return e.JSON(http.StatusBadRequest, "")
	}

	// set claims
	claims := &dto.JwtCustomClaims{
		user.Username,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}

	// create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// generate encoded token and send it as response
	t, err := token.SignedString([]byte(dto.JWTSigningKey))
	if err != nil {
		return e.JSON(http.StatusInternalServerError, err.Error)
	}

	return e.JSON(http.StatusOK, map[string]string{
		"token": t,
	})
}

// Create will create models.User into database.
func (a *UserHandler) Create(e echo.Context) error {
	userCreation := new(dto.UserCreation)

	if err := e.Bind(userCreation); err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}

	user, err := a.UserUsecase.Create(userCreation)
	if err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}

	return e.JSON(http.StatusCreated, user)
}
