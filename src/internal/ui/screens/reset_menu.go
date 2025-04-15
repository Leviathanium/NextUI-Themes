// src/internal/ui/screens/reset_menu.go
// Implementation of the reset menu screen

package screens

import (
	"strings"

	"nextui-themes/internal/app"
	"nextui-themes/internal/logging"
	"nextui-themes/internal/ui"
)

// ResetMenuScreen displays options for resetting various aspects of the UI
func ResetMenuScreen() (string, int) {
	options := []string{
		"Reset Everything",
		"Reset Wallpapers Only",
		"Reset Icons Only",
	}

	return ui.DisplayMinUiList(strings.Join(options, "\n"), "text", "Reset Options")
}

// HandleResetMenu processes user selection for reset options
func HandleResetMenu(selection string, exitCode int) app.Screen {
	logging.LogDebug("HandleResetMenu called with selection: '%s', exitCode: %d", selection, exitCode)

	switch exitCode {
	case 0:
		if selection == "Reset Everything" {
			logging.LogDebug("Selected to reset everything")
			// Set up confirmation for resetting all components
			// This will be implemented as part of the theme system reset
			return app.Screens.ConfirmScreen
		} else if selection == "Reset Wallpapers Only" {
			logging.LogDebug("Selected to delete all backgrounds")
			app.SetSelectedThemeType(app.DefaultTheme)
			app.SetDefaultAction(app.DeleteAction)
			return app.Screens.ConfirmScreen
		} else if selection == "Reset Icons Only" {
			logging.LogDebug("Selected to delete all icons")
			return app.Screens.ClearIconsConfirm
		}
	case 1, 2:
		return app.Screens.MainMenu
	}

	return app.Screens.ResetMenu
}
