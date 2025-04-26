@echo off
set CGO_ENABLED=1
cd /d "%~dp0\build\bin"
echo Running trading dashboard with CGO_ENABLED=%CGO_ENABLED%
echo Working directory: %CD%
"trading-dashboard.exe"
echo Exit code: %ERRORLEVEL%
pause 