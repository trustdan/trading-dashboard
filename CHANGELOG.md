# Changelog

## v1.0.1 (2025-04-25)

### Fixed
- **Critical Fix:** Changed database location to store in user's AppData directory instead of application directory
  - Resolves permission issues requiring admin mode to run
  - Database now stored in `%APPDATA%\TradingDashboard\data\trading.db` on Windows
  - All data will now be preserved between application restarts
- **Improved Build Process:** Enhanced the build script
  - Added `--debug` flag to build with developer tools enabled
  - Added better error detection for CGO issues
- **Documentation:** Updated README with troubleshooting information
  - Added details about database locations
  - Added instructions for debugging UI issues

### Scripts Added
- `direct-run.bat` - Simple script to run the application with proper environment variables
- `run-debug.bat` - Runs application in debug mode with detailed logging
- `fix-installed-app.bat` - Script to fix existing installations

### Developer Notes
- If UI components (Stock Rating, Calendar, etc.) are still freezing:
  1. Build with debug tools: `.\build-trading.bat --debug`
  2. Run with debugging: `.\run-debug.bat`
  3. When the app opens, press F12 or Ctrl+Shift+I to open the developer console
  4. Check for any JavaScript errors in the console when clicking on the problematic UI elements 