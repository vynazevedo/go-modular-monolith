package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vynazevedo/go-modular-monolith/internal/modules/user"
	"github.com/vynazevedo/go-modular-monolith/internal/shared/config"
	"github.com/vynazevedo/go-modular-monolith/internal/shared/database"
	httpHandler "github.com/vynazevedo/go-modular-monolith/internal/shared/http"
	"github.com/vynazevedo/go-modular-monolith/internal/shared/middleware"
	"github.com/vynazevedo/go-modular-monolith/internal/shared/module"
	"github.com/vynazevedo/go-modular-monolith/pkg/logger"
	"gorm.io/gorm"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatalf("Failed to load configuration: %v", err)
	}

	logger.Init(cfg.Logger)
	logger.Info("Logger initialized successfully")

	gin.SetMode(cfg.Server.Mode)
	logger.Infof("Gin mode set to: %s", cfg.Server.Mode)

	db, err := database.Connect(cfg)
	if err != nil {
		logger.Fatalf("Failed to connect to database: %v", err)
	}
	logger.Info("Database connected successfully")

	if cfg.Database.Migrate {
		if err := database.RunMigrations(cfg, "migrations"); err != nil {
			logger.Fatalf("Failed to run migrations: %v", err)
		}
		logger.Info("Database migrations completed successfully")
	}

	userModuleSetup := func(db *gorm.DB) module.Module {
		return user.NewModule(db)
	}
	modules := module.SetupAllModules(db, userModuleSetup)

	router := gin.New()

	router.Use(gin.Recovery())
	router.Use(gin.Logger())
	router.Use(middleware.CORS(cfg))

	healthHandler := httpHandler.NewHandler()
	healthHandler.RegisterRoutes(router)

	api := router.Group("/api/v1")

	module.RegisterModules(api, modules...)

	port := cfg.Server.Port

	logger.Infof("Server starting on port %s", port)
	logger.Fatal(http.ListenAndServe(":"+port, router))
}
