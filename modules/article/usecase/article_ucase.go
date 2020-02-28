package usecase

import (
	"errors"
	"goseed/models"
	"goseed/modules/article"
)

type articleUsecase struct {
	articleRepo article.Repository
}

// NewArticleUsecase will create a new articleUsecase object representation of article.Usecase interface.
func NewArticleUsecase(a article.Repository) article.Usecase {
	return &articleUsecase{
		articleRepo: a,
	}
}

func (a *articleUsecase) Find(tenantID uint) (*[]models.Article, error) {
	return a.articleRepo.Find(tenantID)
}

func (a *articleUsecase) Create(article *models.Article) error {
	// below 'article' variable only exists inside if scope
	if article, _ := a.articleRepo.GetByTitle(article.Title); article != nil {
		return errors.New("title already exists")
	}

	// below article is taken from param
	_, err := a.articleRepo.Create(article)

	return err
}

func (a *articleUsecase) GetByTitle(title string) (*models.Article, error) {
	return a.articleRepo.GetByTitle(title)
}
