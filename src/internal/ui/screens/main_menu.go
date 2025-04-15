// src/internal/ui/screens/main_menu.go
// Implementation of the main menu screen

package screens

import (
	"os"
	"strings"

	"nextui-themes/internal/app"
	"nextui-themes/internal/logging"
	"nextui-themes/internal/ui"
)

// MainMenuScreen displays the main menu options
func MainMenuScreen() (string, int) {
	// Menu items without numbers
	menu := []string{
		"Theme Management", // New consolidated option
		"Customization",
		"Reset",
	}

	return ui.DisplayMinUiList(strings.Join(menu, "\n"), "text", "NextUI Theme Selector", "--cancel-text", "QUIT")
}

// HandleMainMenu processes the user's selection from the main menu
func HandleMainMenu(selection string, exitCode int) app.Screen {
	logging.LogDebug("handleMainMenu called with selection: '%s', exitCode: %d", selection, exitCode)

	switch exitCode {
	case 0:
		// User selected an option
		switch selection {
		case "Theme Management":
			logging.LogDebug("Selected Theme Management")
			return app.Screens.ThemeManagementMenu // New screen

		case "Customization":
			logging.LogDebug("Selected Customization")
			return app.Screens.CustomizationMenu

		case "Reset":
			logging.LogDebug("Selected Reset")
			return app.Screens.ResetMenu

		default:
			logging.LogDebug("Unknown selection: %s", selection)
			return app.Screens.MainMenu
		}

	case 1, 2:
		// User pressed cancel or back
		logging.LogDebug("User cancelled/exited")
		os.Exit(0)
	}

	return app.Screens.MainMenu
}