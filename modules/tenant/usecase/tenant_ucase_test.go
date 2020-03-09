package usecase_test

import (
	"goseed/models"
	"goseed/modules/tenant/mocks"
	"goseed/modules/tenant/usecase"
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreate(t *testing.T) {
	mockTenant := models.Tenant{
		Name: "Tenant #1",
	}

	t.Run("success", func(t *testing.T) {
		mockTenantRepo := new(mocks.Repository)
		mockTenantRepo.On("Create", mock.Anything).Return(int64(1), nil).Once()
		mockTenantRepo.On("GetByName", mockTenant.Name).Return(nil, nil).Once()

		usecase := usecase.NewTenantUsecase(mockTenantRepo)
		err := usecase.Create(&mockTenant)

		assert.NoError(t, err)
	})

	t.Run("existing-title", func(t *testing.T) {
		mockTenantRepo := new(mocks.Repository)
		mockTenantRepo.On("Create", mock.Anything).Return(int64(1), nil).Once()
		mockTenantRepo.On("GetByName", mockTenant.Name).Return(&mockTenant, nil).Once()

		usecase := usecase.NewTenantUsecase(mockTenantRepo)
		err := usecase.Create(&mockTenant)

		assert.Error(t, err)
		assert.Equal(t, err.Error(), "TenantCreateNameExistsMsg")
	})
}

func TestFind(t *testing.T) {
	mockListTenants := make([]models.Tenant, 0)
	var mockTenant models.Tenant
	faker.FakeData(&mockTenant)
	mockListTenants = append(mockListTenants, mockTenant)

	t.Run("success", func(t *testing.T) {
		mockTenantRepo := new(mocks.Repository)
		mockTenantRepo.On("Find").Return(&mockListTenants, nil).Once()

		usecase := usecase.NewTenantUsecase(mockTenantRepo)
		tenants, err := usecase.Find()

		assert.NoError(t, err)
		assert.NotNil(t, tenants)
	})
}
