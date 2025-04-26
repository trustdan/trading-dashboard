# Trading Dashboard

A Wails-based trading dashboard application with SQLite database functionality.

## Build Requirements

This application uses SQLite through the go-sqlite3 package which requires CGO to be enabled during compilation. This means you need:

1. Go programming language
2. Wails framework
3. A C compiler (gcc on Windows, clang on macOS)
4. The sqlite3.dll file (on Windows) in the application directory

## Building on Windows

### Prerequisites

1. Install Go from [golang.org](https://golang.org/dl/)
2. Install Wails CLI: `go install github.com/wailsapp/wails/v2/cmd/wails@latest`
3. Install a GCC compiler. Options include:
   - MSYS2 (recommended): Install from [msys2.org](https://www.msys2.org/) and run `pacman -S mingw-w64-x86_64-toolchain`
   - TDM-GCC: Install from [jmeubank.github.io/tdm-gcc](https://jmeubank.github.io/tdm-gcc/download/)
4. Install NSIS (if you want to create installers): Install from [nsis.sourceforge.io](https://nsis.sourceforge.io/Download)

### Build Steps

To build the application on Windows, you need to enable CGO and ensure GCC is in your PATH:

```powershell
# Set CGO_ENABLED environment variable
$env:CGO_ENABLED=1

# Add GCC to your path (adjust path if necessary)
$env:PATH += ";C:\msys64\ucrt64\bin"  # For MSYS2
# Or
$env:PATH += ";C:\TDM-GCC-64\bin"     # For TDM-GCC

# Build the application
wails build
```

### Creating an NSIS Installer

To build an installer that users can download and run:

```powershell
$env:CGO_ENABLED=1
$env:PATH += ";C:\msys64\ucrt64\bin"  # Adjust path to your compiler
wails build --nsis
```

The installer will be created at `build\bin\trading-dashboard-amd64-installer.exe`.

### Using the Build Script

For convenience, a `build-trading.bat` script is included in the repository. This script:

1. Checks for and downloads the required SQLite DLL if needed
2. Sets the correct environment variables (CGO_ENABLED=1, CC=gcc)
3. Builds the application with NSIS
4. Verifies that CGO was properly enabled in the build

To use it, simply run:

```batch
# For production build with installer
.\build-trading.bat

# For debug build with developer tools
.\build-trading.bat --debug
```

## Running the Application

### Important: CGO_ENABLED and SQLite3.dll

This application requires two critical components to run correctly:

1. The application must have been built with `CGO_ENABLED=1`
2. The `sqlite3.dll` file must be in the same directory as the executable

If the application fails to start or you encounter database-related errors, these are likely the issues.

### Data Storage Location

The application stores its database file in the following location:

- Windows: `%APPDATA%\TradingDashboard\data\trading.db`
- macOS: `~/Library/Application Support/TradingDashboard/data/trading.db`
- Linux: `~/.config/TradingDashboard/data/trading.db`

This allows the application to run without requiring administrator privileges. The application will automatically create this directory and database file on first run.

### Using the Debug Scripts

Several utility scripts are included to help run and debug the application:

#### 1. direct-run.bat

This simple script sets `CGO_ENABLED=1` and runs the application from the build directory:

```batch
.\direct-run.bat
```

#### 2. run-debug.bat

This script provides more detailed debugging information:

- Sets debug environment variables 
- Ensures the SQLite DLL is present
- Creates detailed log files
- Opens the application in a separate window with console output

```batch
.\run-debug.bat
```

#### 3. fix-installed-app.bat

If you've already installed the application and it's not running correctly, this script will fix common issues:

- Locates your installation directory
- Downloads and adds the SQLite3 DLL
- Creates a proper launcher script with CGO_ENABLED=1
- Creates a desktop shortcut

```batch
.\fix-installed-app.bat
```

## Building on macOS

### Prerequisites

1. Install Go from [golang.org](https://golang.org/dl/)
2. Install Wails CLI: `go install github.com/wailsapp/wails/v2/cmd/wails@latest`
3. Install Xcode Command Line Tools: `xcode-select --install`
4. (Optional) Install create-dmg: `brew install create-dmg`

### Build Steps

Building on macOS with CGO is simpler as the C compiler is included with Xcode Command Line Tools:

```bash
# Set CGO_ENABLED environment variable
export CGO_ENABLED=1

# Build the application
wails build
```

### Creating a DMG Installer

Wails provides built-in support for creating macOS .app bundles. To create a DMG installer:

#### Option 1: Using Wails Built-in Packaging

Wails natively supports creating macOS application bundles:

```bash
export CGO_ENABLED=1
wails build -platform darwin/universal
```

This creates an application bundle in `build/bin/trading-dashboard.app`.

#### Option 2: Creating a DMG with create-dmg

To create a DMG installer from the app bundle:

```bash
# First build the app bundle
export CGO_ENABLED=1
wails build -platform darwin/universal

# Navigate to the build directory
cd build/bin

# Create a DMG (requires create-dmg, install with brew install create-dmg)
create-dmg \
  --volname "Trading Dashboard Installer" \
  --volicon "path/to/your/icon.icns" \
  --window-pos 200 120 \
  --window-size 800 400 \
  --icon-size 100 \
  --icon "trading-dashboard.app" 200 190 \
  --hide-extension "trading-dashboard.app" \
  --app-drop-link 600 185 \
  "TradingDashboard.dmg" \
  "trading-dashboard.app"
```

## Troubleshooting

### Common Issues

1. **"Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work"**
   - This error means your application was compiled without CGO support, which is required for SQLite
   - Solution: 
     - Make sure to set `CGO_ENABLED=1` before building
     - Use the provided `build-trading.bat` script
     - If an existing installation has this issue, run `fix-installed-app.bat`

2. **"cgo: C compiler "gcc" not found"**
   - Solution: Ensure your C compiler is installed and in your PATH
   - On Windows, install MSYS2/MinGW or TDM-GCC
   - Make sure the compiler is in your PATH (e.g., C:\msys64\ucrt64\bin)

3. **"The application starts but doesn't function correctly"**
   - Make sure the SQLite3.dll file is in the same directory as the executable
   - Use the provided `run-debug.bat` script to run with debugging enabled
   - Check the generated log files for more information

4. **"UI components freezing or getting stuck"**
   - This may be related to UI JavaScript errors
   - Build in debug mode: `.\build-trading.bat --debug`
   - Run with debug tools: `.\run-debug.bat`
   - Open the developer console in the app with F12 or Ctrl+Shift+I

5. **"Calendar or other UI components not functioning"**
   - Try building with debug mode: `wails build -debug -devtools`
   - Run with debug logging: `$env:WAILS_LOG_LEVEL="DEBUG"; .\build\bin\trading-dashboard.exe`

## Development Tips

For easier builds, create a build script:

### Windows (build.bat)
```batch
@echo off
set CGO_ENABLED=1
set PATH=%PATH%;C:\msys64\ucrt64\bin
wails build %*
```

### macOS (build.sh)
```bash
#!/bin/bash
export CGO_ENABLED=1
wails build "$@"
```

Make the script executable: `chmod +x build.sh`
