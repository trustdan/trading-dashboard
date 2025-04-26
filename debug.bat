@echo off
setlocal enabledelayedexpansion

echo Setting up debug environment for Trading Dashboard...

REM Set environment variables
set WAILS_DEBUG=1
set WAILS_LOG_LEVEL=DEBUG

REM Path to the executable
set "EXE_PATH=%~dp0\build\bin\trading-dashboard.exe"

REM Check if the executable exists
if not exist "%EXE_PATH%" (
    echo ERROR: Application executable not found at: %EXE_PATH%
    echo Building a debug version first...
    call .\build.bat --debug
    if not exist "%EXE_PATH%" (
        echo Build failed. Exiting.
        exit /b 1
    )
)

REM Check data directory
set "DATA_DIR=%APPDATA%\TradingDashboard\data"
echo Checking data directory: %DATA_DIR%
if not exist "%DATA_DIR%" (
    echo Creating data directory...
    mkdir "%DATA_DIR%"
)
echo Data directory contents:
dir "%DATA_DIR%"

REM Launch the application with logging
echo Launching application in debug mode...
cd /d "%~dp0\build\bin"

REM Create log file directory
if not exist "logs" mkdir logs

REM Generate log filename with timestamp
for /f "tokens=2 delims==" %%a in ('wmic OS Get localdatetime /value') do set "dt=%%a"
set "LOGTIME=%dt:~0,8%-%dt:~8,6%"
set "LOG_FILE=logs\debug-%LOGTIME%.txt"

echo Starting application at %date% %time% > "%LOG_FILE%"
echo WAILS_DEBUG=%WAILS_DEBUG% >> "%LOG_FILE%"
echo WAILS_LOG_LEVEL=%WAILS_LOG_LEVEL% >> "%LOG_FILE%"
echo Data directory: %DATA_DIR% >> "%LOG_FILE%"
echo Working directory: %CD% >> "%LOG_FILE%"
dir /b >> "%LOG_FILE%"

echo Launching application with debug output...
echo Log file: %CD%\%LOG_FILE%
start "Trading Dashboard Debug" cmd /k "%EXE_PATH% > "%LOG_FILE%" 2>&1"

echo Application is running in a separate window.
echo Press any key to end debug session...
pause > nul

taskkill /F /IM trading-dashboard.exe 2>nul

echo Debug session ended. Check logs at: %CD%\%LOG_FILE%

endlocal 