package service

import (
	"context"
	"errors"
	"github.com/Project-Sprint-LDH-Team/GoGoManager/internal/models"
	"github.com/Project-Sprint-LDH-Team/GoGoManager/internal/repository"
)

type EmployeeService interface {
	CreateEmployee(ctx context.Context, req *models.CreateEmployeeRequest) (*models.EmployeeResponse, error)
	UpdateEmployee(ctx context.Context, identityNumber string, req *models.UpdateEmployeeRequest) (*models.EmployeeResponse, error)
	DeleteEmployee(ctx context.Context, identityNumber string) error
	ListEmployees(ctx context.Context, filter *models.EmployeeFilter) ([]*models.EmployeeResponse, error)
}

type employeeService struct {
	employeeRepo   repository.EmployeeRepository
	departmentRepo repository.DepartmentRepository
}

func NewEmployeeService(employeeRepo repository.EmployeeRepository, departmentRepo repository.DepartmentRepository) EmployeeService {
	return &employeeService{
		employeeRepo:   employeeRepo,
		departmentRepo: departmentRepo,
	}
}

func (s *employeeService) CreateEmployee(ctx context.Context, req *models.CreateEmployeeRequest) (*models.EmployeeResponse, error) {
	// Check if department exists
	dept, err := s.departmentRepo.FindByID(ctx, req.DepartmentId)
	if err != nil {
		return nil, err
	}
	if dept == nil {
		return nil, errors.New("department not found")
	}

	// Check if identity number exists
	existingEmp, err := s.employeeRepo.FindByIdentityNumber(ctx, req.IdentityNumber)
	if err != nil {
		return nil, err
	}
	if existingEmp != nil {
		return nil, errors.New("identity number already exists")
	}

	// Create employee
	employee := &models.Employee{
		IdentityNumber:   req.IdentityNumber,
		Name:             req.Name,
		EmployeeImageUri: req.EmployeeImageUri,
		Gender:           req.Gender,
		DepartmentID:     req.DepartmentId,
	}

	if err := s.employeeRepo.Create(ctx, employee); err != nil {
		return nil, err
	}

	return employee.ToResponse(), nil
}

func (s *employeeService) UpdateEmployee(ctx context.Context, identityNumber string, req *models.UpdateEmployeeRequest) (*models.EmployeeResponse, error) {
	// Check if employee exists
	employee, err := s.employeeRepo.FindByIdentityNumber(ctx, identityNumber)
	if err != nil {
		return nil, err
	}
	if employee == nil {
		return nil, errors.New("employee not found")
	}

	// Check if department exists
	dept, err := s.departmentRepo.FindByID(ctx, req.DepartmentId)
	if err != nil {
		return nil, err
	}
	if dept == nil {
		return nil, errors.New("department not found")
	}

	// Check if new identity number exists (if changed)
	if req.IdentityNumber != identityNumber {
		existingEmp, err := s.employeeRepo.FindByIdentityNumber(ctx, req.IdentityNumber)
		if err != nil {
			return nil, err
		}
		if existingEmp != nil {
			return nil, errors.New("identity number already exists")
		}
	}

	// Update employee
	employee.IdentityNumber = req.IdentityNumber
	employee.Name = req.Name
	employee.EmployeeImageUri = req.EmployeeImageUri
	employee.Gender = req.Gender
	employee.DepartmentID = req.DepartmentId

	if err := s.employeeRepo.Update(ctx, employee); err != nil {
		return nil, err
	}

	return employee.ToResponse(), nil
}

func (s *employeeService) DeleteEmployee(ctx context.Context, identityNumber string) error {
	employee, err := s.employeeRepo.FindByIdentityNumber(ctx, identityNumber)
	if err != nil {
		return err
	}
	if employee == nil {
		return errors.New("employee not found")
	}

	return s.employeeRepo.Delete(ctx, identityNumber)
}

func (s *employeeService) ListEmployees(ctx context.Context, filter *models.EmployeeFilter) ([]*models.EmployeeResponse, error) {
	employees, err := s.employeeRepo.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	var response []*models.EmployeeResponse
	for _, emp := range employees {
		response = append(response, emp.ToResponse())
	}

	return response, nil
}
