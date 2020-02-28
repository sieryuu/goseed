package tenant

import "goseed/models"

// Repository represents the tenant's repository contract.
type Repository interface {
	Find() (*[]models.Tenant, error)
	First(id int) (*models.Tenant, error)
	Create(tenant *models.Tenant) (int64, error)
	Update(tenant *models.Tenant) (int64, error)
	Delete(tenant *models.Tenant) (int64, error)
	GetByName(title string) (*models.Tenant, error)
}
