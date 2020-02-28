package usecase

import (
	"errors"
	"goseed/models"
	"goseed/modules/tenant"
)

type tenantUsecase struct {
	tenantRepo tenant.Repository
}

// NewTenantUsecase will create a new tenantUsecase object representation of tenant.Usecase interface.
func NewTenantUsecase(a tenant.Repository) tenant.Usecase {
	return &tenantUsecase{
		tenantRepo: a,
	}
}

func (a *tenantUsecase) Find() (*[]models.Tenant, error) {
	return a.tenantRepo.Find()
}

func (a *tenantUsecase) Create(tenant *models.Tenant) error {
	// below 'tenant' variable only exists inside if scope
	if tenant, _ := a.tenantRepo.GetByName(tenant.Name); tenant != nil {
		return errors.New("title already exists")
	}

	// below tenant is taken from param
	_, err := a.tenantRepo.Create(tenant)

	return err
}

func (a *tenantUsecase) GetByName(name string) (*models.Tenant, error) {
	return a.tenantRepo.GetByName(name)
}
