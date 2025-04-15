// src/internal/ui/screens/components_menu.go
// Implementation of the components menu and related screens

package screens

import (
	"strings"

	"nextui-themes/internal/app"
	"nextui-themes/internal/logging"
	"nextui-themes/internal/ui"
)

// ComponentsMenuScreen displays the components menu options
func ComponentsMenuScreen() (string, int) {
	// Menu items
	menu := []string{
		"Apply Component",
		"Export Component",
	}

	return ui.DisplayMinUiList(strings.Join(menu, "\n"), "text", "Components")
}

// HandleComponentsMenu processes the user's selection from the components menu
func HandleComponentsMenu(selection string, exitCode int) app.Screen {
	logging.LogDebug("HandleComponentsMenu called with selection: '%s', exitCode: %d", selection, exitCode)

	switch exitCode {
	case 0:
		// User selected an option
		switch selection {
		case "Apply Component":
			logging.LogDebug("Selected Apply Component")
			return app.Screens.ComponentApplyMenu

		case "Export Component":
			logging.LogDebug("Selected Export Component")
			return app.Screens.ThemeExportTypeMenu

		default:
			logging.LogDebug("Unknown selection: %s", selection)
			return app.Screens.ComponentsMenu
		}

	case 1, 2:
		// User pressed cancel or back
		return app.Screens.MainMenu
	}

	return app.Screens.ComponentsMenu
}

// ComponentApplyMenuScreen displays the component type selection for applying components
func ComponentApplyMenuScreen() (string, int) {
	// Component type options
	menu := []string{
		"Wallpaper Pack",
		"Icon Pack",
		"Accent Pack",
		"LED Pack",
		"Font Pack",
	}

	return ui.DisplayMinUiList(strings.Join(menu, "\n"), "text", "Select Component Type")
}

// HandleComponentApplyMenu processes the user's component type selection
func HandleComponentApplyMenu(selection string, exitCode int) app.Screen {
	logging.LogDebug("HandleComponentApplyMenu called with selection: '%s', exitCode: %d", selection, exitCode)

	switch exitCode {
	case 0:
		// Set the component type based on selection
		switch selection {
		case "Wallpaper Pack":
			app.SetImportComponentType(app.ComponentTypeWallpaper)
			return app.Screens.ThemeImportSelection
		case "Icon Pack":
			app.SetImportComponentType(app.ComponentTypeIcon)
			return app.Screens.ThemeImportSelection
		case "Accent Pack":
			app.SetImportComponentType(app.ComponentTypeAccent)
			return app.Screens.ThemeImportSelection
		case "LED Pack":
			app.SetImportComponentType(app.ComponentTypeLED)
			return app.Screens.ThemeImportSelection
		case "Font Pack":
			app.SetImportComponentType(app.ComponentTypeFont)
			return app.Screens.ThemeImportSelection
		default:
			return app.Screens.ComponentApplyMenu
		}

	case 1, 2:
		// User pressed cancel or back
		return app.Screens.ComponentsMenu
	}

	return app.Screens.ComponentApplyMenu
}
