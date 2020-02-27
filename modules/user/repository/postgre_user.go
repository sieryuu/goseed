package repostory

import (
	"goseed/models"
	"goseed/modules/user"

	"xorm.io/xorm"
)

type postgreUserRepository struct {
	Conn *xorm.Engine
}

// NewPostgreUserRepository will create an object that represents user.Repository interface.
func NewPostgreUserRepository(conn *xorm.Engine) user.Repository {
	return &postgreUserRepository{conn}
}

func (m *postgreUserRepository) Find() (*[]models.User, error) {
	users := new([]models.User)
	err := m.Conn.Find(users)
	return users, err
}

func (m *postgreUserRepository) Create(user *models.User) (int64, error) {
	return m.Conn.InsertOne(user)
}

func (m *postgreUserRepository) GetByUsername(username string) (*models.User, error) {
	user := new(models.User)
	exists, err := m.Conn.Where("username = ?", username).Get(user)
	if exists {
		return user, err
	}

	return nil, err
}
