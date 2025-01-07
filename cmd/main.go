package main

import (
	"github.com/Project-Sprint-LDH-Team/GoGoManager/internal/configs"
	"github.com/Project-Sprint-LDH-Team/GoGoManager/pkg/internalsql"
	"github.com/gofiber/fiber/v2"
	"log"
)

//func main() {
//	// Initialize a new Fiber app
//	app := fiber.New()
//
//	// Define a route for the GET method on the root path '/'
//	app.Get("/", func(c *fiber.Ctx) error {
//		// Send a string response to the client
//		return c.SendString("Hello, World 👋!")
//	})
//
//	// Start the server on port 3000
//	log.Fatal(app.Listen(":3000"))
//}

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

	// Initialize Fiber app
	app := fiber.New()

	// Start server
	err = app.Listen(cfg.Service.Port)
	if err != nil {
		log.Fatalf("error starting server: %+v\n", err)
	}
}
