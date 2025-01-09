package handlers

import (
	"fmt"
	"github.com/Project-Sprint-LDH-Team/GoGoManager/internal/models"
	"github.com/Project-Sprint-LDH-Team/GoGoManager/internal/service"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type EmployeeHandler struct {
	employeeService service.EmployeeService
}

func NewEmployeeHandler(employeeService service.EmployeeService) *EmployeeHandler {
	return &EmployeeHandler{
		employeeService: employeeService,
	}
}

func (h *EmployeeHandler) CreateEmployee(c *fiber.Ctx) error {
	var req models.CreateEmployeeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate request
	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": formatValidationErrors(err),
		})
	}

	response, err := h.employeeService.CreateEmployee(c.Context(), &req)
	if err != nil {
		switch err.Error() {
		case "department not found":
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		case "identity number already exists":
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": err.Error(),
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Internal server error",
			})
		}
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}

func (h *EmployeeHandler) UpdateEmployee(c *fiber.Ctx) error {
	identityNumber := c.Params("identityNumber")
	if identityNumber == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Identity number is required",
		})
	}

	var req models.UpdateEmployeeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": formatValidationErrors(err),
		})
	}

	response, err := h.employeeService.UpdateEmployee(c.Context(), identityNumber, &req)
	if err != nil {
		switch err.Error() {
		case "employee not found":
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		case "department not found":
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		case "identity number already exists":
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": err.Error(),
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Internal server error",
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *EmployeeHandler) DeleteEmployee(c *fiber.Ctx) error {
	identityNumber := c.Params("identityNumber")
	if identityNumber == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Identity number is required",
		})
	}

	err := h.employeeService.DeleteEmployee(c.Context(), identityNumber)
	if err != nil {
		if err.Error() == "employee not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	return c.SendStatus(fiber.StatusOK)
}

func (h *EmployeeHandler) ListEmployees(c *fiber.Ctx) error {
	filter := &models.EmployeeFilter{
		Limit:  5, // default limit
		Offset: 0, // default offset
	}

	// Parse query parameters
	if limit := c.Query("limit"); limit != "" {
		if val, err := strconv.Atoi(limit); err == nil && val > 0 {
			filter.Limit = val
		}
	}

	if offset := c.Query("offset"); offset != "" {
		if val, err := strconv.Atoi(offset); err == nil && val >= 0 {
			filter.Offset = val
		}
	}

	// Optional filters
	filter.IdentityNumber = c.Query("identityNumber")
	filter.Name = c.Query("name")
	filter.Gender = c.Query("gender")
	filter.DepartmentID = c.Query("departmentId") // Langsung assign string departmentId

	employees, err := h.employeeService.ListEmployees(c.Context(), filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(employees)
}

// Helper function untuk format validation errors
func formatValidationErrors(err error) []map[string]string {
	var errors []map[string]string
	for _, err := range err.(validator.ValidationErrors) {
		errors = append(errors, map[string]string{
			"field":   err.Field(),
			"message": formatValidationMessage(err),
		})
	}
	return errors
}

func formatValidationMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", err.Field())
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", err.Field(), err.Param())
	case "max":
		return fmt.Sprintf("%s must not exceed %s characters", err.Field(), err.Param())
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", err.Field(), err.Param())
	case "uri":
		return fmt.Sprintf("%s must be a valid URI", err.Field())
	default:
		return fmt.Sprintf("%s is invalid", err.Field())
	}
}
