package handlers

import (
	"fmt"
	"github.com/Project-Sprint-LDH-Team/GoGoManager/internal/models"
	"github.com/Project-Sprint-LDH-Team/GoGoManager/internal/service"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"strings"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) HandleAuth(c *fiber.Ctx) error {
	var req models.AuthRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate request
	if validationErrors := h.validateRequest(&req); len(validationErrors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": validationErrors,
		})
	}

	// Process authentication
	response, err := h.authService.Authenticate(c.Context(), &req)
	if err != nil {
		switch err.Error() {
		case "email already exists":
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": err.Error(),
			})
		case "user not found":
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Internal server error",
			})
		}
	}

	// Return appropriate status based on action
	status := fiber.StatusOK
	if req.Action == "create" {
		status = fiber.StatusCreated
	}

	return c.Status(status).JSON(response)
}

var validate = validator.New()

// ValidationError untuk format error yang konsisten
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (h *AuthHandler) validateRequest(req *models.AuthRequest) []ValidationError {
	var validationErrors []ValidationError

	err := validate.Struct(req)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ValidationError
			element.Field = strings.ToLower(err.Field())

			switch err.Tag() {
			case "required":
				element.Message = fmt.Sprintf("%s is required", err.Field())
			case "email":
				element.Message = "Invalid email format"
			case "min":
				element.Message = fmt.Sprintf("%s must be at least %s characters long", err.Field(), err.Param())
			case "max":
				element.Message = fmt.Sprintf("%s must not exceed %s characters", err.Field(), err.Param())
			case "oneof":
				element.Message = fmt.Sprintf("%s must be one of: %s", err.Field(), err.Param())
			default:
				element.Message = fmt.Sprintf("%s is not valid", err.Field())
			}

			validationErrors = append(validationErrors, element)
		}
	}

	// Custom validation untuk format email jika diperlukan
	if req.Email != "" {
		if !strings.Contains(req.Email, "@") || !strings.Contains(req.Email, ".") {
			validationErrors = append(validationErrors, ValidationError{
				Field:   "email",
				Message: "Email must be in valid format",
			})
		}
	}

	// Custom validation untuk action
	if req.Action != "" && req.Action != "create" && req.Action != "login" {
		validationErrors = append(validationErrors, ValidationError{
			Field:   "action",
			Message: "Action must be either 'create' or 'login'",
		})
	}

	return validationErrors
}
