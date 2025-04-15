// src/internal/ui/screens/icons_menu.go
// Implementation of the icons menu screen

package screens

import (
	"fmt"
	"strings"

	"nextui-themes/internal/app"
	"nextui-themes/internal/icons"
	"nextui-themes/internal/logging"
	"nextui-themes/internal/ui"
)

// IconsMenuScreen displays the icons menu options
func IconsMenuScreen() (string, int) {
	// Menu items
	menu := []string{
		"Icon Packs",
		"Clear Icons",
	}

	return ui.DisplayMinUiList(strings.Join(menu, "\n"), "text", "Icons Menu")
}

// HandleIconsMenu processes the user's selection from the icons menu
func HandleIconsMenu(selection string, exitCode int) app.Screen {
	logging.LogDebug("HandleIconsMenu called with selection: '%s', exitCode: %d", selection, exitCode)

	switch exitCode {
	case 0:
		// User selected an option
		switch selection {
		case "Icon Packs":
			logging.LogDebug("Selected Icon Packs")
			return app.Screens.IconSelection

		case "Clear Icons":
			logging.LogDebug("Selected Clear Icons")
			return app.Screens.ClearIconsConfirm

		default:
			logging.LogDebug("Unknown selection: %s", selection)
			return app.Screens.IconsMenuNew
		}

	case 1, 2:
		// User pressed cancel or back
		return app.Screens.ComponentsMenu
	}

	return app.Screens.IconsMenuNew
}

// ClearIconsConfirmScreen displays a confirmation prompt for clearing icons
func ClearIconsConfirmScreen() (string, int) {
	message := "Are you sure you want to remove all system icons?"

	options := []string{
		"Yes",
		"No",
	}

	return ui.DisplayMinUiList(strings.Join(options, "\n"), "text", message)
}

// HandleClearIconsConfirm processes the user's confirmation for clearing icons
func HandleClearIconsConfirm(selection string, exitCode int) app.Screen {
	logging.LogDebug("HandleClearIconsConfirm called with selection: '%s', exitCode: %d", selection, exitCode)

	switch exitCode {
	case 0:
		if selection == "Yes" {
			// Clear all icons
			logging.LogDebug("User confirmed, clearing all icons")
			err := icons.DeleteAllIcons()
			if err != nil {
				logging.LogDebug("Error clearing icons: %v", err)
				ui.ShowMessage(fmt.Sprintf("Error: %s", err), "3")
			} else {
				ui.ShowMessage("All system icons have been cleared", "3")
			}
		}
		// Return to reset menu regardless of Yes/No selection
		return app.Screens.ResetMenu

	case 1, 2:
		// User pressed cancel or back
		logging.LogDebug("User cancelled, returning to reset menu")
		return app.Screens.ResetMenu
	}

	return app.Screens.ResetMenu
}
