@echo off
setlocal enabledelayedexpansion

REM Parse command line arguments
set DEBUG_MODE=0
set INSTALLER=1

if "%1"=="--debug" (
    set DEBUG_MODE=1
    set INSTALLER=0
    echo Building in DEBUG mode with developer tools
) else (
    echo Building in PRODUCTION mode with installer
)

REM Check if NSIS is available for installer builds
if %INSTALLER%==1 (
    where makensis >nul 2>&1
    if %ERRORLEVEL% NEQ 0 (
        echo Checking for NSIS in common locations...
        if exist "C:\Program Files (x86)\NSIS\makensis.exe" (
            set "PATH=%PATH%;C:\Program Files (x86)\NSIS"
        ) else if exist "C:\Program Files\NSIS\makensis.exe" (
            set "PATH=%PATH%;C:\Program Files\NSIS"
        ) else (
            echo WARNING: NSIS not found. Please install NSIS from https://nsis.sourceforge.io/Download
            set INSTALLER=0
        )
    )
)

REM Delete any existing application files to ensure fresh build
echo Cleaning previous build...
if exist "build\bin\trading-dashboard.exe" del "build\bin\trading-dashboard.exe"
if exist "build\bin\trading-dashboard-amd64-installer.exe" del /f "build\bin\trading-dashboard-amd64-installer.exe"

REM Build with appropriate options
echo Running wails build...

if %DEBUG_MODE%==1 (
    REM Debug build with devtools
    wails build -debug -devtools
) else (
    if %INSTALLER%==1 (
        REM Production build with NSIS installer
        wails build --nsis
    ) else (
        REM Standard production build
        wails build
    )
)

echo.
echo Build Results:
echo =============

if exist "build\bin\trading-dashboard.exe" (
    echo ✓ Executable: %CD%\build\bin\trading-dashboard.exe
    
    if %INSTALLER%==1 (
        if exist "build\bin\trading-dashboard-amd64-installer.exe" (
            echo ✓ Installer: %CD%\build\bin\trading-dashboard-amd64-installer.exe
        ) else (
            echo ✗ Installer was not created. Check if NSIS is properly installed.
        )
    )
) else (
    echo ✗ Build failed. No executable was created.
    exit /b 1
)

echo.
echo Build completed successfully!
echo.
echo Usage:
echo - Default (no args): Builds with installer
echo - --debug: Builds with developer tools

endlocal 