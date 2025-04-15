// src/internal/ui/screens/theme_management_menu.go
// Implementation of the theme management menu screen

package screens

import (
	"strings"

	"nextui-themes/internal/app"
	"nextui-themes/internal/logging"
	"nextui-themes/internal/ui"
)

// ThemeManagementMenuScreen displays the theme management options
func ThemeManagementMenuScreen() (string, int) {
	menu := []string{
		"Import",
		"Export",
		"Convert Theme",
	}

	return ui.DisplayMinUiList(strings.Join(menu, "\n"), "text", "Theme Management")
}

// HandleThemeManagementMenu processes the user's selection from the theme management menu
func HandleThemeManagementMenu(selection string, exitCode int) app.Screen {
	logging.LogDebug("HandleThemeManagementMenu called with selection: '%s', exitCode: %d", selection, exitCode)

	switch exitCode {
	case 0:
		switch selection {
		case "Import":
			logging.LogDebug("Selected Import")
			return app.Screens.ThemeImportTypeMenu

		case "Export":
			logging.LogDebug("Selected Export")
			return app.Screens.ThemeExportTypeMenu

		case "Convert Theme":
			logging.LogDebug("Selected Convert Theme")
			return app.Screens.ThemeConvertSelection

		default:
			logging.LogDebug("Unknown selection: %s", selection)
			return app.Screens.ThemeManagementMenu
		}

	case 1, 2:
		// User pressed cancel or back
		return app.Screens.MainMenu
	}

	return app.Screens.ThemeManagementMenu
}