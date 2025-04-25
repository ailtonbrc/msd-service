# Passo a passo para cadastrar uma Nova entidade

## 1. Criar o arquivo Model base
- **Path:** `internal/models/arquivo.go`
- Definir a estrutura da tabela com tags GORM e JSON
- Definir estruturas para `CreateRequest` e `UpdateRequest` com validações

## 2. Criar o arquivo Model DTO
- **Path:** `internal/models/arquivo_dto.go`
- Incluir DTOs para listagem e detalhes
- Implementar métodos de conversão (`ToDTO`, `ToDetailDTO`)

## 3. Criar o arquivo de Repositório
- **Path:** `internal/repository/arquivo_repository.go`
- Definir interface e implementação GORM

## 4. Criar o arquivo de Validação
- **Path:** `internal/validator/arquivo_validator.go`
- Implementar validações de negócio

## 5. Criar o arquivo de Services
- **Path:** `internal/service/arquivo_service.go`
- Implementar lógica de negócio usando repository e validator

## 6. Criar o arquivo de Handlers
- **Path:** `internal/api/handlers/arquivo.go`
- Implementar endpoints com anotações Swagger

## 7. Criar o arquivo de Rotas
- **Path:** `internal/api/routes/arquivo.go`
- Configurar rotas com middlewares de autenticação e autorização

## 8. Ajustar o arquivo `Server.go`
- **Path:** `internal/api/server/server.go`
- Incluir as novas rotas no método `SetupServer`

## 9. Atualizar o arquivo de Migrações
- Garantir que a nova entidade seja incluída no processo de migração automática
- Verificar se é necessário adicionar no `internal/models/migrate.go`

## 10. Configurar o Seeder
- **Path:** `internal/models/seeder.go`
- Adicionar dados iniciais para a nova entidade
- Adicionar novas permissões relacionadas à entidade

## 11. Atualizar permissões
- Adicionar novas permissões para a entidade (`view`, `create`, `edit`, `delete`)
- Associar as permissões aos perfis existentes conforme necessário

## 12. Testes unitários (opcional, mas recomendado)
- Criar testes para repository, validator e service

---

Este passo a passo cobre todos os aspectos necessários para adicionar uma nova entidade completa ao sistema, desde o modelo de dados até a API e permissões.

