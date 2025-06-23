// Package module provides an interface and utility functions for managing application modules in a modular monolith architecture
package module

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Module interface {
	RegisterRoutes(router *gin.RouterGroup)
	GetModels() []any
}

func RegisterModules(router *gin.RouterGroup, modules ...Module) {
	for _, module := range modules {
		module.RegisterRoutes(router)
	}
}

func GetAllModels(modules ...Module) []any {
	var models []any
	for _, module := range modules {
		models = append(models, module.GetModels()...)
	}
	return models
}

type Setup func(db *gorm.DB) Module

func SetupAllModules(db *gorm.DB, moduleSetups ...Setup) []Module {
	modules := make([]Module, len(moduleSetups))
	for i, setup := range moduleSetups {
		modules[i] = setup(db)
	}
	return modules
}
