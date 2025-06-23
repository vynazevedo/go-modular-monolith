package http

import "github.com/gin-gonic/gin"

func (h *UserHandlers) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/", h.CreateUser)
	router.GET("/:id", h.GetUser)
	router.PUT("/:id", h.UpdateUser)
	router.DELETE("/:id", h.DeleteUser)
	router.PUT("/:id/activate", h.ActivateUser)
	router.PUT("/:id/deactivate", h.DeactivateUser)
	router.GET("/", h.ListUsers)
}
