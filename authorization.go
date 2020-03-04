package main

import (
	"goseed/modules/user/delivery/dto"
	"strings"

	"github.com/casbin/casbin/v2"
	xormadapter "github.com/casbin/xorm-adapter/v2"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"xorm.io/xorm"
)

// Enforcer represents casbin enforcer model.
type Enforcer struct {
	enforcer *casbin.Enforcer
}

// NewCasbinEnforcer will create an Enforcer.Enforce that overide default casbin Enforce func.
func NewCasbinEnforcer(engine *xorm.Engine) Enforcer {
	// Casbin
	xormAdapter, err := xormadapter.NewAdapterByEngine(engine)
	if err != nil {
		logger.Error("failed to init xormAdapter", zap.Error(err))
		panic(err)
	}
	casbinEnforcer, err := casbin.NewEnforcer("./rbac_with_tenant.conf", xormAdapter)
	if err != nil {
		logger.Error("failed to init casbinEnforcer", zap.Error(err))
		panic(err)
	}

	// example
	casbinEnforcer.DeleteRole("user")
	casbinEnforcer.AddPolicy("admin", "1", "/v1/*/articles", "GET")
	casbinEnforcer.AddPolicy("admin", "2", "/v1/*/articles", "GET")
	casbinEnforcer.AddPolicy("admin", "1", "/v1/*/articles", "POST")
	casbinEnforcer.AddPolicy("admin", "2", "/v1/*/articles", "POST")
	casbinEnforcer.AddGroupingPolicy("user", "admin", "1")
	casbinEnforcer.AddGroupingPolicy("user", "admin", "2")

	casbinEnforcer.LoadPolicy()

	return Enforcer{enforcer: casbinEnforcer}
}

// Enforce will overide enforcer Enforce method.
// https://klotzandrew.com/blog/authorization-with-casbin
func (e *Enforcer) Enforce(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(*dto.JwtCustomClaims)
		username := claims.Username

		method := c.Request().Method
		path := c.Request().URL.Path
		tenant := getTenantFromPath(path)

		result, _ := e.enforcer.Enforce(username, tenant, path, method)

		if result {
			return next(c)
		}
		return echo.ErrForbidden
	}
}

// getTenantFromPath will return tenant string from path
// example:  /v1/1/articles
//   index: 0/1 /2/3
// (second index of path is the domain)
func getTenantFromPath(path string) string {
	splits := strings.Split(path, "/")
	if len(splits) >= 2 {
		return splits[2]
	}
	return ""
}
