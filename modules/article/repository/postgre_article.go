package repository

import (
	"goseed/models"
	"goseed/modules/article"

	"xorm.io/xorm"
)

type postgreArticleRepository struct {
	Conn *xorm.Engine
}

// NewPostgreArticleRepository will create an object that represents the article.Repository interface.
func NewPostgreArticleRepository(conn *xorm.Engine) article.Repository {
	return &postgreArticleRepository{conn}
}

func (m *postgreArticleRepository) Find() (*[]models.Article, error) {
	articles := new([]models.Article)
	err := m.Conn.Find(articles)
	return articles, err
}

func (m *postgreArticleRepository) First(id int) (*models.Article, error) {
	article := new(models.Article)
	_, err := m.Conn.Where("id = ?", id).Get(article)
	return article, err
}

func (m *postgreArticleRepository) Create(article *models.Article) (int64, error) {
	return m.Conn.Insert(article)
}

func (m *postgreArticleRepository) Update(article *models.Article) (int64, error) {
	return m.Conn.ID(article.ID).Update(article)
}

func (m *postgreArticleRepository) Delete(article *models.Article) (int64, error) {
	return m.Conn.ID(article.ID).Delete(article)
}

func (m *postgreArticleRepository) GetByTitle(title string) (*models.Article, error) {
	article := new(models.Article)
	exists, err := m.Conn.Where("title = ?", title).Get(article)
	if exists {
		return article, err
	}

	return nil, err
}
