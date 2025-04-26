@echo off
setlocal enabledelayedexpansion

echo Setting up environment for Trading Dashboard...

REM Set environment variables
set WAILS_DEBUG=1
set WAILS_LOG_LEVEL=DEBUG

REM Path to the executable
set "EXE_PATH=%~dp0\build\bin\trading-dashboard.exe"

REM Check if the executable exists
if not exist "%EXE_PATH%" (
    echo ERROR: Application executable not found at: %EXE_PATH%
    echo Make sure you build the application first with: wails build
    exit /b 1
)

echo Environment prepared. Starting application...
echo WAILS_DEBUG=%WAILS_DEBUG%
echo WAILS_LOG_LEVEL=%WAILS_LOG_LEVEL%

REM Launch the application with redirection to log files
cd /d "%~dp0\build\bin"
echo Starting application at %date% %time% > debug-log.txt
echo WAILS_DEBUG=%WAILS_DEBUG% >> debug-log.txt
echo WAILS_LOG_LEVEL=%WAILS_LOG_LEVEL% >> debug-log.txt
echo Working directory: %CD% >> debug-log.txt
dir /b >> debug-log.txt

echo Launching application from: %CD%
echo Log files will be created in this directory.
echo IMPORTANT: The application window might be blank or not visible.
echo Check debug-output.txt for any errors.

REM Run the app with output to console window for better visibility
start "Trading Dashboard Debug" cmd /k "%EXE_PATH%"

echo The application is running in a separate window.
echo Press any key to end debug session...
pause > nul

taskkill /F /IM trading-dashboard.exe 2>nul

endlocal 