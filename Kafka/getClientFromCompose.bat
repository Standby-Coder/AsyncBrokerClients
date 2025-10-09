@echo off 
REM This script allows you to get a shell inside a specific service container defined in the compose.yaml file.
REM Usage:
REM   getClientFromCompose.bat producer-1
REM   getClientFromCompose.bat consumer-1

if "%1"=="" (
    echo Please provide the service name as an argument.
    echo Example: getClientFromCompose.bat producer-1
    exit /b 1
)

SET SERVICE_NAME=%1

IF "%SERVICE_NAME%"=="producer-1" (
    SET EXEC_NAME=./producer
) ELSE IF "%SERVICE_NAME%"=="consumer-1" (
    SET EXEC_NAME=./consumer
) ELSE (
    echo Unknown service: %SERVICE_NAME%
    exit /b 1
)

REM Check if compose is up and running, if not, start it
docker compose ps | findstr /C:"%SERVICE_NAME%" >nul
IF ERRORLEVEL 1 (
    echo Starting Docker Compose services...
    docker compose up -d
)

docker compose exec %SERVICE_NAME% sh -c "%EXEC_NAME%"