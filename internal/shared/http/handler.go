package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) RegisterRoutes(router gin.IRouter) {
	router.GET("/health-check/alive", HealthCheckHandler)
}

// HealthCheckHandler rota para verificar a saúde do serviço
func HealthCheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "UP",
	})
}
