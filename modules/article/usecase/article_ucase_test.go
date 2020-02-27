package usecase_test

import (
	"goseed/models"
	"goseed/modules/article/mocks"
	"goseed/modules/article/usecase"
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreate(t *testing.T) {
	mockArticleRepo := new(mocks.Repository)
	mockArticle := models.Article{
		Title: "Article #1",
	}

	t.Run("success", func(t *testing.T) {
		mockArticleRepo.On("Create", mock.Anything).Return(int64(1), nil).Once()
		mockArticleRepo.On("GetByTitle", mockArticle.Title).Return(nil, nil).Once()

		usecase := usecase.NewArticleUsecase(mockArticleRepo)
		err := usecase.Create(&mockArticle)

		assert.NoError(t, err)
	})

	t.Run("existing-title", func(t *testing.T) {
		mockArticleRepo.On("Create", mock.Anything).Return(int64(1), nil).Once()
		mockArticleRepo.On("GetByTitle", mockArticle.Title).Return(&mockArticle, nil).Once()

		usecase := usecase.NewArticleUsecase(mockArticleRepo)
		err := usecase.Create(&mockArticle)

		assert.Error(t, err)
		assert.Equal(t, err.Error(), "title already exists")
	})
}

func TestFind(t *testing.T) {
	mockArticleRepo := new(mocks.Repository)

	mockListArticles := make([]models.Article, 0)
	var mockArticle models.Article
	faker.FakeData(&mockArticle)
	mockListArticles = append(mockListArticles, mockArticle)

	t.Run("success", func(t *testing.T) {
		mockArticleRepo.On("Find").Return(&mockListArticles, nil).Once()

		usecase := usecase.NewArticleUsecase(mockArticleRepo)
		articles, err := usecase.Find()

		assert.NoError(t, err)
		assert.NotNil(t, articles)
	})
}
