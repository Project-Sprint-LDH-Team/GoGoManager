package main

import (
	"context"
	"log"

	"github.com/Project-Sprint-LDH-Team/GoGoManager/internal/configs"
	"github.com/Project-Sprint-LDH-Team/GoGoManager/internal/handlers"
	"github.com/Project-Sprint-LDH-Team/GoGoManager/internal/middleware"
	"github.com/Project-Sprint-LDH-Team/GoGoManager/internal/repository"
	"github.com/Project-Sprint-LDH-Team/GoGoManager/internal/routes"
	"github.com/Project-Sprint-LDH-Team/GoGoManager/internal/service"
	"github.com/Project-Sprint-LDH-Team/GoGoManager/pkg/internalsql"
	"github.com/Project-Sprint-LDH-Team/GoGoManager/pkg/jwt"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gofiber/fiber/v2"
)

func main() {
	var (
		cfg *configs.Config
	)

	// Initialize configs
	err := configs.Init(
		configs.WithConfigFolder([]string{
			"./configs/",
			"./internal/configs/",
		}),
		configs.WithConfigFile("config"),
		configs.WithConfigType("yaml"),
	)
	if err != nil {
		log.Fatalf("error initializing configs: %+v\n", err)
	}
	cfg = configs.Get()

	// Connect to database
	db, err := internalsql.Connect(cfg.Database.DataSourceName)
	if err != nil {
		log.Fatalf("error connecting to database %+v\n", err)
	}

	// Initialize AWS S3 client
	awsCfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(cfg.AWS.Region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			cfg.AWS.AccessKeyID,
			cfg.AWS.SecretAccessKey,
			"",
		)),
	)
	if err != nil {
		log.Fatal("\033[31mUnable to load AWS SDK config:\033[0m", err)
	}

	s3Client := s3.NewFromConfig(awsCfg)

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	employeeRepo := repository.NewEmployeeRepository(db)
	departmentRepo := repository.NewDepartmentRepository(db)

	// Initialize JWT maker
	jwtMaker := jwt.NewJWTMaker(cfg.Service.SecretJWT)

	// Initialize services
	authService := service.NewAuthService(userRepo, jwtMaker)
	profileService := service.NewProfileService(userRepo)
	storageService := service.NewS3StorageService(s3Client, cfg.AWS.Bucket)
	fileService := service.NewFileService(storageService)
	employeeService := service.NewEmployeeService(employeeRepo, departmentRepo)
	departmentService := service.NewDepartmentService(departmentRepo)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	profileHandler := handlers.NewProfileHandler(profileService)
	fileHandler := handlers.NewFileHandler(fileService, profileService)
	employeeHandler := handlers.NewEmployeeHandler(employeeService)
	departmentHandler := handlers.NewDepartmentHandler(departmentService)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(jwtMaker)

	// Initialize Fiber app
	app := fiber.New()

	// setup routes
	routes.SetupRoutes(
		app,
		authMiddleware,
		authHandler,
		profileHandler,
		fileHandler,
		employeeHandler,
		departmentHandler,
	)

	// Start server
	err = app.Listen(":" + cfg.Service.Port)
	if err != nil {
		log.Fatalf("error starting server: %+v\n", err)
	}

	//port := fmt.Sprintf(":%d", cfg.Service.Port)
	//log.Printf("Server is starting on port %s", port)
	//log.Fatal(app.Listen(port))
}
