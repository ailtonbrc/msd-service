@echo off
echo Compilando...
go build -o bin/main.exe ./cmd/api
if %errorlevel% neq 0 (
    echo Erro ao compilar o programa.
    exit /b %errorlevel%
)
echo Build concluído!

echo Executando...
bin\main.exe
