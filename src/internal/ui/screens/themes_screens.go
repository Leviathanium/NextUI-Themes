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
	"nextui-themes/internal/ui"
)

// ThemesMenuScreen displays the streamlined themes menu
func ThemesMenuScreen() (string, int) {
	// Menu items
	menu := []string{
		"Apply Theme",
		"Extract Theme Components",
		"Export Current Theme",
	}

	return ui.DisplayMinUiList(strings.Join(menu, "\n"), "text", "Themes")
}

// HandleThemesMenu processes the user's selection from the themes menu
func HandleThemesMenu(selection string, exitCode int) app.Screen {
	logging.LogDebug("HandleThemesMenu called with selection: '%s', exitCode: %d", selection, exitCode)

	var nextScreen app.Screen

	switch exitCode {
	case 0:
		// User selected an option
		switch selection {
		case "Apply Theme":
			logging.LogDebug("Selected Apply Theme")
			nextScreen = app.Screens.ThemeApplyMenu

		case "Extract Theme Components":
			logging.LogDebug("Selected Extract Theme Components")
			nextScreen = app.Screens.ThemeExtractMenu

		case "Export Current Theme":
			logging.LogDebug("Selected Export Current Theme")
			nextScreen = app.Screens.ThemeExportConfirm

		default:
			logging.LogDebug("Unknown selection: %s", selection)
			nextScreen = app.Screens.ThemesMenu
		}

	case 1, 2:
		// User pressed cancel or back
		nextScreen = app.Screens.MainMenu
	default:
		// Default case
		nextScreen = app.Screens.ThemesMenu
	}

	logging.LogDebug("HandleThemesMenu returning screen: %d", nextScreen)
	return nextScreen
}

// ThemeApplyMenuScreen displays the theme selection for applying themes
func ThemeApplyMenuScreen() (string, int) {
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
	return ui.DisplayMinUiList(strings.Join(themesList, "\n"), "text", "Select Theme to Apply")
}

// HandleThemeApplyMenu processes the user's theme selection for applying
func HandleThemeApplyMenu(selection string, exitCode int) app.Screen {
	logging.LogDebug("HandleThemeApplyMenu called with selection: '%s', exitCode: %d", selection, exitCode)

	switch exitCode {
	case 0:
		// User selected a theme
		app.SetSelectedTheme(selection)
		app.SetSelectedThemeType(app.GlobalTheme)

		// Set all components to be imported
		app.SetImportAllComponents(true)

		// Prepare for confirmation
		return app.Screens.ThemeImportConfirm

	case 1, 2:
		// User pressed cancel or back
		return app.Screens.ThemesMenu
	}

	return app.Screens.ThemeApplyMenu
}

// ThemeExtractMenuScreen displays available themes for component extraction
func ThemeExtractMenuScreen() (string, int) {
	// Re-use the same theme listing logic from ThemeApplyMenuScreen
	return ThemeApplyMenuScreen()
}

// HandleThemeExtractMenu processes the user's theme selection for extraction
func HandleThemeExtractMenu(selection string, exitCode int) app.Screen {
	logging.LogDebug("HandleThemeExtractMenu called with selection: '%s', exitCode: %d", selection, exitCode)

	switch exitCode {
	case 0:
		// User selected a theme
		app.SetSelectedTheme(selection)
		return app.Screens.ThemeConvertComponentSelection

	case 1, 2:
		// User pressed cancel or back
		return app.Screens.ThemesMenu
	}

	return app.Screens.ThemeExtractMenu
}

// ExportCurrentThemeScreen displays the theme export confirmation for current theme
func ExportCurrentThemeScreen() (string, int) {
	// Simple confirmation message
	message := "Export current theme settings?\nThis will create a theme package in Themes/Exports."
	options := []string{
		"Yes",
		"No",
	}

	return ui.DisplayMinUiList(strings.Join(options, "\n"), "text", message)
}

// HandleExportCurrentTheme processes the user's choice to export current theme
func HandleExportCurrentTheme(selection string, exitCode int) app.Screen {
	logging.LogDebug("HandleExportCurrentTheme called with selection: '%s', exitCode: %d", selection, exitCode)

	switch exitCode {
	case 0:
		if selection == "Yes" {
			// Set component type to full theme for export
			app.SetExportComponentType(app.ComponentTypeFullTheme)

			// Prompt for theme name
			return app.Screens.ThemeExportName
		}
		// Return to themes menu
		return app.Screens.ThemesMenu

	case 1, 2:
		// User pressed cancel or back
		return app.Screens.ThemesMenu
	}

	return app.Screens.ThemeExportConfirm
}
