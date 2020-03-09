package echoutil

import (
	"goseed/modules/user/delivery/dto"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

const (
	// TenantParam represents "tenant" param path in url
	TenantParam = "tenant"
)

type (
	// CustomValidator represents struct to overide echo validator
	CustomValidator struct {
		Validator *validator.Validate
	}
)

// Validate will overide echo super 'Validate' func
func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.Validator.Struct(i)
}

// GetTenant will return tenant string from url param
func GetTenant(c echo.Context) uint {
	tenantID, err := strconv.ParseUint(c.Param(TenantParam), 10, 32)
	if err != nil {
		return 0
	}
	return uint(tenantID)
}

// GetLoggedInUser will return logged in username
func GetLoggedInUser(c echo.Context) string {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*dto.JwtCustomClaims)
	return claims.Username
}

// GetLocalizerFromEchoContext will return *i18n.Localizer depends on http header 'Accept-Language'
func GetLocalizerFromEchoContext(b *i18n.Bundle, e echo.Context) *i18n.Localizer {
	accept := e.Request().Header.Get("Accept-Language")
	return i18n.NewLocalizer(b, accept)
}

// GenLoginTknString will return token string by given params
func GenLoginTknString(username string) (string, error) {
	token := GetLoginTkn(username)
	// generate encoded token and send it as response
	t, err := token.SignedString([]byte(dto.JWTSigningKey))
	return t, err
}

// GetLoginTkn will return *jwt.token by given params
func GetLoginTkn(username string) *jwt.Token {
	// set claims
	claims := &dto.JwtCustomClaims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
}

// Mocking Code below

// MockEcho will return echo.Echo pointer for Unit Testing Purpose
func MockEcho() *echo.Echo {
	e := echo.New()
	e.Validator = &CustomValidator{Validator: validator.New()}
	return e
}

// MockTenant will insert tenantID into echo.Context for testing purpose
func MockTenant(c echo.Context, tenantID string) {
	c.SetParamNames(TenantParam)
	c.SetParamValues(tenantID)
}
