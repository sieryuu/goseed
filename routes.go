package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"xorm.io/xorm"

	"goseed/models"
	_articleHttpHandler "goseed/modules/article/delivery/http"
	_articleRepo "goseed/modules/article/repository"
	_articleUsecase "goseed/modules/article/usecase"

	"goseed/modules/user/delivery/dto"
	_userHttpHandler "goseed/modules/user/delivery/http"
	_userRepo "goseed/modules/user/repository"
	_userUsecase "goseed/modules/user/usecase"
)

// ConfigureEcho will configure middleware and initialize all modules endpoints.
func ConfigureEcho(e *echo.Echo, db *xorm.Engine, enforcer Enforcer) {
	me := e.Group("/me")
	v1 := e.Group("/v1/:domain")

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	// Jwt middleware
	config := middleware.JWTConfig{
		Claims:     &dto.JwtCustomClaims{},
		SigningKey: []byte(dto.JWTSigningKey),
	}

	// v1 middleware
	v1.Use(middleware.JWTWithConfig(config))
	v1.Use(enforcer.Enforce)

	// Routes
	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	db.Sync(new(models.User))
	db.Sync(new(models.Article))

	userRepo := _userRepo.NewPostgreUserRepository(db)
	userUsecase := _userUsecase.NewUserUsecase(userRepo)
	_userHttpHandler.NewUserHandler(me, userUsecase)

	articleRepo := _articleRepo.NewPostgreArticleRepository(db)
	articleUsecase := _articleUsecase.NewArticleUsecase(articleRepo)
	_articleHttpHandler.NewArticleHandler(v1, articleUsecase)

	// Seeding
	userUsecase.Create(&dto.UserCreation{
		Username:  "admin",
		FirstName: "John",
		LastName:  "Snow",
		Password:  "123qwe",
	})
}
