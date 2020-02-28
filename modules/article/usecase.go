package article

import "goseed/models"

// Usecase represents the article's usecase.
type Usecase interface {
	Find(tenantID uint) (*[]models.Article, error)
	Create(article *models.Article) error
	GetByTitle(title string) (*models.Article, error)
}
