package handlers

import (
	"fmt"
	"github.com/Project-Sprint-LDH-Team/GoGoManager/internal/models"
	"github.com/Project-Sprint-LDH-Team/GoGoManager/internal/service"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type ProfileHandler struct {
	profileService service.ProfileService
}

func NewProfileHandler(profileService service.ProfileService) *ProfileHandler {
	return &ProfileHandler{
		profileService: profileService,
	}
}

func (h *ProfileHandler) GetProfile(c *fiber.Ctx) error {
	// Get userID from context (set by auth middleware)
	userID := c.Locals("userID").(uint)

	user, err := h.profileService.GetProfile(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(user.ToProfileResponse())
}

func (h *ProfileHandler) UpdateProfile(c *fiber.Ctx) error {
	// Get userID from context (set by auth middleware)
	userID := c.Locals("userID").(uint)

	var req models.UpdateProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate request
	if errors := validateUpdateProfileRequest(&req); len(errors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": errors,
		})
	}

	// Update profile
	user, err := h.profileService.UpdateProfile(c.Context(), userID, &req)
	if err != nil {
		switch err.Error() {
		case "email is already used by another user":
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": err.Error(),
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(user.ToProfileResponse())
}

func validateUpdateProfileRequest(req *models.UpdateProfileRequest) []ValidationError {
	var validationErrors []ValidationError

	err := validate.Struct(req)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ValidationError
			element.Field = err.Field()

			switch err.Tag() {
			case "required":
				element.Message = fmt.Sprintf("%s is required", err.Field())
			case "email":
				element.Message = "Invalid email format"
			case "min":
				element.Message = fmt.Sprintf("%s must be at least %s characters", err.Field(), err.Param())
			case "max":
				element.Message = fmt.Sprintf("%s must not exceed %s characters", err.Field(), err.Param())
			case "uri":
				element.Message = fmt.Sprintf("%s must be a valid URI", err.Field())
			}

			validationErrors = append(validationErrors, element)
		}
	}

	return validationErrors
}
