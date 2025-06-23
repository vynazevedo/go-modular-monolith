package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vynazevedo/go-modular-monolith/internal/modules/user"
	"github.com/vynazevedo/go-modular-monolith/internal/shared/config"
	"github.com/vynazevedo/go-modular-monolith/internal/shared/database"
	httpHandler "github.com/vynazevedo/go-modular-monolith/internal/shared/http"
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

	userModuleSetup := func(db *gorm.DB) module.Module {
		return user.NewModule(db)
	}
	modules := module.SetupAllModules(db, userModuleSetup)

	if err := database.AutoMigrate(db, module.GetAllModels(modules...)...); err != nil {
		logger.Fatalf("Failed to auto migrate models: %v", err)
	}
	logger.Info("Database migration completed successfully")

	router := gin.Default()

	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	healthHandler := httpHandler.NewHandler()
	healthHandler.RegisterRoutes(router)

	api := router.Group("/api/v1")

	module.RegisterModules(api, modules...)

	port := cfg.Server.Port

	logger.Infof("Server starting on port %s", port)
	logger.Fatal(http.ListenAndServe(":"+port, router))
}
