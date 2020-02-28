package article

import "goseed/models"

// Repository represents the article's repository contract.
type Repository interface {
	Find(tenantID uint) (*[]models.Article, error)
	First(id int) (*models.Article, error)
	Create(article *models.Article) (int64, error)
	Update(article *models.Article) (int64, error)
	Delete(article *models.Article) (int64, error)
	GetByTitle(title string) (*models.Article, error)
}
