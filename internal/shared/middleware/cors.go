package middleware

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/vynazevedo/go-modular-monolith/internal/shared/config"
)

func CORS(cfg *config.Config) gin.HandlerFunc {
	c := cors.Config{
		AllowOrigins:     cfg.CORS.AllowedOrigins,
		AllowMethods:     cfg.CORS.AllowedMethods,
		AllowHeaders:     cfg.CORS.AllowedHeaders,
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           time.Duration(cfg.CORS.MaxAge) * time.Second,
	}

	return cors.New(c)
}
