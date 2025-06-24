package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/vynazevedo/go-modular-monolith/internal/shared/config"
	"github.com/vynazevedo/go-modular-monolith/internal/shared/database"
	"github.com/vynazevedo/go-modular-monolith/pkg/logger"
)

func main() {
	var (
		action     = flag.String("action", "up", "Ação da migração: up, down, status, force")
		version    = flag.Int("version", -1, "Versão para forçar (apenas com action=force)")
		migrateDir = flag.String("dir", "migrations", "Diretório das migrações")
	)
	flag.Parse()

	logger.Init(logger.Config{
		Level:  "info",
		Format: "text",
	})

	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatalf("Erro ao carregar configuração: %v", err)
	}

	migrationService, err := database.GetMigrationService(cfg, *migrateDir)
	if err != nil {
		logger.Fatalf("Erro ao criar serviço de migração: %v", err)
	}

	switch *action {
	case "up":
		logger.Info("Executando migrações...")
		if err := migrationService.Up(); err != nil {
			logger.Fatalf("Erro ao executar migrações: %v", err)
		}
		logger.Info("Migrações executadas com sucesso!")

	case "down":
		logger.Info("Desfazendo uma migração...")
		if err := migrationService.Down(); err != nil {
			logger.Fatalf("Erro ao desfazer migração: %v", err)
		}
		logger.Info("Migração desfeita com sucesso!")

	case "status":
		if err := migrationService.Status(); err != nil {
			logger.Fatalf("Erro ao verificar status: %v", err)
		}

	case "force":
		if *version < 0 {
			logger.Fatal("Versão é obrigatória para action=force")
		}
		logger.Warnf("Forçando versão %d...", *version)
		if err := migrationService.Force(*version); err != nil {
			logger.Fatalf("Erro ao forçar versão: %v", err)
		}
		logger.Info("Versão forçada com sucesso!")

	default:
		fmt.Printf("Uso: %s -action=<up|down|status|force> [-version=<num>] [-dir=<path>]\n", os.Args[0])
		fmt.Println("\nAções disponíveis:")
		fmt.Println("  up     - Executa todas as migrações pendentes")
		fmt.Println("  down   - Desfaz a última migração")
		fmt.Println("  status - Mostra o status atual das migrações")
		fmt.Println("  force  - Força uma versão específica (requer -version)")
		fmt.Println("\nExemplos:")
		fmt.Printf("  %s -action=up\n", os.Args[0])
		fmt.Printf("  %s -action=down\n", os.Args[0])
		fmt.Printf("  %s -action=status\n", os.Args[0])
		fmt.Printf("  %s -action=force -version=1\n", os.Args[0])
		os.Exit(1)
	}
}
