package config

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/spf13/viper"
	"github.com/vynazevedo/go-modular-monolith/pkg/logger"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Logger   logger.Config
	CORS     CORSConfig
}

type ServerConfig struct {
	Port string
	Mode string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	Migrate  bool
}

type CORSConfig struct {
	AllowedOrigins []string
	AllowedMethods []string
	AllowedHeaders []string
	MaxAge         int
}

func LoadConfig() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	viper.SetDefault("PORT", "8080")
	viper.SetDefault("GIN_MODE", "debug")
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", "3306")
	viper.SetDefault("DB_USER", "root")
	viper.SetDefault("DB_PASSWORD", "password")
	viper.SetDefault("DB_NAME", "modular_monolith")
	viper.SetDefault("DB_AUTO_MIGRATE", false)
	viper.SetDefault("LOG_LEVEL", "info")
	viper.SetDefault("LOG_FORMAT", "text")
	viper.SetDefault("SERVICE_NAME", "go-modular-monolith")
	viper.SetDefault("CORS_ALLOWED_ORIGINS", "http://localhost:3000,http://localhost:5173")
	viper.SetDefault("CORS_ALLOWED_METHODS", "GET,POST,PUT,DELETE,OPTIONS")
	viper.SetDefault("CORS_ALLOWED_HEADERS", "Origin,Content-Type,Accept,Authorization,X-API-Key")
	viper.SetDefault("CORS_MAX_AGE", 86400)

	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if !errors.As(err, &configFileNotFoundError) {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
		log.Println("No .env file found, using environment variables and defaults")
	}

	config := &Config{
		Server: ServerConfig{
			Port: viper.GetString("PORT"),
			Mode: viper.GetString("GIN_MODE"),
		},
		Database: DatabaseConfig{
			Host:     viper.GetString("DB_HOST"),
			Port:     viper.GetString("DB_PORT"),
			User:     viper.GetString("DB_USER"),
			Password: viper.GetString("DB_PASSWORD"),
			Name:     viper.GetString("DB_NAME"),
			Migrate:  viper.GetBool("DB_AUTO_MIGRATE"),
		},
		Logger: logger.Config{
			Level:       viper.GetString("LOG_LEVEL"),
			Format:      viper.GetString("LOG_FORMAT"),
			ServiceName: viper.GetString("SERVICE_NAME"),
		},
		CORS: CORSConfig{
			AllowedOrigins: strings.Split(viper.GetString("CORS_ALLOWED_ORIGINS"), ","),
			AllowedMethods: strings.Split(viper.GetString("CORS_ALLOWED_METHODS"), ","),
			AllowedHeaders: strings.Split(viper.GetString("CORS_ALLOWED_HEADERS"), ","),
			MaxAge:         viper.GetInt("CORS_MAX_AGE"),
		},
	}

	return config, nil
}
