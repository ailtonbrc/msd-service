@echo off
echo ================================
echo Iniciando push para o GitHub...
echo ================================

REM Navega até a raiz do projeto (caso esteja fora)
cd /d %~dp0

REM Adiciona todos os arquivos modificados
git add .

REM Solicita uma mensagem de commit ao usuário
set /p msg=Digite a mensagem de commit: 

REM Faz o commit
git commit -m "%msg%"

REM Faz o push para o repositório remoto
git push

echo.
echo Enviado para o git (Push) concluído com sucesso!
pause