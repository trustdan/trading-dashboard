@echo off
cd /d "%~dp0"
set CGO_ENABLED=1
start "" "trading-dashboard.exe"