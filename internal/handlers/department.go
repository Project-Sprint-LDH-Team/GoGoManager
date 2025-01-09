package handlers

import (
	"github.com/Project-Sprint-LDH-Team/GoGoManager/internal/models"
	"github.com/Project-Sprint-LDH-Team/GoGoManager/internal/service"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type DepartmentHandler struct {
	departmentService service.DepartmentService
}

func NewDepartmentHandler(departmentService service.DepartmentService) *DepartmentHandler {
	return &DepartmentHandler{
		departmentService: departmentService,
	}
}

func (h *DepartmentHandler) CreateDepartment(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var req models.CreateDepartmentRequest
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

	response, err := h.departmentService.CreateDepartment(c.Context(), userID, &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create department",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}

func (h *DepartmentHandler) UpdateDepartment(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	departmentID, err := strconv.ParseUint(c.Params("departmentId"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid department ID",
		})
	}

	var req models.UpdateDepartmentRequest
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

	response, err := h.departmentService.UpdateDepartment(c.Context(), userID, string(departmentID), &req)
	if err != nil {
		switch err.Error() {
		case "department not found":
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		case "unauthorized access to department":
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": err.Error(),
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to update department",
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *DepartmentHandler) DeleteDepartment(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	departmentID, err := strconv.ParseUint(c.Params("departmentId"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid department ID",
		})
	}

	err = h.departmentService.DeleteDepartment(c.Context(), userID, string(departmentID))
	if err != nil {
		switch err.Error() {
		case "department not found":
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		case "unauthorized access to department":
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": err.Error(),
			})
		case "department still contains employees":
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": err.Error(),
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to delete department",
			})
		}
	}

	return c.SendStatus(fiber.StatusOK)
}

func (h *DepartmentHandler) ListDepartments(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	filter := &models.DepartmentFilter{
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

	filter.Name = c.Query("name")

	departments, err := h.departmentService.ListDepartments(c.Context(), userID, filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch departments",
		})
	}

	return c.Status(fiber.StatusOK).JSON(departments)
}
