// src/cmd/theme-manager/main.go
// Main entry point for the NextUI Theme Manager application

package main

import (
	"fmt"
	"nextui-themes/internal/app"
	"nextui-themes/internal/logging"
	"nextui-themes/internal/ui/screens"
	"os"
	"path/filepath"
	"runtime"
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

		// Ensure screen value is valid - update the upper bound to include all new screens
		if currentScreen < app.Screens.MainMenu || currentScreen > 500 {
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

		// New screens for streamlined UI
		case app.Screens.ThemesMenu:
			logging.LogDebug("Showing themes menu")
			selection, exitCode = screens.ThemesMenuScreen()
			nextScreen = screens.HandleThemesMenu(selection, exitCode)
			logging.LogDebug("Themes menu returned next screen: %d", nextScreen)

		case app.Screens.ComponentsMenu:
			logging.LogDebug("Showing components menu")
			selection, exitCode = screens.ComponentsMenuScreen()
			nextScreen = screens.HandleComponentsMenu(selection, exitCode)
			logging.LogDebug("Components menu returned next screen: %d", nextScreen)

		// Theme browse/download screens
		case app.Screens.ThemeBrowseMenu:
			logging.LogDebug("Showing theme browse menu")
			selection, exitCode = screens.ThemeBrowseMenuScreen()
			nextScreen = screens.HandleThemeBrowseMenu(selection, exitCode)
			logging.LogDebug("Theme browse menu returned next screen: %d", nextScreen)

		case app.Screens.ThemeDownloadMenu:
			logging.LogDebug("Showing theme download menu")
			selection, exitCode = screens.ThemeDownloadMenuScreen()
			nextScreen = screens.HandleThemeDownloadMenu(selection, exitCode)
			logging.LogDebug("Theme download menu returned next screen: %d", nextScreen)

		case app.Screens.ThemeDownloadConfirm:
			logging.LogDebug("Processing theme download")
			selection, exitCode = screens.ThemeDownloadConfirmScreen()
			// Always return to Themes Menu after download
			nextScreen = app.Screens.ThemesMenu
			logging.LogDebug("Theme download confirmation completed")

		// Component type screens
		case app.Screens.WallpapersMenu:
			logging.LogDebug("Showing wallpapers menu")
			selection, exitCode = screens.WallpapersMenuScreen()
			nextScreen = screens.HandleWallpapersMenu(selection, exitCode)
			logging.LogDebug("Wallpapers menu returned next screen: %d", nextScreen)

		case app.Screens.IconsMenuNew:
			logging.LogDebug("Showing icons menu")
			selection, exitCode = screens.IconsMenuNew()
			nextScreen = screens.HandleIconsMenuNew(selection, exitCode)
			logging.LogDebug("Icons menu returned next screen: %d", nextScreen)

		case app.Screens.AccentsMenu:
			logging.LogDebug("Showing accents menu")
			selection, exitCode = screens.AccentsMenuScreen()
			nextScreen = screens.HandleAccentsMenu(selection, exitCode)
			logging.LogDebug("Accents menu returned next screen: %d", nextScreen)

		case app.Screens.LEDsMenu:
			logging.LogDebug("Showing LEDs menu")
			selection, exitCode = screens.LEDsMenuScreen()
			nextScreen = screens.HandleLEDsMenu(selection, exitCode)
			logging.LogDebug("LEDs menu returned next screen: %d", nextScreen)

		case app.Screens.FontsMenu:
			logging.LogDebug("Showing fonts menu")
			selection, exitCode = screens.FontsMenuScreen()
			nextScreen = screens.HandleFontsMenu(selection, exitCode)
			logging.LogDebug("Fonts menu returned next screen: %d", nextScreen)

		// Component options screens
		case app.Screens.ComponentBrowseMenu:
			logging.LogDebug("Showing component browse menu")
			selection, exitCode = screens.ComponentBrowseMenuScreen()
			nextScreen = screens.HandleComponentBrowseMenu(selection, exitCode)
			logging.LogDebug("Component browse menu returned next screen: %d", nextScreen)

		case app.Screens.ComponentDownloadMenu:
			logging.LogDebug("Showing component download menu")
			selection, exitCode = screens.ComponentDownloadMenuScreen()
			nextScreen = screens.HandleComponentDownloadMenu(selection, exitCode)
			logging.LogDebug("Component download menu returned next screen: %d", nextScreen)

		case app.Screens.ComponentDownloadConfirm:
			logging.LogDebug("Processing component download")
			selection, exitCode = screens.ComponentDownloadConfirmScreen()
			// Return to component context menu after confirmation
			componentType := app.GetComponentContext()
			if componentType == app.ComponentTypeWallpaper {
				nextScreen = app.Screens.WallpapersMenu
			} else if componentType == app.ComponentTypeIcon {
				nextScreen = app.Screens.IconsMenuNew
			} else if componentType == app.ComponentTypeAccent {
				nextScreen = app.Screens.AccentsMenu
			} else if componentType == app.ComponentTypeLED {
				nextScreen = app.Screens.LEDsMenu
			} else if componentType == app.ComponentTypeFont {
				nextScreen = app.Screens.FontsMenu
			} else {
				nextScreen = app.Screens.ComponentsMenu
			}
			logging.LogDebug("Component download confirmation completed")

		case app.Screens.ComponentExportMenu:
			logging.LogDebug("Showing component export menu")
			selection, exitCode = screens.ComponentExportMenuScreen()
			nextScreen = screens.HandleComponentExportMenu(selection, exitCode)
			logging.LogDebug("Component export menu returned next screen: %d", nextScreen)

		// Replace the old ThemeApplyMenu with a redirect to ThemeBrowseMenu
		case app.Screens.ThemeApplyMenu:
			logging.LogDebug("Redirecting from old theme apply menu to new theme browse menu")
			app.SetCurrentScreen(app.Screens.ThemeBrowseMenu)
			continue

		case app.Screens.ThemeExtractMenu:
			logging.LogDebug("Showing theme extract menu")
			selection, exitCode = screens.ThemeExtractMenuScreen()
			nextScreen = screens.HandleThemeExtractMenu(selection, exitCode)
			logging.LogDebug("Theme extract menu returned next screen: %d", nextScreen)

		case app.Screens.ComponentApplyMenu:
			logging.LogDebug("Redirecting from old component apply menu to component browse menu")
			// Set a default component type if none is set
			if app.GetComponentContext() == 0 {
				app.SetComponentContext(app.ComponentTypeWallpaper)
			}
			app.SetCurrentScreen(app.Screens.ComponentBrowseMenu)
			continue

		// For backward compatibility - redirect to new Components menu
		case app.Screens.CustomizationMenu:
			logging.LogDebug("Redirecting from old customization menu to new Components menu")
			app.SetCurrentScreen(app.Screens.ComponentsMenu)
			continue

		// Add GlobalOptionsMenu case
		case app.Screens.GlobalOptionsMenu:
			logging.LogDebug("Showing global options menu")
			selection, exitCode = screens.GlobalOptionsMenuScreen()
			nextScreen = screens.HandleGlobalOptionsMenu(selection, exitCode)
			logging.LogDebug("Global options menu returned next screen: %d", nextScreen)

		// Add SystemOptionsMenu case
		case app.Screens.SystemOptionsMenu:
			logging.LogDebug("Showing system options menu")
			selection, exitCode = screens.SystemOptionsMenuScreen()
			nextScreen = screens.HandleSystemOptionsMenu(selection, exitCode)
			logging.LogDebug("System options menu returned next screen: %d", nextScreen)

		// Add SystemOptionsForSelectedSystem case
		case app.Screens.SystemOptionsForSelectedSystem:
			logging.LogDebug("Showing system options for selected system")
			selection, exitCode = screens.SystemOptionsForSelectedSystemScreen()
			nextScreen = screens.HandleSystemOptionsForSelectedSystem(selection, exitCode)
			logging.LogDebug("System options for selected system returned next screen: %d", nextScreen)

		// Add WallpaperSelection case
		case app.Screens.WallpaperSelection:
			logging.LogDebug("Showing wallpaper selection")
			selection, exitCode = screens.WallpaperSelectionScreen()
			nextScreen = screens.HandleWallpaperSelection(selection, exitCode)
			logging.LogDebug("Wallpaper selection returned next screen: %d", nextScreen)

		// Add WallpaperConfirm case
		case app.Screens.WallpaperConfirm:
			logging.LogDebug("Showing wallpaper confirmation")
			selection, exitCode = screens.WallpaperConfirmScreen()
			nextScreen = screens.HandleWallpaperConfirm(selection, exitCode)
			logging.LogDebug("Wallpaper confirmation returned next screen: %d", nextScreen)

		// Add FontSelection case
		case app.Screens.FontSelection:
			logging.LogDebug("Showing font selection")
			selection, exitCode = screens.FontSelectionScreen()
			nextScreen = screens.HandleFontSelection(selection, exitCode)
			logging.LogDebug("Font selection returned next screen: %d", nextScreen)

		// Add FontList case
		case app.Screens.FontList:
			logging.LogDebug("Showing font list")
			selection, exitCode = screens.FontListScreen()
			nextScreen = screens.HandleFontList(selection, exitCode)
			logging.LogDebug("Font list returned next screen: %d", nextScreen)

		// Add FontPreview case
		case app.Screens.FontPreview:
			logging.LogDebug("Showing font preview")
			selection, exitCode = screens.FontPreviewScreen()
			nextScreen = screens.HandleFontPreview(selection, exitCode)
			logging.LogDebug("Font preview returned next screen: %d", nextScreen)

		// Add AccentMenu case
		case app.Screens.AccentMenu:
			logging.LogDebug("Showing accent menu")
			selection, exitCode = screens.AccentMenuScreen()
			nextScreen = screens.HandleAccentMenu(selection, exitCode)
			logging.LogDebug("Accent menu returned next screen: %d", nextScreen)

		// Add AccentSelection case
		case app.Screens.AccentSelection:
			logging.LogDebug("Showing accent selection")
			selection, exitCode = screens.AccentSelectionScreen()
			nextScreen = screens.HandleAccentSelection(selection, exitCode)
			logging.LogDebug("Accent selection returned next screen: %d", nextScreen)

		// Add AccentExport case
		case app.Screens.AccentExport:
			logging.LogDebug("Showing accent export")
			selection, exitCode = screens.AccentExportScreen()
			nextScreen = screens.HandleAccentExport(selection, exitCode)
			logging.LogDebug("Accent export returned next screen: %d", nextScreen)

		// Add LEDMenu case
		case app.Screens.LEDMenu:
			logging.LogDebug("Showing LED menu")
			selection, exitCode = screens.LEDMenuScreen()
			nextScreen = screens.HandleLEDMenu(selection, exitCode)
			logging.LogDebug("LED menu returned next screen: %d", nextScreen)

		// Add LEDSelection case
		case app.Screens.LEDSelection:
			logging.LogDebug("Showing LED selection")
			selection, exitCode = screens.LEDSelectionScreen()
			nextScreen = screens.HandleLEDSelection(selection, exitCode)
			logging.LogDebug("LED selection returned next screen: %d", nextScreen)

		// Add LEDExport case
		case app.Screens.LEDExport:
			logging.LogDebug("Showing LED export")
			selection, exitCode = screens.LEDExportScreen()
			nextScreen = screens.HandleLEDExport(selection, exitCode)
			logging.LogDebug("LED export returned next screen: %d", nextScreen)

		// Add IconSelection case
		case app.Screens.IconSelection:
			logging.LogDebug("Showing icon selection")
			selection, exitCode = screens.IconSelectionScreen()
			nextScreen = screens.HandleIconSelection(selection, exitCode)
			logging.LogDebug("Icon selection returned next screen: %d", nextScreen)

		// Add IconConfirm case
		case app.Screens.IconConfirm:
			logging.LogDebug("Showing icon confirmation")
			selection, exitCode = screens.IconConfirmScreen()
			nextScreen = screens.HandleIconConfirm(selection, exitCode)
			logging.LogDebug("Icon confirmation returned next screen: %d", nextScreen)

		// Add ClearIconsConfirm case
		case app.Screens.ClearIconsConfirm:
			logging.LogDebug("Showing clear icons confirmation")
			selection, exitCode = screens.ClearIconsConfirmScreen()
			nextScreen = screens.HandleClearIconsConfirm(selection, exitCode)
			logging.LogDebug("Clear icons confirmation returned next screen: %d", nextScreen)

		// Add SystemIconSelection case
		case app.Screens.SystemIconSelection:
			logging.LogDebug("Showing system icon selection")
			selection, exitCode = screens.SystemIconSelectionScreen()
			nextScreen = screens.HandleSystemIconSelection(selection, exitCode)
			logging.LogDebug("System icon selection returned next screen: %d", nextScreen)

		// Add SystemIconConfirm case
		case app.Screens.SystemIconConfirm:
			logging.LogDebug("Showing system icon confirmation")
			selection, exitCode = screens.SystemIconConfirmScreen()
			nextScreen = screens.HandleSystemIconConfirm(selection, exitCode)
			logging.LogDebug("System icon confirmation returned next screen: %d", nextScreen)

		// Add ThemeSelection case
		case app.Screens.ThemeSelection:
			logging.LogDebug("Showing theme selection")
			selection, exitCode = screens.ThemeSelectionScreen()
			nextScreen = screens.HandleThemeSelection(selection, exitCode)
			logging.LogDebug("Theme selection returned next screen: %d", nextScreen)

		// Add DefaultThemeOptions case
		case app.Screens.DefaultThemeOptions:
			logging.LogDebug("Showing default theme options")
			selection, exitCode = screens.DefaultThemeOptionsScreen()
			nextScreen = screens.HandleDefaultThemeOptions(selection, exitCode)
			logging.LogDebug("Default theme options returned next screen: %d", nextScreen)

		// Add ConfirmScreen case
		case app.Screens.ConfirmScreen:
			logging.LogDebug("Showing confirmation screen")
			selection, exitCode = screens.ConfirmScreen()
			nextScreen = screens.HandleConfirmScreen(selection, exitCode)
			logging.LogDebug("Confirmation screen returned next screen: %d", nextScreen)

		case app.Screens.ResetMenu:
			logging.LogDebug("Showing reset menu")
			selection, exitCode = screens.ResetMenuScreen()
			nextScreen = screens.HandleResetMenu(selection, exitCode)
			logging.LogDebug("Reset menu returned next screen: %d", nextScreen)

		// For backward compatibility - redirect to new Themes menu
		case app.Screens.ThemeManagementMenu:
			logging.LogDebug("Redirecting from old theme management menu to new Themes menu")
			app.SetCurrentScreen(app.Screens.ThemesMenu)
			continue

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
			logging.LogDebug("Showing theme export confirm screen")
			selection, exitCode = screens.ThemeExportScreen()
			nextScreen = screens.HandleThemeExport(selection, exitCode)
			logging.LogDebug("Theme export confirm returned next screen: %d", nextScreen)

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

		default:
			logging.LogDebug("CRITICAL ERROR: Invalid screen value: %d, resetting to MainMenu", currentScreen)
			app.SetCurrentScreen(app.Screens.MainMenu)
			continue
		}

		// Add extra debug logging
		logging.LogDebug("Current screen: %d, Next screen: %d", currentScreen, nextScreen)

		// Updated range check for valid screen values
        if nextScreen < app.Screens.MainMenu || nextScreen > 500 {
			logging.LogDebug("ERROR: Invalid next screen value: %d, defaulting to MainMenu", nextScreen)
			nextScreen = app.Screens.MainMenu
		}

		// Update the current screen - add extra debugging
		logging.LogDebug("Setting next screen to: %d", nextScreen)
		app.SetCurrentScreen(nextScreen)
		logging.LogDebug("Screen set to: %d", app.GetCurrentScreen())
	}
}
