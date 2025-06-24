# Guia de Migrações de Banco de Dados

Este projeto utiliza [golang-migrate/migrate](https://github.com/golang-migrate/migrate) para gerenciamento de migrações de banco de dados, seguindo as melhores práticas da indústria.

## Estrutura de Arquivos

```
migrations/
├── 000001_create_users_table.up.sql     # Migração para frente
├── 000001_create_users_table.down.sql   # Rollback da migração
├── 000002_add_user_profile.up.sql       # Próxima migração
└── 000002_add_user_profile.down.sql     # Rollback correspondente
```

### Convenção de Nomenclatura

```
{version}_{description}.{direction}.sql

Onde:
- version: Número sequencial (6 dígitos): 000001, 000002, etc.
- description: Descrição em snake_case: create_users_table
- direction: up (aplicar) ou down (reverter)
```

## Comandos Disponíveis

### Executar Migrações
```bash
# Aplicar todas as migrações pendentes
make migrate-up

# Verificar status atual
make migrate-status
```

### Reverter Migrações
```bash
# Desfazer a última migração
make migrate-down
```

### Criar Nova Migração
```bash
# Criar arquivos de migração
make migrate-create NAME=add_user_avatar

# Isso criará:
# migrations/000002_add_user_avatar.up.sql
# migrations/000002_add_user_avatar.down.sql
```

### Resolução de Problemas
```bash
# Forçar versão específica (use com cuidado!)
make migrate-force VERSION=1

# Verificar status antes de forçar
make migrate-status
```
---

## Comandos de Referência Rápida

```bash
# Comandos básicos
make migrate-up                    # Aplicar migrações
make migrate-down                  # Desfazer última migração  
make migrate-status                # Ver status atual
make migrate-create NAME=exemplo   # Criar nova migração
make migrate-force VERSION=1       # Forçar versão (emergência)

# Comandos diretos (alternativa)
go run cmd/migrate/main.go -action=up
go run cmd/migrate/main.go -action=down
go run cmd/migrate/main.go -action=status
go run cmd/migrate/main.go -action=force -version=1
```

**Lembre-se**: Migrações são irreversíveis em produção. Sempre teste tudo localmente primeiro!