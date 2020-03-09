package http_test

import (
	"encoding/json"
	"errors"
	"goseed/models"
	"goseed/modules/article/mocks"
	"goseed/utils/echoutil"
	"goseed/utils/i18nutil"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bxcodec/faker"
	"github.com/labstack/echo/v4"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	_articlehttp "goseed/modules/article/delivery/http"
)

func TestFind(t *testing.T) {
	var mockArticle models.Article
	err := faker.FakeData(&mockArticle)
	assert.NoError(t, err)

	mockListArticle := make([]models.Article, 0)
	mockListArticle = append(mockListArticle, mockArticle)

	mockUsecase := new(mocks.Usecase)
	mockUsecase.On("Find", mock.Anything).Return(&mockListArticle, nil)

	e := echo.New()
	req, err := http.NewRequest(echo.GET, "/v1/1/articles", strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := _articlehttp.ArticleHandler{
		ArticleUsecase: mockUsecase,
	}
	err = handler.Find(c)
	// require.NoError will halt the test if found any error
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)
	mockUsecase.AssertExpectations(t)
}

func TestCreate(t *testing.T) {
	bundle := i18nutil.GetTestingi18nBundle()
	l := i18nutil.GetTestingLocalizer()

	mockArticle := models.Article{
		Title: "Article Title",
	}

	json, err := json.Marshal(mockArticle)
	assert.NoError(t, err)

	e := echoutil.MockEcho()

	t.Run("success", func(t *testing.T) {
		mockUsecase := new(mocks.Usecase)
		mockUsecase.On("Create", mock.Anything).Return(nil).Once()

		req, err := http.NewRequest(echo.POST, "/v1/1/articles", strings.NewReader(string(json)))
		assert.NoError(t, err)

		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		res := httptest.NewRecorder()
		c := e.NewContext(req, res)
		echoutil.MockTenant(c, "1")
		c.Set("user", echoutil.GetLoginTkn("admin"))

		handler := _articlehttp.ArticleHandler{
			ArticleUsecase: mockUsecase,
			I18nBundle:     bundle,
		}
		err = handler.Create(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, res.Code)

		body, err := ioutil.ReadAll(res.Body)
		assert.NoError(t, err)

		title, err := l.LocalizeMessage(&i18n.Message{ID: "ArticleCreateSuccessTitle"})
		assert.NoError(t, err)
		msg, err := l.Localize(&i18n.LocalizeConfig{
			MessageID: "ArticleCreateSuccessMsg",
			TemplateData: map[string]string{
				"title": mockArticle.Title,
			},
		})
		assert.NoError(t, err)

		assert.Contains(t, string(body), title)
		assert.Contains(t, string(body), msg)

		mockUsecase.AssertExpectations(t)
	})

	t.Run("title-exists", func(t *testing.T) {
		mockUsecase := new(mocks.Usecase)
		mockUsecase.On("Create", mock.Anything).Return(errors.New("ArticleCreateTitleExistMsg")).Once()

		req, err := http.NewRequest(echo.POST, "/v1/1/articles", strings.NewReader(string(json)))
		assert.NoError(t, err)

		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		res := httptest.NewRecorder()
		c := e.NewContext(req, res)
		echoutil.MockTenant(c, "1")
		c.Set("user", echoutil.GetLoginTkn("admin"))

		handler := _articlehttp.ArticleHandler{
			ArticleUsecase: mockUsecase,
			I18nBundle:     bundle,
		}
		err = handler.Create(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusConflict, res.Code)

		body, err := ioutil.ReadAll(res.Body)
		assert.NoError(t, err)

		title, err := l.LocalizeMessage(&i18n.Message{ID: "ArticleCreateFail"})
		assert.NoError(t, err)
		msg, err := l.LocalizeMessage(&i18n.Message{ID: "ArticleCreateTitleExistMsg"})
		assert.NoError(t, err)

		assert.Contains(t, string(body), title)
		assert.Contains(t, string(body), msg)

		mockUsecase.AssertExpectations(t)
	})
}
