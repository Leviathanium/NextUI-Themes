// src/cmd/theme-manager/main.go
// Main entry point for the NextUI Theme Manager application

package main

import (
	"os"
	"path/filepath"
    "runtime"
	"nextui-themes/internal/app"
	"nextui-themes/internal/logging"
	"nextui-themes/internal/ui/screens"
	"fmt"
)

func main() {
	// Recover from panics
	defer func() {
		if r := recover(); r != nil {
			// Get stack trace
			buf := make([]byte, 4096)
			n := runtime.Stack(buf, false)
			stackTrace := string(buf[:n])

			// Log the panic
			fmt.Fprintf(os.Stderr, "PANIC: %v\n\nStack Trace:\n%s\n", r, stackTrace)

			// Also try to log to file if possible
			if logging.IsLoggerInitialized() {
				logging.LogDebug("PANIC: %v\n\nStack Trace:\n%s\n", r, stackTrace)
			}

			// Exit with error
			os.Exit(1)
		}
	}()

	// Initialize the logger
	defer logging.CloseLogger()

	logging.LogDebug("Application started")
	logging.SetLoggerInitialized() // Explicitly mark logger as initialized

	// Get current directory
	cwd, err := os.Getwd()
	if err != nil {
		logging.LogDebug("Error getting current directory: %v", err)
		return
	}

	// Check if minui-list exists in the application directory
	minuiListPath := filepath.Join(cwd, "minui-list")
	_, err = os.Stat(minuiListPath)
	if err != nil {
		logging.LogDebug("minui-list not found at %s: %v", minuiListPath, err)
		return
	}

	// Check if minui-presenter exists in the application directory
	minuiPresenterPath := filepath.Join(cwd, "minui-presenter")
	_, err = os.Stat(minuiPresenterPath)
	if err != nil {
		logging.LogDebug("minui-presenter not found at %s: %v", minuiPresenterPath, err)
		return
	}

	// Initialize application
	if err := app.Initialize(); err != nil {
		logging.LogDebug("Failed to initialize application: %v", err)
		return
	}

	logging.LogDebug("Starting main loop")

    // Main application loop
    for {
        var selection string
        var exitCode int
        var nextScreen app.Screen

        // Log current screen
        currentScreen := app.GetCurrentScreen()
        logging.LogDebug("Current screen: %d", currentScreen)

        // Ensure screen value is valid
        if currentScreen < app.Screens.MainMenu || currentScreen > app.Screens.ThemeConvertConfirm {
            logging.LogDebug("CRITICAL ERROR: Invalid screen value: %d, resetting to MainMenu", currentScreen)
            app.SetCurrentScreen(app.Screens.MainMenu)
            continue
        }

        // Process current screen
        switch currentScreen {
        case app.Screens.MainMenu:
            logging.LogDebug("Showing main menu")
            selection, exitCode = screens.MainMenuScreen()
            nextScreen = screens.HandleMainMenu(selection, exitCode)
            logging.LogDebug("Main menu returned next screen: %d", nextScreen)

        // New theme management screens
        case app.Screens.ThemeManagementMenu:
            logging.LogDebug("Showing theme management menu")
            selection, exitCode = screens.ThemeManagementMenuScreen()
            nextScreen = screens.HandleThemeManagementMenu(selection, exitCode)
            logging.LogDebug("Theme management menu returned next screen: %d", nextScreen)

        // Import screens
        case app.Screens.ThemeImportTypeMenu:
            logging.LogDebug("Showing theme import type menu")
            selection, exitCode = screens.ThemeImportTypeMenuScreen()
            nextScreen = screens.HandleThemeImportTypeMenu(selection, exitCode)
            logging.LogDebug("Theme import type menu returned next screen: %d", nextScreen)

        case app.Screens.ThemeImportSelection:
            logging.LogDebug("Showing theme import selection")
            selection, exitCode = screens.ThemeImportSelectionScreen()
            nextScreen = screens.HandleThemeImportSelection(selection, exitCode)
            logging.LogDebug("Theme import selection returned next screen: %d", nextScreen)

        case app.Screens.ThemeImportOptions:
            logging.LogDebug("Showing theme import options")
            selection, exitCode = screens.ThemeImportOptionsScreen()
            nextScreen = screens.HandleThemeImportOptions(selection, exitCode)
            logging.LogDebug("Theme import options returned next screen: %d", nextScreen)

        case app.Screens.ThemeImportComponentSelection:
            logging.LogDebug("Showing theme import component selection")
            selection, exitCode = screens.ThemeImportComponentSelectionScreen()
            nextScreen = screens.HandleThemeImportComponentSelection(selection, exitCode)
            logging.LogDebug("Theme import component selection returned next screen: %d", nextScreen)

        case app.Screens.ThemeImportConfirm:
            logging.LogDebug("Showing theme import confirmation")
            selection, exitCode = screens.ThemeImportConfirmScreen()
            nextScreen = screens.HandleThemeImportConfirm(selection, exitCode)
            logging.LogDebug("Theme import confirmation returned next screen: %d", nextScreen)

        // Export screens
        case app.Screens.ThemeExportTypeMenu:
            logging.LogDebug("Showing theme export type menu")
            selection, exitCode = screens.ThemeExportTypeMenuScreen()
            nextScreen = screens.HandleThemeExportTypeMenu(selection, exitCode)
            logging.LogDebug("Theme export type menu returned next screen: %d", nextScreen)

        case app.Screens.ThemeExportName:
            logging.LogDebug("Showing theme export name selection")
            selection, exitCode = screens.ThemeExportNameScreen()
            nextScreen = screens.HandleThemeExportName(selection, exitCode)
            logging.LogDebug("Theme export name selection returned next screen: %d", nextScreen)

        case app.Screens.ThemeExportOptions:
            logging.LogDebug("Showing theme export options")
            selection, exitCode = screens.ThemeExportOptionsScreen()
            nextScreen = screens.HandleThemeExportOptions(selection, exitCode)
            logging.LogDebug("Theme export options returned next screen: %d", nextScreen)

        case app.Screens.ThemeExportComponentSelection:
            logging.LogDebug("Showing theme export component selection")
            selection, exitCode = screens.ThemeExportComponentSelectionScreen()
            nextScreen = screens.HandleThemeExportComponentSelection(selection, exitCode)
            logging.LogDebug("Theme export component selection returned next screen: %d", nextScreen)

        case app.Screens.ThemeExportConfirm:
            logging.LogDebug("Showing theme export confirmation")
            selection, exitCode = screens.ThemeExportConfirmScreen()
            nextScreen = screens.HandleThemeExportConfirm(selection, exitCode)
            logging.LogDebug("Theme export confirmation returned next screen: %d", nextScreen)

        // Convert screens
        case app.Screens.ThemeConvertSelection:
            logging.LogDebug("Showing theme convert selection")
            selection, exitCode = screens.ThemeConvertSelectionScreen()
            nextScreen = screens.HandleThemeConvertSelection(selection, exitCode)
            logging.LogDebug("Theme convert selection returned next screen: %d", nextScreen)

        case app.Screens.ThemeConvertOptions:
            logging.LogDebug("Showing theme convert options")
            selection, exitCode = screens.ThemeConvertOptionsScreen()
            nextScreen = screens.HandleThemeConvertOptions(selection, exitCode)
            logging.LogDebug("Theme convert options returned next screen: %d", nextScreen)

        case app.Screens.ThemeConvertComponentSelection:
            logging.LogDebug("Showing theme convert component selection")
            selection, exitCode = screens.ThemeConvertComponentSelectionScreen()
            nextScreen = screens.HandleThemeConvertComponentSelection(selection, exitCode)
            logging.LogDebug("Theme convert component selection returned next screen: %d", nextScreen)

        case app.Screens.ThemeConvertConfirm:
            logging.LogDebug("Showing theme convert confirmation")
            selection, exitCode = screens.ThemeConvertConfirmScreen()
            nextScreen = screens.HandleThemeConvertConfirm(selection, exitCode)
            logging.LogDebug("Theme convert confirmation returned next screen: %d", nextScreen)

        // Legacy screens - maintained for compatibility
        case app.Screens.ThemesMenu:
            logging.LogDebug("Showing themes menu (legacy)")
            selection, exitCode = screens.ThemesMenuScreen()
            nextScreen = screens.HandleThemesMenu(selection, exitCode)
            logging.LogDebug("Themes menu returned next screen: %d", nextScreen)

        case app.Screens.ThemeImport:
            logging.LogDebug("Showing theme import selection (legacy)")
            selection, exitCode = screens.ThemeImportScreen()
            nextScreen = screens.HandleThemeImport(selection, exitCode)
            logging.LogDebug("Theme import returned next screen: %d", nextScreen)

        case app.Screens.ThemeExport:
            logging.LogDebug("Showing theme export screen (legacy)")
            selection, exitCode = screens.ThemeExportScreen()
            nextScreen = screens.HandleThemeExport(selection, exitCode)
            logging.LogDebug("Theme export returned next screen: %d", nextScreen)

		case app.Screens.ThemeSelection:
			logging.LogDebug("Showing theme selection")
			selection, exitCode = screens.ThemeSelectionScreen()
			nextScreen = screens.HandleThemeSelection(selection, exitCode)
            logging.LogDebug("Theme selection returned next screen: %d", nextScreen)

		case app.Screens.ResetMenu:
			logging.LogDebug("Showing reset menu")
			selection, exitCode = screens.ResetMenuScreen()
			nextScreen = screens.HandleResetMenu(selection, exitCode)
            logging.LogDebug("Reset menu returned next screen: %d", nextScreen)

		case app.Screens.ConfirmScreen:
			logging.LogDebug("Showing confirmation screen")
			selection, exitCode = screens.ConfirmScreen()
			nextScreen = screens.HandleConfirmScreen(selection, exitCode)
            logging.LogDebug("Confirm screen returned next screen: %d", nextScreen)

		case app.Screens.FontSelection:
			logging.LogDebug("Showing font slot selection")
			selection, exitCode = screens.FontSelectionScreen()
			nextScreen = screens.HandleFontSelection(selection, exitCode)
            logging.LogDebug("Font selection returned next screen: %d", nextScreen)

		case app.Screens.FontList:
			logging.LogDebug("Showing font list")
			selection, exitCode = screens.FontListScreen()
			nextScreen = screens.HandleFontList(selection, exitCode)
            logging.LogDebug("Font list returned next screen: %d", nextScreen)

		case app.Screens.FontPreview:
			logging.LogDebug("Showing font preview")
			selection, exitCode = screens.FontPreviewScreen()
			nextScreen = screens.HandleFontPreview(selection, exitCode)
            logging.LogDebug("Font preview returned next screen: %d", nextScreen)

		case app.Screens.CustomizationMenu:
			logging.LogDebug("Showing customization menu")
			selection, exitCode = screens.CustomizationMenuScreen()
			nextScreen = screens.HandleCustomizationMenu(selection, exitCode)
            logging.LogDebug("Customization menu returned next screen: %d", nextScreen)

		case app.Screens.GlobalOptionsMenu:
			logging.LogDebug("Showing global options menu")
			selection, exitCode = screens.GlobalOptionsMenuScreen()
			nextScreen = screens.HandleGlobalOptionsMenu(selection, exitCode)
            logging.LogDebug("Global options menu returned next screen: %d", nextScreen)

		case app.Screens.SystemOptionsMenu:
			logging.LogDebug("Showing system options menu")
			selection, exitCode = screens.SystemOptionsMenuScreen()
			nextScreen = screens.HandleSystemOptionsMenu(selection, exitCode)
            logging.LogDebug("System options menu returned next screen: %d", nextScreen)

		case app.Screens.SystemOptionsForSelectedSystem:
			logging.LogDebug("Showing options for selected system")
			selection, exitCode = screens.SystemOptionsForSelectedSystemScreen()
			nextScreen = screens.HandleSystemOptionsForSelectedSystem(selection, exitCode)
            logging.LogDebug("System options for selected system returned next screen: %d", nextScreen)

		case app.Screens.WallpaperSelection:
			logging.LogDebug("Showing wallpaper selection")
			selection, exitCode = screens.WallpaperSelectionScreen()
			nextScreen = screens.HandleWallpaperSelection(selection, exitCode)
            logging.LogDebug("Wallpaper selection returned next screen: %d", nextScreen)

		case app.Screens.AccentMenu:
			logging.LogDebug("Showing accent menu")
			selection, exitCode = screens.AccentMenuScreen()
			nextScreen = screens.HandleAccentMenu(selection, exitCode)
            logging.LogDebug("Accent menu returned next screen: %d", nextScreen)

		case app.Screens.AccentSelection:
			logging.LogDebug("Showing accent selection")
			selection, exitCode = screens.AccentSelectionScreen()
			nextScreen = screens.HandleAccentSelection(selection, exitCode)
            logging.LogDebug("Accent selection returned next screen: %d", nextScreen)

		case app.Screens.AccentExport:
			logging.LogDebug("Handling accent export")
			selection, exitCode = screens.AccentExportScreen()
			nextScreen = screens.HandleAccentExport(selection, exitCode)
            logging.LogDebug("Accent export returned next screen: %d", nextScreen)

		case app.Screens.LEDMenu:
			logging.LogDebug("Showing LED menu")
			selection, exitCode = screens.LEDMenuScreen()
			nextScreen = screens.HandleLEDMenu(selection, exitCode)
            logging.LogDebug("LED menu returned next screen: %d", nextScreen)

		case app.Screens.LEDSelection:
			logging.LogDebug("Showing LED selection")
			selection, exitCode = screens.LEDSelectionScreen()
			nextScreen = screens.HandleLEDSelection(selection, exitCode)
            logging.LogDebug("LED selection returned next screen: %d", nextScreen)

		case app.Screens.LEDExport:
			logging.LogDebug("Handling LED export")
			selection, exitCode = screens.LEDExportScreen()
			nextScreen = screens.HandleLEDExport(selection, exitCode)
            logging.LogDebug("LED export returned next screen: %d", nextScreen)

		case app.Screens.IconsMenu:
			logging.LogDebug("Showing icons menu")
			selection, exitCode = screens.IconsMenuScreen()
			nextScreen = screens.HandleIconsMenu(selection, exitCode)
            logging.LogDebug("Icons menu returned next screen: %d", nextScreen)

		case app.Screens.IconSelection:
			logging.LogDebug("Showing icon selection")
			selection, exitCode = screens.IconSelectionScreen()
			nextScreen = screens.HandleIconSelection(selection, exitCode)
            logging.LogDebug("Icon selection returned next screen: %d", nextScreen)

        case app.Screens.SystemIconSelection:
            logging.LogDebug("Showing system icon selection")
            selection, exitCode = screens.SystemIconSelectionScreen()
            nextScreen = screens.HandleSystemIconSelection(selection, exitCode)
            logging.LogDebug("System icon selection returned next screen: %d", nextScreen)

        case app.Screens.SystemIconConfirm:
            logging.LogDebug("Showing system icon confirmation")
            selection, exitCode = screens.SystemIconConfirmScreen()
            nextScreen = screens.HandleSystemIconConfirm(selection, exitCode)
            logging.LogDebug("System icon confirmation returned next screen: %d", nextScreen)

		case app.Screens.IconConfirm:
			logging.LogDebug("Showing icon confirmation")
			selection, exitCode = screens.IconConfirmScreen()
			nextScreen = screens.HandleIconConfirm(selection, exitCode)
            logging.LogDebug("Icon confirmation returned next screen: %d", nextScreen)

		case app.Screens.ClearIconsConfirm:
			logging.LogDebug("Showing clear icons confirmation")
			selection, exitCode = screens.ClearIconsConfirmScreen()
			nextScreen = screens.HandleClearIconsConfirm(selection, exitCode)
            logging.LogDebug("Clear icons confirmation returned next screen: %d", nextScreen)

        default:
            logging.LogDebug("Unknown screen type: %d, defaulting to MainMenu", currentScreen)
            nextScreen = app.Screens.MainMenu
        }

        // Add extra debug logging
        logging.LogDebug("Current screen: %d, Next screen: %d", currentScreen, nextScreen)

        // Updated range check for valid screen values
        if nextScreen < app.Screens.MainMenu || nextScreen > app.Screens.ThemeConvertConfirm {
            logging.LogDebug("ERROR: Invalid next screen value: %d, defaulting to MainMenu", nextScreen)
            nextScreen = app.Screens.MainMenu
        }

        // Update the current screen - add extra debugging
        logging.LogDebug("Setting next screen to: %d", nextScreen)
        app.SetCurrentScreen(nextScreen)
        logging.LogDebug("Screen set to: %d", app.GetCurrentScreen())
	}
}