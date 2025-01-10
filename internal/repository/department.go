package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/Project-Sprint-LDH-Team/GoGoManager/internal/models"
	"gorm.io/gorm"
)

type DepartmentRepository interface {
	Create(ctx context.Context, department *models.Department) error
	Update(ctx context.Context, department *models.Department) error
	Delete(ctx context.Context, departmentID string) error
	FindByDepartmentID(ctx context.Context, departmentID string) (*models.Department, error)
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
	// Cek apakah department_id sudah ada (termasuk yang soft deleted)
	var count int64
	if err := r.db.WithContext(ctx).Unscoped().
		Model(&models.Department{}).
		Where("department_id = ?", department.DepartmentID).
		Count(&count).Error; err != nil {
		return err
	}

	// Jika ditemukan, update ID increment untuk mencegah duplikasi
	if count > 0 {
		// Ambil ID terakhir
		var lastDepartment models.Department
		if err := r.db.WithContext(ctx).Unscoped().
			Order("id desc").
			First(&lastDepartment).Error; err != nil && err != gorm.ErrRecordNotFound {
			return err
		}

		// Set ID baru
		department.DepartmentID = fmt.Sprintf("DEP-%02d", lastDepartment.ID+1)
	}

	// Create department baru
	return r.db.WithContext(ctx).Create(department).Error
}

func (r *departmentRepository) Update(ctx context.Context, department *models.Department) error {
	return r.db.WithContext(ctx).Save(department).Error
}

func (r *departmentRepository) Delete(ctx context.Context, departmentID string) error {
	// Gunakan Unscoped() jika ingin melihat semua data termasuk yang soft deleted
	result := r.db.WithContext(ctx).
		Where("department_id = ?", departmentID).
		Delete(&models.Department{}) // GORM akan otomatis mengisi deleted_at

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("department not found")
	}

	return nil
}

func (r *departmentRepository) FindByDepartmentID(ctx context.Context, departmentID string) (*models.Department, error) {
	var department models.Department

	// Tambahkan Unscoped() jika ingin melihat semua data termasuk yang soft deleted
	err := r.db.WithContext(ctx).
		Where("department_id = ? AND deleted_at IS NULL", departmentID). // Tambahkan pengecekan deleted_at
		First(&department).Error

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &department, err
}

func (r *departmentRepository) List(ctx context.Context, userID uint, filter *models.DepartmentFilter) ([]*models.Department, error) {
	var departments []*models.Department
	query := r.db.WithContext(ctx).Where("user_id = ? AND deleted_at IS NULL", userID)

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
