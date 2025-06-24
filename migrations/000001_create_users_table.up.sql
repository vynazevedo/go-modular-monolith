-- Criação da tabela de usuários
CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(36) NOT NULL PRIMARY KEY COMMENT 'UUID do usuário',
    email VARCHAR(255) NOT NULL UNIQUE COMMENT 'Email único do usuário',
    name VARCHAR(255) NOT NULL COMMENT 'Nome do usuário',
    status VARCHAR(20) NOT NULL DEFAULT 'active' COMMENT 'Status do usuário (active/inactive)',
    created_at BIGINT NOT NULL COMMENT 'Timestamp de criação em Unix time',
    
    INDEX idx_users_email (email),
    INDEX idx_users_status (status),
    INDEX idx_users_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Tabela de usuários do sistema';