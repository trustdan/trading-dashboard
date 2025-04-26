@echo off
setlocal enabledelayedexpansion

REM Set environment variable
set CGO_ENABLED=1

REM Change to the application directory
cd /d "%~dp0"

REM Create log file
echo Starting application at %date% %time% > "%~dp0trading-dashboard-log.txt"
echo CGO_ENABLED=%CGO_ENABLED% >> "%~dp0trading-dashboard-log.txt"
echo Working directory: %CD% >> "%~dp0trading-dashboard-log.txt"

REM List files in directory
dir /b >> "%~dp0trading-dashboard-log.txt"

REM Launch application
echo Launching application... >> "%~dp0trading-dashboard-log.txt"
"%~dp0trading-dashboard.exe" > "%~dp0app-output.txt" 2>&1

REM Check exit code
if %ERRORLEVEL% NEQ 0 (
    echo Application exited with code %ERRORLEVEL% >> "%~dp0trading-dashboard-log.txt"
) else (
    echo Application exited normally >> "%~dp0trading-dashboard-log.txt"
)

endlocal
