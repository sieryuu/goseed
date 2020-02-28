package tenant

import "goseed/models"

// Usecase represents the tenant's usecase.
type Usecase interface {
	Find() (*[]models.Tenant, error)
	Create(tenant *models.Tenant) error
	GetByName(title string) (*models.Tenant, error)
}
