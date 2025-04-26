!macro customInstall
  # Set CGO_ENABLED for this installation
  WriteRegStr HKCU "Environment" "CGO_ENABLED" "1"
  
  # Include the SQLite DLL
  File "sqlite3.dll"
  
  # Create a launcher batch file that sets CGO_ENABLED for running the application
  FileOpen $0 "$INSTDIR\run-trading-dashboard.bat" w
  FileWrite $0 "@echo off$\r$\n"
  FileWrite $0 "setlocal enabledelayedexpansion$\r$\n"
  FileWrite $0 "set CGO_ENABLED=1$\r$\n"
  FileWrite $0 "cd /d `"%~dp0`"$\r$\n"
  FileWrite $0 "echo Starting application at %date% %time% > `"%~dp0\trading-dashboard-log.txt`"$\r$\n"
  FileWrite $0 "echo CGO_ENABLED=%CGO_ENABLED% >> `"%~dp0\trading-dashboard-log.txt`"$\r$\n"
  FileWrite $0 "echo Working directory: %CD% >> `"%~dp0\trading-dashboard-log.txt`"$\r$\n"
  FileWrite $0 "dir /b `"%~dp0`" >> `"%~dp0\trading-dashboard-log.txt`"$\r$\n"
  FileWrite $0 "start `"`" `"%~dp0trading-dashboard.exe`" > `"%~dp0\app-output.txt`" 2>&1$\r$\n"
  FileWrite $0 "endlocal$\r$\n"
  FileClose $0
  
  # Create a direct executable launcher
  FileOpen $0 "$INSTDIR\run-direct.bat" w
  FileWrite $0 "@echo off$\r$\n"
  FileWrite $0 "set CGO_ENABLED=1$\r$\n"
  FileWrite $0 "cd /d `"%~dp0`"$\r$\n"
  FileWrite $0 "`"%~dp0trading-dashboard.exe`"$\r$\n"
  FileClose $0
!macroend

!macro customShortcuts
  # Create shortcuts to the batch file instead of directly to the exe
  CreateShortcut "$SMPROGRAMS\${PRODUCT_NAME}\${PRODUCT_NAME}.lnk" "$INSTDIR\run-trading-dashboard.bat" "" "$INSTDIR\${PRODUCT_FILENAME}.exe" 0
  CreateShortcut "$DESKTOP\${PRODUCT_NAME}.lnk" "$INSTDIR\run-trading-dashboard.bat" "" "$INSTDIR\${PRODUCT_FILENAME}.exe" 0
  
  # Create a direct launcher shortcut
  CreateShortcut "$SMPROGRAMS\${PRODUCT_NAME}\${PRODUCT_NAME} (Direct).lnk" "$INSTDIR\run-direct.bat" "" "$INSTDIR\${PRODUCT_FILENAME}.exe" 0
  CreateShortcut "$DESKTOP\${PRODUCT_NAME} (Direct).lnk" "$INSTDIR\run-direct.bat" "" "$INSTDIR\${PRODUCT_FILENAME}.exe" 0
!macroend

!macro customUnInstall
  # Clean up the environment variable we set
  DeleteRegValue HKCU "Environment" "CGO_ENABLED"
  
  # Remove our additional files
  Delete "$INSTDIR\run-trading-dashboard.bat"
  Delete "$INSTDIR\run-direct.bat"
  Delete "$INSTDIR\trading-dashboard-log.txt"
  Delete "$INSTDIR\app-output.txt"
  Delete "$INSTDIR\sqlite3.dll"
!macroend 