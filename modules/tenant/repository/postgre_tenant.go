package repository

import (
	"goseed/models"
	"goseed/modules/tenant"

	"xorm.io/xorm"
)

type postgreTenantRepository struct {
	Conn *xorm.Engine
}

// NewPostgreTenantRepository will create an object that represents the tenant.Repository interface.
func NewPostgreTenantRepository(conn *xorm.Engine) tenant.Repository {
	return &postgreTenantRepository{conn}
}

func (m *postgreTenantRepository) Find() (*[]models.Tenant, error) {
	tenants := new([]models.Tenant)
	err := m.Conn.Find(tenants)
	return tenants, err
}

func (m *postgreTenantRepository) First(id int) (*models.Tenant, error) {
	tenant := new(models.Tenant)
	_, err := m.Conn.Where("id = ?", id).Get(tenant)
	return tenant, err
}

func (m *postgreTenantRepository) Create(tenant *models.Tenant) (int64, error) {
	return m.Conn.Insert(tenant)
}

func (m *postgreTenantRepository) Update(tenant *models.Tenant) (int64, error) {
	return m.Conn.ID(tenant.ID).Update(tenant)
}

func (m *postgreTenantRepository) Delete(tenant *models.Tenant) (int64, error) {
	return m.Conn.ID(tenant.ID).Delete(tenant)
}

func (m *postgreTenantRepository) GetByName(name string) (*models.Tenant, error) {
	tenant := new(models.Tenant)
	exists, err := m.Conn.Where("name = ?", name).Get(tenant)
	if exists {
		return tenant, err
	}

	return nil, err
}
