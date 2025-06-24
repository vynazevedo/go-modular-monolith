package database

import (
	"database/sql"
	"fmt"

	"github.com/vynazevedo/go-modular-monolith/internal/shared/config"
	"github.com/vynazevedo/go-modular-monolith/internal/shared/migration"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Connect(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func ConnectWithEnv() (*gorm.DB, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	return Connect(cfg)
}

// RunMigrations executa as migrações de banco de dados
func RunMigrations(cfg *config.Config, migrationsDir string) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
	)

	sqlDB, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("erro ao conectar ao banco para migração: %w", err)
	}
	defer sqlDB.Close()

	migrationService := migration.NewService(sqlDB, migrationsDir)
	return migrationService.Up()
}

// GetMigrationService retorna uma instância do serviço de migração
func GetMigrationService(cfg *config.Config, migrationsDir string) (*migration.Service, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
	)

	sqlDB, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("erro ao conectar ao banco: %w", err)
	}

	return migration.NewService(sqlDB, migrationsDir), nil
}
