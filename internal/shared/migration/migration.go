package migration

import (
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/vynazevedo/go-modular-monolith/pkg/logger"
)

type Service struct {
	db            *sql.DB
	migrationsDir string
}

// NewService cria um novo serviço de migração
func NewService(db *sql.DB, migrationsDir string) *Service {
	return &Service{
		db:            db,
		migrationsDir: migrationsDir,
	}
}

// Up executa todas as migrações pendentes
func (s *Service) Up() error {
	m, err := s.createMigrator()
	if err != nil {
		return fmt.Errorf("erro ao criar migrator: %w", err)
	}
	defer m.Close()

	logger.Info("Executando migrações...")

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			logger.Info("Nenhuma migração pendente encontrada")
			return nil
		}
		return fmt.Errorf("erro ao executar migrações: %w", err)
	}

	logger.Info("Migrações executadas com sucesso")
	return nil
}

// Down desfaz uma migração
func (s *Service) Down() error {
	m, err := s.createMigrator()
	if err != nil {
		return fmt.Errorf("erro ao criar migrator: %w", err)
	}
	defer m.Close()

	logger.Warn("Desfazendo uma migração...")

	if err := m.Steps(-1); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			logger.Info("Nenhuma migração para desfazer")
			return nil
		}
		return fmt.Errorf("erro ao desfazer migração: %w", err)
	}

	logger.Info("Migração desfeita com sucesso")
	return nil
}

// Force força a versão da migração
func (s *Service) Force(version int) error {
	m, err := s.createMigrator()
	if err != nil {
		return fmt.Errorf("erro ao criar migrator: %w", err)
	}
	defer func(m *migrate.Migrate) {
		err, _ := m.Close()
		if err != nil {
			logger.Errorf("Erro ao fechar migrator: %v", err)
		} else {
			logger.Info("Migrator fechado com sucesso")
		}
	}(m)

	logger.Warnf("Forçando versão da migração para: %d", version)

	if err := m.Force(version); err != nil {
		return fmt.Errorf("erro ao forçar versão: %w", err)
	}

	logger.Info("Versão da migração forçada com sucesso")
	return nil
}

// Version retorna a versão atual da migração
func (s *Service) Version() (uint, bool, error) {
	m, err := s.createMigrator()
	if err != nil {
		return 0, false, fmt.Errorf("erro ao criar migrator: %w", err)
	}
	defer func(m *migrate.Migrate) {
		err, _ := m.Close()
		if err != nil {

		}
	}(m)

	version, dirty, err := m.Version()
	if err != nil {
		if errors.Is(err, migrate.ErrNilVersion) {
			return 0, false, nil
		}
		return 0, false, fmt.Errorf("erro ao obter versão: %w", err)
	}

	return version, dirty, nil
}

// Status retorna informações sobre o status das migrações
func (s *Service) Status() error {
	version, dirty, err := s.Version()
	if err != nil {
		return err
	}

	if version == 0 {
		logger.Info("Status: Nenhuma migração executada")
	} else {
		status := "limpo"
		if dirty {
			status = "inconsistente (requer intervenção manual)"
		}
		logger.Infof("Status: Versão %d (%s)", version, status)
	}

	return nil
}

// createMigrator cria uma instância do migrator
func (s *Service) createMigrator() (*migrate.Migrate, error) {
	driver, err := mysql.WithInstance(s.db, &mysql.Config{})
	if err != nil {
		return nil, fmt.Errorf("erro ao criar driver do banco: %w", err)
	}

	migrationsPath, err := filepath.Abs(s.migrationsDir)
	if err != nil {
		return nil, fmt.Errorf("erro ao resolver caminho das migrações: %w", err)
	}

	sourceURL := fmt.Sprintf("file://%s", migrationsPath)

	m, err := migrate.NewWithDatabaseInstance(sourceURL, "mysql", driver)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar instância de migração: %w", err)
	}

	return m, nil
}
