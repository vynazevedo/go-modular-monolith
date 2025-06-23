package main

import (
	"errors"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/vynazevedo/template-go-modular/internal/modules/user"
	"github.com/vynazevedo/template-go-modular/internal/shared/config"
	"github.com/vynazevedo/template-go-modular/internal/shared/database"
	"github.com/vynazevedo/template-go-modular/internal/shared/http"
	"github.com/vynazevedo/template-go-modular/internal/shared/module"
	"gorm.io/gorm"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	userModuleSetup := func(db *gorm.DB) module.Module {
		return user.NewModule(db)
	}
	modules := module.SetupAllModules(db, userModuleSetup)

	if err := database.AutoMigrate(db, module.GetAllModels(modules...)...); err != nil {
		log.Fatalf("Failed to auto migrate models: %v", err)
	}

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			var e *fiber.Error
			if errors.As(err, &e) {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	app.Use(recover.New())
	app.Use(logger.New())

	healthHandler := http.NewHandler()
	healthHandler.RegisterRoutes(app)

	api := app.Group("/api/v1")

	module.RegisterModules(api, modules...)

	port := cfg.Server.Port

	log.Printf("Server starting on port %s", port)
	log.Fatal(app.Listen(":" + port))
}
