// Package config provides configuration management for the applications
package config

import (
	"errors"
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
}

type ServerConfig struct {
	Port string
}

type DatabaseConfig struct {
	URL string
}

func LoadConfig() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	viper.SetDefault("PORT", "8080")
	viper.SetDefault("DATABASE_URL", "root:password@tcp(localhost:3306)/modular_monolith?charset=utf8mb4&parseTime=True&loc=Local")

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
		},
		Database: DatabaseConfig{
			URL: viper.GetString("DATABASE_URL"),
		},
	}

	return config, nil
}
