package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vynazevedo/go-modular-monolith/pkg/logger"
)

const (
	HeaderAPIKey   = "X-API-Key"
	RequiredAPIKey = "api-key-exemplo"
)

// ValidateAPIKey middleware que valida o header X-API-Key
func ValidateAPIKey() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader(HeaderAPIKey)

		if apiKey == "" {
			logger.Warn("API Key missing in request")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "API Key é obrigatório",
				"code":  "MISSING_API_KEY",
			})
			c.Abort()
			return
		}

		if apiKey != RequiredAPIKey {
			logger.WithField("provided_key", apiKey).Warn("Invalid API Key provided")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "API Key inválido",
				"code":  "INVALID_API_KEY",
			})
			c.Abort()
			return
		}

		logger.Debug("API Key validated successfully")
		c.Next()
	}
}

// OptionalAPIKey middleware que valida API Key apenas se fornecido
func OptionalAPIKey() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader(HeaderAPIKey)

		if apiKey != "" && apiKey != RequiredAPIKey {
			logger.WithField("provided_key", apiKey).Warn("Invalid optional API Key provided")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "API Key fornecido é inválido",
				"code":  "INVALID_API_KEY",
			})
			c.Abort()
			return
		}

		if apiKey == RequiredAPIKey {
			logger.Debug("Optional API Key validated successfully")
			c.Set("authenticated", true)
		}

		c.Next()
	}
}
