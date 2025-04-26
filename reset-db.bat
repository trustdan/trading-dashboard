@echo off
setlocal enabledelayedexpansion

echo WARNING: This will delete all database files and reset the application to a fresh state.
echo All your data will be lost.
echo.
echo Press Ctrl+C to cancel, or any key to continue...
pause > nul

REM Get the data directory path
set "DATA_DIR=%APPDATA%\TradingDashboard\data"
echo Using data directory: %DATA_DIR%

REM Check if process is running and try to close it
tasklist /fi "imagename eq trading-dashboard.exe" | find "trading-dashboard.exe" > nul
if %ERRORLEVEL% EQU 0 (
    echo Trading Dashboard is currently running. Attempting to close it...
    taskkill /F /IM trading-dashboard.exe
    timeout /t 2 > nul
)

REM Delete database files
echo Deleting database files...
if exist "%DATA_DIR%" (
    del /F /Q "%DATA_DIR%\*.vlog" 2>nul
    del /F /Q "%DATA_DIR%\*.sst" 2>nul
    del /F /Q "%DATA_DIR%\MANIFEST" 2>nul
    del /F /Q "%DATA_DIR%\KEYREGISTRY" 2>nul
    del /F /Q "%DATA_DIR%\DISCARD" 2>nul
    del /F /Q "%DATA_DIR%\trading.db" 2>nul
    
    echo Database files deleted.
) else (
    echo Data directory not found. Nothing to delete.
)

REM Build and run the application in debug mode
echo Building application in debug mode...
call .\build.bat --debug

echo Starting application...
start "Trading Dashboard" "build\bin\trading-dashboard.exe"

echo Done. Application has been reset and restarted.
endlocal 