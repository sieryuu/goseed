package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"xorm.io/xorm"

	"goseed/models"
	_articleHttpHandler "goseed/modules/article/delivery/http"
	_articleRepo "goseed/modules/article/repository"
	_articleUsecase "goseed/modules/article/usecase"
	"goseed/utils/echoutil"

	"goseed/modules/user/delivery/dto"
	_userHttpHandler "goseed/modules/user/delivery/http"
	_userRepo "goseed/modules/user/repository"
	_userUsecase "goseed/modules/user/usecase"

	_tenantHttpHandler "goseed/modules/tenant/delivery/http"
	_tenantRepo "goseed/modules/tenant/repository"
	_tenantUsecase "goseed/modules/tenant/usecase"
)

// Router represents app
type Router struct {
	echo       *echo.Echo   // http endpoints
	db         *xorm.Engine // database connector
	enforcer   Enforcer     // casbin authorization enforcer
	i18nBundle *i18n.Bundle // i18n localization
}

// Init will configure middleware and initialize all modules endpoints.
func (r *Router) Init() {
	me1 := r.echo.Group("/me/v1")
	v1 := r.echo.Group("/v1/:" + echoutil.TenantParam)

	// Middleware
	r.echo.Use(middleware.Logger())
	r.echo.Use(middleware.Recover())
	// Jwt middleware
	config := middleware.JWTConfig{
		Claims:     &dto.JwtCustomClaims{},
		SigningKey: []byte(dto.JWTSigningKey),
	}

	// v1 middleware
	v1.Use(middleware.JWTWithConfig(config))
	v1.Use(r.enforcer.Enforce)

	// Routes
	r.echo.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	userRepo := _userRepo.NewPostgreUserRepository(r.db)
	userUsecase := _userUsecase.NewUserUsecase(userRepo)
	_userHttpHandler.NewMeHandler(me1, userUsecase, r.i18nBundle)
	_userHttpHandler.NewUserHandler(v1, userUsecase, r.i18nBundle)

	tenantRepo := _tenantRepo.NewPostgreTenantRepository(r.db)
	tenantUsecase := _tenantUsecase.NewTenantUsecase(tenantRepo)
	_tenantHttpHandler.NewTenantHandler(v1, tenantUsecase, r.i18nBundle)

	articleRepo := _articleRepo.NewPostgreArticleRepository(r.db)
	articleUsecase := _articleUsecase.NewArticleUsecase(articleRepo)
	_articleHttpHandler.NewArticleHandler(v1, articleUsecase)

	// Seeding
	userUsecase.Create(&models.User{
		Username:     "user",
		FirstName:    "John",
		LastName:     "Snow",
		PasswordHash: "$2a$04$obP4vv8R8AICEC1.Jkkibe43sa45txflPU7eeJWUaDo5JqmLPETIe", // 123qwe
	})
}
