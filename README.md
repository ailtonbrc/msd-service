# Simple ERP Service

Backend em Golang para um sistema ERP simples, utilizando Gin e GORM.

## Requisitos

- Go 1.21 ou superior
- PostgreSQL 12 ou superior

## Configuração

1. Crie o Banco de Dados no SGBD
2. Configure o arquivo `.env` com suas credenciais de banco de dados
3. Execute o script build.bat na raiz do projeto para compilar `.\build.bat`
4. Ele vai Executar as Migrações e perguntar se deseja executar os seeders:
   Deseja rodar os seeders? (s/n):
   Na primeira vez coloque 's' para rodas nas próximas não precisa, apenas se criar novos seeders no sistema.

## Configuração GIT

Antes de dar o Commit de alterações precisa rodar os comandos.
`git config --local user.name seuUsuarioGit`
`git config --local user.email seuEmailGit`
