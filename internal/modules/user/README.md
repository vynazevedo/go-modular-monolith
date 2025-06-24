# Módulo User

> Módulo de gestão de usuários

## API Endpoints

| Método | Endpoint | Auth | Descrição |
|--------|----------|------|-----------|
| POST | `/users/` | Obrigatória | Criar usuário |
| GET | `/users/:id` | Opcional | Buscar por ID |
| GET | `/users/` | Opcional | Listar (paginado) |
| PUT | `/users/:id` | Obrigatória | Atualizar nome |
| DELETE | `/users/:id` | Obrigatória | Remover usuário |
| PUT | `/users/:id/activate` | Obrigatória | Ativar usuário |
| PUT | `/users/:id/deactivate` | Obrigatória | Desativar usuário |

**Auth**: Rotas protegidas exigem header `X-API-Key: api-key-exemplo`
