@echo off
setlocal enabledelayedexpansion

echo Trading Dashboard Installation Fix Script
echo =======================================
echo.

REM Try to find the installation directory
echo Looking for the installation directory...
for /f "tokens=*" %%a in ('reg query "HKLM\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall" /s /f "Trading Dashboard" ^| findstr "InstallLocation"') do (
    for /f "tokens=2*" %%b in ('echo %%a') do (
        set "INSTALL_DIR=%%c"
    )
)

if not defined INSTALL_DIR (
    echo Installation not found in registry. Searching file system...
    for /f "tokens=*" %%a in ('dir /b /s /ad "C:\Program Files\*Trading Dashboard*" 2^>nul') do (
        set "INSTALL_DIR=%%a"
    )
)

if not defined INSTALL_DIR (
    echo Could not locate Trading Dashboard installation.
    echo Please enter the installation path manually:
    set /p INSTALL_DIR="Installation path: "
)

if not exist "!INSTALL_DIR!" (
    echo Installation directory does not exist: !INSTALL_DIR!
    exit /b 1
)

echo Found installation at: !INSTALL_DIR!

REM Create proper launcher script
echo Creating launcher script...
(
    echo @echo off
    echo setlocal enabledelayedexpansion
    echo.
    echo cd /d "%%~dp0"
    echo.
    echo REM Create log file
    echo echo Starting application at %%date%% %%time%% ^> "%%~dp0trading-dashboard-log.txt"
    echo echo Working directory: %%CD%% ^>^> "%%~dp0trading-dashboard-log.txt"
    echo.
    echo REM List files
    echo dir /b ^>^> "%%~dp0trading-dashboard-log.txt"
    echo.
    echo REM Start application
    echo "%%~dp0trading-dashboard.exe"
    echo.
    echo REM Check exit code
    echo if %%ERRORLEVEL%% NEQ 0 (
    echo     echo Application exited with error code %%ERRORLEVEL%% ^>^> "%%~dp0trading-dashboard-log.txt"
    echo     echo Please check the log file for more information.
    echo     pause
    echo ^)
    echo.
    echo endlocal
) > "!INSTALL_DIR!\run-trading-dashboard.bat"

echo Creating desktop shortcut...
powershell -Command "$WshShell = New-Object -ComObject WScript.Shell; $Shortcut = $WshShell.CreateShortcut([Environment]::GetFolderPath('Desktop') + '\Trading Dashboard.lnk'); $Shortcut.TargetPath = '!INSTALL_DIR!\run-trading-dashboard.bat'; $Shortcut.IconLocation = '!INSTALL_DIR!\trading-dashboard.exe,0'; $Shortcut.WorkingDirectory = '!INSTALL_DIR!'; $Shortcut.Save()"

echo.
echo Fix completed. Please try running the application using:
echo 1. The new desktop shortcut
echo 2. !INSTALL_DIR!\run-trading-dashboard.bat
echo.
echo If you still encounter issues, please check the log file at:
echo !INSTALL_DIR!\trading-dashboard-log.txt

endlocal 