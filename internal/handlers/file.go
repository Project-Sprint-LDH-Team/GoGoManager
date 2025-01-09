package handlers

import (
	"github.com/Project-Sprint-LDH-Team/GoGoManager/internal/models"
	"github.com/Project-Sprint-LDH-Team/GoGoManager/internal/service"
	"github.com/gofiber/fiber/v2"
)

type FileHandler struct {
	fileService    service.FileService
	profileService service.ProfileService
}

func NewFileHandler(fileService service.FileService, profileService service.ProfileService) *FileHandler {
	return &FileHandler{
		fileService:    fileService,
		profileService: profileService,
	}
}

func (h *FileHandler) UploadFile(c *fiber.Ctx) error {
	// Get user ID from context (set by auth middleware)
	userID := c.Locals("userID").(uint)

	// Get user profile for email
	userProfile, err := h.profileService.GetProfile(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get user profile",
		})
	}

	// Get file from form
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No file uploaded",
		})
	}

	// Upload file
	response, err := h.fileService.UploadFile(c.Context(), file, userID, userProfile.Email)
	if err != nil {
		// Check if it's a validation error
		if validationErrors, ok := err.(models.FileValidationErrors); ok {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"errors": validationErrors,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to upload file",
		})
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
