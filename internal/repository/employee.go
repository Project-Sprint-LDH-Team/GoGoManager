package repository

import (
	"context"
	"github.com/Project-Sprint-LDH-Team/GoGoManager/internal/models"
	"gorm.io/gorm"
)

type EmployeeRepository interface {
	Create(ctx context.Context, employee *models.Employee) error
	Update(ctx context.Context, employee *models.Employee) error
	Delete(ctx context.Context, identityNumber string) error
	FindByIdentityNumber(ctx context.Context, identityNumber string) (*models.Employee, error)
	List(ctx context.Context, filter *models.EmployeeFilter) ([]*models.Employee, error)
	CheckIdentityExists(ctx context.Context, identityNumber string, excludeID uint) (bool, error)
}

type employeeRepository struct {
	db *gorm.DB
}

func NewEmployeeRepository(db *gorm.DB) EmployeeRepository {
	return &employeeRepository{
		db: db,
	}
}

func (r *employeeRepository) Create(ctx context.Context, employee *models.Employee) error {
	return r.db.WithContext(ctx).Create(employee).Error
}

func (r *employeeRepository) Update(ctx context.Context, employee *models.Employee) error {
	return r.db.WithContext(ctx).Save(employee).Error
}

func (r *employeeRepository) Delete(ctx context.Context, identityNumber string) error {
	return r.db.WithContext(ctx).Where("identity_number = ?", identityNumber).Delete(&models.Employee{}).Error
}

func (r *employeeRepository) FindByIdentityNumber(ctx context.Context, identityNumber string) (*models.Employee, error) {
	var employee models.Employee
	err := r.db.WithContext(ctx).Where("identity_number = ?", identityNumber).First(&employee).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &employee, err
}

func (r *employeeRepository) List(ctx context.Context, filter *models.EmployeeFilter) ([]*models.Employee, error) {
	var employees []*models.Employee
	query := r.db.WithContext(ctx)

	// Filter by identity number (prefix search)
	if filter.IdentityNumber != "" {
		query = query.Where("identity_number ILIKE ?", filter.IdentityNumber+"%")
	}

	// Filter by name (prefix-suffix search, case insensitive)
	if filter.Name != "" {
		query = query.Where("name ILIKE ?", "%"+filter.Name+"%")
	}

	// Filter by gender
	if filter.Gender != "" {
		query = query.Where("gender = ?", filter.Gender)
	}

	// Filter by department
	if filter.DepartmentID != "" {
		query = query.Where("department_id = ?", filter.DepartmentID)
	}

	// Apply pagination
	query = query.Limit(filter.Limit).Offset(filter.Offset)

	// Order by identity_number
	query = query.Order("identity_number ASC")

	err := query.Find(&employees).Error
	return employees, err
}

func (r *employeeRepository) CheckIdentityExists(ctx context.Context, identityNumber string, excludeID uint) (bool, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&models.Employee{}).Where("identity_number = ?", identityNumber)

	if excludeID != 0 {
		query = query.Where("id != ?", excludeID)
	}

	err := query.Count(&count).Error
	return count > 0, err
}
