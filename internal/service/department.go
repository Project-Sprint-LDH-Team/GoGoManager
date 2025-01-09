package service

import (
	"context"
	"errors"
	"github.com/Project-Sprint-LDH-Team/GoGoManager/internal/models"
	"github.com/Project-Sprint-LDH-Team/GoGoManager/internal/repository"
	"github.com/google/uuid"
)

type DepartmentService interface {
	CreateDepartment(ctx context.Context, userID uint, req *models.CreateDepartmentRequest) (*models.DepartmentResponse, error)
	UpdateDepartment(ctx context.Context, userID uint, departmentID string, req *models.UpdateDepartmentRequest) (*models.DepartmentResponse, error)
	DeleteDepartment(ctx context.Context, userID uint, departmentID string) error
	ListDepartments(ctx context.Context, userID uint, filter *models.DepartmentFilter) ([]*models.DepartmentResponse, error)
}

type departmentService struct {
	departmentRepo repository.DepartmentRepository
}

func NewDepartmentService(departmentRepo repository.DepartmentRepository) DepartmentService {
	return &departmentService{
		departmentRepo: departmentRepo,
	}
}

func (s *departmentService) CreateDepartment(ctx context.Context, userID uint, req *models.CreateDepartmentRequest) (*models.DepartmentResponse, error) {
	department := &models.Department{
		ID:     uuid.New().String(), // Generate UUID for department ID
		UserID: userID,
		Name:   req.Name,
	}

	if err := s.departmentRepo.Create(ctx, department); err != nil {
		return nil, err
	}

	return department.ToResponse(), nil
}

func (s *departmentService) UpdateDepartment(ctx context.Context, userID uint, departmentID string, req *models.UpdateDepartmentRequest) (*models.DepartmentResponse, error) {
	department, err := s.departmentRepo.FindByID(ctx, departmentID)
	if err != nil {
		return nil, err
	}
	if department == nil {
		return nil, errors.New("department not found")
	}

	// Verify ownership
	if department.UserID != userID {
		return nil, errors.New("unauthorized access to department")
	}

	department.Name = req.Name

	if err := s.departmentRepo.Update(ctx, department); err != nil {
		return nil, err
	}

	return department.ToResponse(), nil
}

func (s *departmentService) DeleteDepartment(ctx context.Context, userID uint, departmentID string) error {
	department, err := s.departmentRepo.FindByID(ctx, departmentID)
	if err != nil {
		return err
	}
	if department == nil {
		return errors.New("department not found")
	}

	// Verify ownership
	if department.UserID != userID {
		return errors.New("unauthorized access to department")
	}

	// Check if department has employees
	hasEmployees, err := s.departmentRepo.HasEmployees(ctx, departmentID)
	if err != nil {
		return err
	}
	if hasEmployees {
		return errors.New("department still contains employees")
	}

	return s.departmentRepo.Delete(ctx, departmentID)
}

func (s *departmentService) ListDepartments(ctx context.Context, userID uint, filter *models.DepartmentFilter) ([]*models.DepartmentResponse, error) {
	filter.Normalize()

	departments, err := s.departmentRepo.List(ctx, userID, filter)
	if err != nil {
		return nil, err
	}

	var response []*models.DepartmentResponse
	for _, dept := range departments {
		response = append(response, dept.ToResponse())
	}

	return response, nil
}
