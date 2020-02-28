package echohlpr

import (
	"goseed/modules/user/delivery/dto"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

const (
	// TenantParam represents "tenant" param path in url
	TenantParam = "tenant"
)

// GetTenant will return tenant string from url param
func GetTenant(c echo.Context) uint {
	tenantID, err := strconv.ParseUint(c.Param(TenantParam), 10, 32)
	if err != nil {
		return 0
	}
	return uint(tenantID)
}

// GetUser will return logged in username
func GetUser(c echo.Context) string {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*dto.JwtCustomClaims)
	return claims.Username
}
