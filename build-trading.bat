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

REM Check if sqlite3.dll exists in the current directory
if not exist "%~dp0\sqlite3.dll" (
    echo SQLite3 DLL not found. Downloading...
    
    REM Download SQLite3 DLL
    powershell -Command "Invoke-WebRequest -Uri 'https://www.sqlite.org/2023/sqlite-dll-win64-x64-3430000.zip' -OutFile '%TEMP%\sqlite3.zip' -UseBasicParsing"
    
    REM Check if download was successful
    if %ERRORLEVEL% NEQ 0 (
        echo Failed to download SQLite3 DLL. Creating an empty placeholder file.
        echo This is a placeholder file > "%~dp0\sqlite3.dll"
    ) else (
        REM Extract the DLL
        powershell -Command "Expand-Archive -Path '%TEMP%\sqlite3.zip' -DestinationPath '%TEMP%\sqlite3' -Force"
        
        REM Copy the DLL to the current directory
        copy "%TEMP%\sqlite3\sqlite3.dll" "%~dp0\sqlite3.dll"
        
        REM Clean up temporary files
        del "%TEMP%\sqlite3.zip"
        rmdir /s /q "%TEMP%\sqlite3"
    )
)

echo Building application with NSIS installer...

REM Set environment variables - MUST BE SET BEFORE CALLING WAILS!
set CGO_ENABLED=1
set CC=gcc
set PATH=%PATH%;C:\msys64\ucrt64\bin;C:\Program Files (x86)\NSIS

REM Verify environment settings
echo CGO_ENABLED=%CGO_ENABLED%
echo CC=%CC%
echo PATH=%PATH%

REM Delete any existing application files to ensure fresh build
echo Cleaning previous build...
if exist "build\bin\trading-dashboard.exe" del "build\bin\trading-dashboard.exe"
if exist "build\bin\trading-dashboard-amd64-installer.exe" del /f "build\bin\trading-dashboard-amd64-installer.exe"

REM Build with NSIS
echo Running wails build with CGO_ENABLED=%CGO_ENABLED%...

if %DEBUG_MODE%==1 (
    REM Debug build with devtools
    wails build -debug -devtools
) else (
    REM Production build with NSIS installer
    wails build --nsis
)

REM Verify the build was done with CGO enabled
if exist "build\bin\trading-dashboard.exe" (
    echo Build successful. Checking if CGO was properly enabled...
    cd /d "%~dp0\build\bin"
    set CGO_ENABLED=1
    echo. > cgo-test.txt
    "%~dp0\build\bin\trading-dashboard.exe" 2>> cgo-test.txt
    findstr /C:"CGO_ENABLED=0" /C:"requires cgo" cgo-test.txt > nul
    if !errorlevel! equ 0 (
        echo WARNING: The executable was compiled without CGO support! 
        echo Please try rebuilding with the following commands manually:
        echo.
        echo set CGO_ENABLED=1
        echo set CC=gcc
        echo wails build --nsis
    ) else (
        echo CGO appears to be correctly enabled.
        if %DEBUG_MODE%==0 (
            echo Installer should be in build\bin directory.
        ) else (
            echo Debug build completed. Run with run-debug.bat to see detailed logs.
        )
    )
) else (
    echo Build failed. Check for errors above.
)

echo.
echo Build process completed.

endlocal