@echo off
setlocal enabledelayedexpansion

REM Find the installed folder
echo Looking for the installation directory...
set "PRODUCT_NAME=Trading Dashboard"
for /f "tokens=*" %%a in ('dir /b /s /a:d "C:\Program Files\*%PRODUCT_NAME%*" 2^>nul') do (
    set "INSTALL_DIR=%%a"
    echo Found installation directory: !INSTALL_DIR!
)

if not defined INSTALL_DIR (
    echo Trading Dashboard installation not found.
    echo Creating new installation directory...
    
    REM Create installation folder
    set "INSTALL_DIR=C:\Program Files\TrustDan\Trading Dashboard"
    mkdir "!INSTALL_DIR!" 2>nul
)

REM Create a proper launcher script
echo Creating launcher script...
(
    echo @echo off
    echo setlocal enabledelayedexpansion
    echo.
    echo REM Change to the application directory
    echo cd /d "%%~dp0"
    echo.
    echo REM Create log file
    echo echo Starting application at %%date%% %%time%% ^> "%%~dp0trading-dashboard-log.txt"
    echo echo Working directory: %%CD%% ^>^> "%%~dp0trading-dashboard-log.txt"
    echo.
    echo REM List files in directory
    echo dir /b ^>^> "%%~dp0trading-dashboard-log.txt"
    echo.
    echo REM Launch application
    echo echo Launching application... ^>^> "%%~dp0trading-dashboard-log.txt"
    echo "%%~dp0trading-dashboard.exe" ^> "%%~dp0app-output.txt" 2^>^&1
    echo.
    echo REM Check exit code
    echo if %%ERRORLEVEL%% NEQ 0 (
    echo     echo Application exited with code %%ERRORLEVEL%% ^>^> "%%~dp0trading-dashboard-log.txt"
    echo ^) else (
    echo     echo Application exited normally ^>^> "%%~dp0trading-dashboard-log.txt"
    echo ^)
    echo.
    echo endlocal
) > "%~dp0run-trading-dashboard.bat"

echo Copying files to installation directory...
REM Copy launcher script to installation folder
copy /Y "%~dp0run-trading-dashboard.bat" "!INSTALL_DIR!\run-trading-dashboard.bat"

echo Creating shortcut...
REM Create a shortcut to the launcher on desktop
powershell -Command "$WshShell = New-Object -ComObject WScript.Shell; $Shortcut = $WshShell.CreateShortcut([Environment]::GetFolderPath('Desktop') + '\Trading Dashboard.lnk'); $Shortcut.TargetPath = '!INSTALL_DIR!\run-trading-dashboard.bat'; $Shortcut.IconLocation = '!INSTALL_DIR!\trading-dashboard.exe,0'; $Shortcut.WorkingDirectory = '!INSTALL_DIR!'; $Shortcut.Save()"

echo Done. Try running the Trading Dashboard application using the desktop shortcut or from:
echo !INSTALL_DIR!\run-trading-dashboard.bat

endlocal 