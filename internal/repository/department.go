package repository

import (
	"context"
	"github.com/Project-Sprint-LDH-Team/GoGoManager/internal/models"
	"gorm.io/gorm"
)

type DepartmentRepository interface {
	Create(ctx context.Context, department *models.Department) error
	Update(ctx context.Context, department *models.Department) error
	Delete(ctx context.Context, id string) error
	FindByID(ctx context.Context, id string) (*models.Department, error) // Mengubah uint menjadi string
	List(ctx context.Context, userID uint, filter *models.DepartmentFilter) ([]*models.Department, error)
	HasEmployees(ctx context.Context, departmentID string) (bool, error)
}

type departmentRepository struct {
	db *gorm.DB
}

func NewDepartmentRepository(db *gorm.DB) DepartmentRepository {
	return &departmentRepository{
		db: db,
	}
}

func (r *departmentRepository) Create(ctx context.Context, department *models.Department) error {
	return r.db.WithContext(ctx).Create(department).Error
}

func (r *departmentRepository) Update(ctx context.Context, department *models.Department) error {
	return r.db.WithContext(ctx).Save(department).Error
}

func (r *departmentRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.Department{}, "id = ?", id).Error
}

func (r *departmentRepository) FindByID(ctx context.Context, id string) (*models.Department, error) {
	var department models.Department
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&department).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &department, err
}

func (r *departmentRepository) List(ctx context.Context, userID uint, filter *models.DepartmentFilter) ([]*models.Department, error) {
	var departments []*models.Department
	query := r.db.WithContext(ctx).Where("user_id = ?", userID)

	// Filter by name (prefix-suffix search, case insensitive)
	if filter.Name != "" {
		query = query.Where("name ILIKE ?", "%"+filter.Name+"%")
	}

	// Apply pagination
	query = query.Limit(filter.Limit).Offset(filter.Offset)

	// Order by id
	query = query.Order("id ASC")

	err := query.Find(&departments).Error
	return departments, err
}

func (r *departmentRepository) HasEmployees(ctx context.Context, departmentID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Employee{}).
		Where("department_id = ?", departmentID).
		Count(&count).Error
	return count > 0, err
}
