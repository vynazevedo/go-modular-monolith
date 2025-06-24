package http

import (
	"github.com/gin-gonic/gin"
	"github.com/vynazevedo/go-modular-monolith/internal/shared/middleware"
)

func (h *UserHandlers) RegisterRoutes(router *gin.RouterGroup) {
	// Rotas protegidas
	protected := router.Group("/", middleware.ValidateAPIKey())
	{
		protected.POST("/", h.CreateUser)
		protected.PUT("/:id", h.UpdateUser)
		protected.DELETE("/:id", h.DeleteUser)
		protected.PUT("/:id/activate", h.ActivateUser)
		protected.PUT("/:id/deactivate", h.DeactivateUser)
	}

	// Rotas públicas (apenas leitura)
	public := router.Group("/", middleware.OptionalAPIKey())
	{
		public.GET("/:id", h.GetUser)
		public.GET("/", h.ListUsers)
	}
}
