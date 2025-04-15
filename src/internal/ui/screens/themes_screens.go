// src/internal/ui/screens/themes_screens.go
// Implementation of screens for theme management

package screens

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"nextui-themes/internal/app"
	"nextui-themes/internal/logging"
	"nextui-themes/internal/themes"
	"nextui-themes/internal/ui"
)

// ThemesMenuScreen displays the main themes menu
func ThemesMenuScreen() (string, int) {
	// Menu items
	menu := []string{
		"Import Theme",
		"Export Current Settings",
	}

	return ui.DisplayMinUiList(strings.Join(menu, "\n"), "text", "Themes")
}

// Replace the HandleThemesMenu function with this:
func HandleThemesMenu(selection string, exitCode int) app.Screen {
	logging.LogDebug("HandleThemesMenu called with selection: '%s', exitCode: %d", selection, exitCode)

	var nextScreen app.Screen

	switch exitCode {
	case 0:
		// User selected an option
		switch selection {
		case "Import Theme":
			logging.LogDebug("Selected Import Theme")
			nextScreen = app.Screens.ThemeImportTypeMenu

		case "Export Current Settings":
			logging.LogDebug("Selected Export Current Settings")
			nextScreen = app.Screens.ThemeExportTypeMenu

		default:
			logging.LogDebug("Unknown selection: %s", selection)
			nextScreen = app.Screens.ThemeManagementMenu
		}

	case 1, 2:
		// User pressed cancel or back
		nextScreen = app.Screens.MainMenu
	default:
		// Default case
		nextScreen = app.Screens.ThemeManagementMenu
	}

	logging.LogDebug("HandleThemesMenu returning screen: %d", nextScreen)
	return nextScreen
}

// ThemeImportScreen displays available themes from the Imports directory
func ThemeImportScreen() (string, int) {
	// Get the current directory
	cwd, err := os.Getwd()
	if err != nil {
		logging.LogDebug("Error getting current directory: %v", err)
		ui.ShowMessage(fmt.Sprintf("Error: %s", err), "3")
		return "", 1
	}

	// Path to Themes/Imports directory
	importsDir := filepath.Join(cwd, "Themes", "Imports")

	// Ensure directory exists
	if err := os.MkdirAll(importsDir, 0755); err != nil {
		logging.LogDebug("Error creating imports directory: %v", err)
		ui.ShowMessage(fmt.Sprintf("Error: %s", err), "3")
		return "", 1
	}

	// List available themes
	entries, err := os.ReadDir(importsDir)
	if err != nil {
		logging.LogDebug("Error reading imports directory: %v", err)
		ui.ShowMessage(fmt.Sprintf("Error: %s", err), "3")
		return "", 1
	}

	// Filter for theme directories
	var themesList []string
	for _, entry := range entries {
		if entry.IsDir() && strings.HasSuffix(entry.Name(), ".theme") {
			themesList = append(themesList, entry.Name())
		}
	}

	if len(themesList) == 0 {
		logging.LogDebug("No themes found")
		ui.ShowMessage("No themes found in Imports directory", "3")
		return "", 1
	}

	logging.LogDebug("Found %d themes", len(themesList))
	return ui.DisplayMinUiList(strings.Join(themesList, "\n"), "text", "Select Theme to Import")
}

// HandleThemeImport processes the user's theme selection
func HandleThemeImport(selection string, exitCode int) app.Screen {
	logging.LogDebug("HandleThemeImport called with selection: '%s', exitCode: %d", selection, exitCode)

	switch exitCode {
	case 0:
		// User selected a theme
		app.SetSelectedTheme(selection)
		return app.Screens.ThemeImportConfirm

	case 1, 2:
		// User pressed cancel or back
		return app.Screens.ThemeManagementMenu
	}

	return app.Screens.ThemeImportSelection
}

// ThemeExportScreen displays the theme export confirmation
func ThemeExportScreen() (string, int) {
	// Simple confirmation message
	message := "Export current theme settings?\nThis will create a theme package in Themes/Exports."
	options := []string{
		"Yes",
		"No",
	}

	return ui.DisplayMinUiList(strings.Join(options, "\n"), "text", message)
}

// HandleThemeExport processes the user's choice to export a theme
func HandleThemeExport(selection string, exitCode int) app.Screen {
	logging.LogDebug("HandleThemeExport called with selection: '%s', exitCode: %d", selection, exitCode)

	switch exitCode {
	case 0:
		if selection == "Yes" {
			// Perform theme export
			if err := themes.ExportTheme(); err != nil {
				logging.LogDebug("Error exporting theme: %v", err)
				ui.ShowMessage(fmt.Sprintf("Error: %s", err), "3")
			}
		}
		// Return to themes menu
		return app.Screens.ThemeManagementMenu

	case 1, 2:
		// User pressed cancel or back
		return app.Screens.ThemeManagementMenu
	}

	return app.Screens.ThemeExportTypeMenu
}
