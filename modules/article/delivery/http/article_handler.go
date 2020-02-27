package http

import (
	"goseed/models"
	"goseed/modules/article"
	"net/http"

	"github.com/labstack/echo/v4"
)

// ArticleHandler represents httpHandler for Article.
type ArticleHandler struct {
	ArticleUsecase article.Usecase
}

// NewArticleHandler will initialize the articles / resources endpoint.
func NewArticleHandler(e *echo.Group, usecase article.Usecase) {
	handler := &ArticleHandler{
		ArticleUsecase: usecase,
	}
	e.GET("/articles", handler.Find)
	e.POST("/articles", handler.Create)
}

// Find will find all articles.
func (a *ArticleHandler) Find(e echo.Context) error {
	articles, _ := a.ArticleUsecase.Find()
	return e.JSON(http.StatusOK, articles)
}

// Create will insert a new article.
func (a *ArticleHandler) Create(e echo.Context) error {
	article := new(models.Article)

	if err := e.Bind(article); err != nil {
		return e.JSON(http.StatusBadRequest, err)
	}

	err := a.ArticleUsecase.Create(article)
	if err != nil {
		return e.JSON(http.StatusBadRequest, err.Error)
	}
	return e.JSON(http.StatusCreated, article)
}
