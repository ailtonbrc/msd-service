@echo off
setlocal

echo.
echo ===============================
echo Compilando e executando backend...
echo ===============================
cd server
go build -o clinica_server.exe ./cmd/api

echo.
echo Iniciando API...
echo ===============================
clinica_server.exe


endlocal
pause
