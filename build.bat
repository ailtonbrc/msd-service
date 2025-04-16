@echo off
echo Compilando...
go build -o bin/clinica_server.exe ./cmd/api
if %errorlevel% neq 0 (
    echo Erro ao compilar o programa.
    exit /b %errorlevel%
)
echo Build concluído!

echo Executando...
bin\clinica_server.exe
