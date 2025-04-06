@echo off
setlocal

echo.
echo ===============================
echo Iniciando ambiente Docker...
echo ===============================

:: Verifica se Docker Desktop está rodando
tasklist /FI "IMAGENAME eq Docker Desktop.exe" | find /I "Docker Desktop.exe" >nul
if errorlevel 1 (
    echo  Docker Desktop nao esta rodando. Iniciando...
    start "" "C:\Program Files\Docker\Docker\Docker Desktop.exe"
    timeout /t 15 >nul
)

:: Verifica se o daemon está pronto
:esperando_docker
docker info >nul 2>&1
if errorlevel 1 (
    echo Aguardando Docker inicializar...
    timeout /t 3 >nul
    goto esperando_docker
)

echo --->>>> Docker está ativo!

echo.
echo ===============================
echo Subindo containers...
echo ===============================
cd /d D:\Desenvolvimento_React\msd-service docker-compose up -d

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
