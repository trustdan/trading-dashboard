@echo off
setlocal enabledelayedexpansion

REM Check if sqlite3.dll exists in the current directory
if not exist "%~dp0\sqlite3.dll" (
    echo SQLite3 DLL not found. Downloading...
    
    REM Download SQLite3 DLL from the official site
    powershell -Command "Invoke-WebRequest -Uri 'https://www.sqlite.org/2023/sqlite-dll-win64-x64-3420000.zip' -OutFile '%TEMP%\sqlite3.zip'"
    
    REM Check if download was successful
    if %ERRORLEVEL% NEQ 0 (
        echo Failed to download SQLite3 DLL. Exiting.
        exit /b 1
    ) else (
        REM Extract the DLL
        powershell -Command "Expand-Archive -Path '%TEMP%\sqlite3.zip' -DestinationPath '%TEMP%\sqlite3' -Force"
        
        REM Copy the DLL to the current directory
        copy "%TEMP%\sqlite3\sqlite3.dll" "%~dp0\sqlite3.dll"
        
        REM Clean up temporary files
        del "%TEMP%\sqlite3.zip"
        rmdir /s /q "%TEMP%\sqlite3"
        
        echo SQLite3 DLL downloaded successfully.
    )
)

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
    echo REM Set environment variable
    echo set CGO_ENABLED=1
    echo.
    echo REM Change to the application directory
    echo cd /d "%%~dp0"
    echo.
    echo REM Create log file
    echo echo Starting application at %%date%% %%time%% ^> "%%~dp0trading-dashboard-log.txt"
    echo echo CGO_ENABLED=%%CGO_ENABLED%% ^>^> "%%~dp0trading-dashboard-log.txt"
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
REM Copy SQLite DLL to installation folder
copy /Y "%~dp0sqlite3.dll" "!INSTALL_DIR!\sqlite3.dll"

REM Copy launcher script to installation folder
copy /Y "%~dp0run-trading-dashboard.bat" "!INSTALL_DIR!\run-trading-dashboard.bat"

echo Creating shortcut...
REM Create a shortcut to the launcher on desktop
powershell -Command "$WshShell = New-Object -ComObject WScript.Shell; $Shortcut = $WshShell.CreateShortcut([Environment]::GetFolderPath('Desktop') + '\Trading Dashboard.lnk'); $Shortcut.TargetPath = '!INSTALL_DIR!\run-trading-dashboard.bat'; $Shortcut.IconLocation = '!INSTALL_DIR!\trading-dashboard.exe,0'; $Shortcut.WorkingDirectory = '!INSTALL_DIR!'; $Shortcut.Save()"

echo Done. Try running the Trading Dashboard application using the desktop shortcut or from:
echo !INSTALL_DIR!\run-trading-dashboard.bat

endlocal 