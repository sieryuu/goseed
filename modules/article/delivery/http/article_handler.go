package http

import (
	"goseed/models"
	"goseed/modules/article"
	"goseed/utils/echoutil"
	"goseed/utils/httputil"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

// ArticleHandler represents httpHandler for Article.
type ArticleHandler struct {
	ArticleUsecase article.Usecase
	I18nBundle     *i18n.Bundle
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
	tenantID := echoutil.GetTenant(e)

	articles, _ := a.ArticleUsecase.Find(tenantID)
	return e.JSON(http.StatusOK, httputil.ResponseBody{
		Data: articles,
	})
}

// Create will insert a new article.
func (a *ArticleHandler) Create(e echo.Context) error {
	l := echoutil.GetLocalizerFromEchoContext(a.I18nBundle, e)

	article := new(models.Article)

	createFailTitle, _ := l.LocalizeMessage(&i18n.Message{ID: "ArticleCreateFail"})
	if err := e.Bind(article); err != nil {
		return e.JSON(http.StatusBadRequest, httputil.ResponseBody{
			Title: createFailTitle,
			Msg:   err.Error(),
		})
	}

	article.TenantID = echoutil.GetTenant(e)
	article.CreatedBy = echoutil.GetLoggedInUser(e)
	article.LastUpdatedBy = echoutil.GetLoggedInUser(e)

	if err := e.Validate(article); err != nil {
		return e.JSON(http.StatusBadRequest, httputil.ResponseBody{
			Title: createFailTitle,
			Msg:   err.Error(),
		})
	}

	err := a.ArticleUsecase.Create(article)
	if err != nil {
		statusCode := http.StatusInternalServerError
		switch err.Error() {
		case "ArticleCreateTitleExistMsg":
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

	createSuccessTitle, _ := l.LocalizeMessage(&i18n.Message{ID: "ArticleCreateSuccessTitle"})
	return e.JSON(http.StatusCreated, httputil.ResponseBody{
		Title: createSuccessTitle,
		Msg: l.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "ArticleCreateSuccessMsg",
			TemplateData: map[string]string{
				"title": article.Title,
			},
		}),
		Data: article,
	})
}
