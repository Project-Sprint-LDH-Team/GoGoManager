package routes

import (
	"github.com/Project-Sprint-LDH-Team/GoGoManager/internal/handlers"
	"github.com/Project-Sprint-LDH-Team/GoGoManager/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(
	app *fiber.App,
	authMiddleware *middleware.AuthMiddleware,
	authHandler *handlers.AuthHandler,
	profileHandler *handlers.ProfileHandler,
	fileHandler *handlers.FileHandler,
	employeeHandler *handlers.EmployeeHandler,
	departmentHandler *handlers.DepartmentHandler,
) {
	// Routes
	api := app.Group("/v1")

	// Auth routes
	api.Post("/auth", authHandler.HandleAuth)

	// File upload route (protected)
	api.Post("/file", authMiddleware.AuthRequired(), fileHandler.UploadFile)

	// Profile routes (protected)
	profileRoutes := api.Group("/user", authMiddleware.AuthRequired())
	{
		profileRoutes.Get("/", profileHandler.GetProfile)
		profileRoutes.Patch("/", profileHandler.UpdateProfile)
	}

	// Employee routes (protected)
	employeeRoutes := api.Group("/employee", authMiddleware.AuthRequired())
	{
		employeeRoutes.Post("/", employeeHandler.CreateEmployee)
		employeeRoutes.Get("/", employeeHandler.ListEmployees)
		employeeRoutes.Patch("/:identityNumber", employeeHandler.UpdateEmployee)
		employeeRoutes.Delete("/:identityNumber", employeeHandler.DeleteEmployee)
	}

	// Department routes (protected)
	departmentRoutes := api.Group("/departement", authMiddleware.AuthRequired())
	{
		departmentRoutes.Post("/", departmentHandler.CreateDepartment)
		departmentRoutes.Get("/", departmentHandler.ListDepartments)
		departmentRoutes.Patch("/:departmentId", departmentHandler.UpdateDepartment)
		departmentRoutes.Delete("/:departmentId", departmentHandler.DeleteDepartment)
	}
}
