!macro customInstallShortcuts
  # Create shortcuts to the batch file instead of directly to the exe
  CreateShortcut "$SMPROGRAMS\${PRODUCT_NAME}\${PRODUCT_NAME}.lnk" "$INSTDIR\run-trading-dashboard.bat" "" "$INSTDIR\${PRODUCT_FILENAME}.exe" 0
  CreateShortcut "$DESKTOP\${PRODUCT_NAME}.lnk" "$INSTDIR\run-trading-dashboard.bat" "" "$INSTDIR\${PRODUCT_FILENAME}.exe" 0
!macroend

!macro installActionsBeforeInstall
  # Custom actions to perform before files are installed
!macroend

!macro installActionsAfterInstall
  # Custom actions to perform after files are installed
  !insertmacro customInstall
!macroend

!macro installActionsBeforeShortcuts
  # Custom actions to perform before shortcuts are created
!macroend

!macro installActionsAfterShortcuts
  # Custom actions to perform after shortcuts are created
  !insertmacro customShortcuts
!macroend

!macro uninstallActionsBeforeUninstall
  # Custom actions to perform before files are uninstalled
!macroend

!macro uninstallActionsAfterUninstall
  # Custom actions to perform after files are uninstalled
  !insertmacro customUnInstall
!macroend 