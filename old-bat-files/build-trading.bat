@echo off
setlocal enabledelayedexpansion

REM Parse command line arguments
set DEBUG_MODE=0
if "%1"=="--debug" (
    set DEBUG_MODE=1
    echo Building in DEBUG mode with developer tools enabled
) else (
    echo Building in PRODUCTION mode
)

echo Building application with NSIS installer...

REM Set environment variables
set PATH=%PATH%;C:\Program Files (x86)\NSIS

REM Delete any existing application files to ensure fresh build
echo Cleaning previous build...
if exist "build\bin\trading-dashboard.exe" del "build\bin\trading-dashboard.exe"
if exist "build\bin\trading-dashboard-amd64-installer.exe" del /f "build\bin\trading-dashboard-amd64-installer.exe"

REM Build with NSIS
echo Running wails build...

if %DEBUG_MODE%==1 (
    REM Debug build with devtools
    wails build -debug -devtools
) else (
    REM Production build with NSIS installer
    wails build --nsis
)

if exist "build\bin\trading-dashboard.exe" (
    echo Build successful.
    if %DEBUG_MODE%==0 (
        echo Installer should be in build\bin directory.
    ) else (
        echo Debug build completed. Run with run-debug.bat to see detailed logs.
    )
) else (
    echo Build failed. Check for errors above.
)

echo.
echo Build process completed.

endlocal