!ifndef CONSTANTS_NSH
!define CONSTANTS_NSH

# Windows message constants
!define WM_SETTINGCHANGE 0x001A

# Define a unique name for the installer
!define /date TIMESTAMP "%Y%m%d-%H%M%S"
!define UNIQUE_INSTALLER_NAME "trading-dashboard-amd64-installer-${TIMESTAMP}.exe"

!endif 