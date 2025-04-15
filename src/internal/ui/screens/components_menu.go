// src/internal/ui/screens/components_menu.go
// Implementation of the components menu and related screens

package screens

import (
	"fmt"
	"strings"

	"nextui-themes/internal/app"
	"nextui-themes/internal/logging"
	"nextui-themes/internal/ui"
)

// ComponentsMenuScreen displays the component types directly
func ComponentsMenuScreen() (string, int) {
	// Show component types directly
	menu := []string{
		"Wallpapers",
		"Icons",
		"Accents",
		"LEDs",
		"Fonts",
	}

	return ui.DisplayMinUiList(strings.Join(menu, "\n"), "text", "Components")
}

// HandleComponentsMenu processes the user's component type selection
func HandleComponentsMenu(selection string, exitCode int) app.Screen {
	logging.LogDebug("HandleComponentsMenu called with selection: '%s', exitCode: %d", selection, exitCode)

	switch exitCode {
	case 0:
		// Set the selected component type
		switch selection {
		case "Wallpapers":
			app.SetComponentContext(app.ComponentTypeWallpaper)
			return app.Screens.WallpapersMenu
		case "Icons":
			app.SetComponentContext(app.ComponentTypeIcon)
			return app.Screens.IconsMenuNew
		case "Accents":
			app.SetComponentContext(app.ComponentTypeAccent)
			return app.Screens.AccentsMenu
		case "LEDs":
			app.SetComponentContext(app.ComponentTypeLED)
			return app.Screens.LEDsMenu
		case "Fonts":
			app.SetComponentContext(app.ComponentTypeFont)
			return app.Screens.FontsMenu
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

// Helper function to get component type name
func getComponentTypeName() string {
	componentType := app.GetComponentContext()
	switch componentType {
	case app.ComponentTypeWallpaper:
		return "Wallpaper"
	case app.ComponentTypeIcon:
		return "Icon"
	case app.ComponentTypeAccent:
		return "Accent"
	case app.ComponentTypeLED:
		return "LED"
	case app.ComponentTypeFont:
		return "Font"
	default:
		return "Component"
	}
}

// WallpapersMenuScreen displays the options for wallpapers
func WallpapersMenuScreen() (string, int) {
	return ComponentOptionsMenuScreen()
}

// HandleWallpapersMenu handles wallpaper menu selections
func HandleWallpapersMenu(selection string, exitCode int) app.Screen {
	return HandleComponentOptionsMenu(selection, exitCode)
}

// IconsMenuNew displays the options for icons
func IconsMenuNew() (string, int) {
	return ComponentOptionsMenuScreen()
}

// HandleIconsMenuNew handles icon menu selections
func HandleIconsMenuNew(selection string, exitCode int) app.Screen {
	return HandleComponentOptionsMenu(selection, exitCode)
}

// AccentsMenuScreen displays the options for accents
func AccentsMenuScreen() (string, int) {
	return ComponentOptionsMenuScreen()
}

// HandleAccentsMenu handles accent menu selections
func HandleAccentsMenu(selection string, exitCode int) app.Screen {
	return HandleComponentOptionsMenu(selection, exitCode)
}

// LEDsMenuScreen displays the options for LEDs
func LEDsMenuScreen() (string, int) {
	return ComponentOptionsMenuScreen()
}

// HandleLEDsMenu handles LED menu selections
func HandleLEDsMenu(selection string, exitCode int) app.Screen {
	return HandleComponentOptionsMenu(selection, exitCode)
}

// FontsMenuScreen displays the options for fonts
func FontsMenuScreen() (string, int) {
	return ComponentOptionsMenuScreen()
}

// HandleFontsMenu handles font menu selections
func HandleFontsMenu(selection string, exitCode int) app.Screen {
	return HandleComponentOptionsMenu(selection, exitCode)
}

// ComponentOptionsMenuScreen displays the component options (Browse, Download, Export)
func ComponentOptionsMenuScreen() (string, int) {
	menu := []string{
		"Browse",
		"Download",
		"Export",
	}

	componentTypeName := getComponentTypeName()
	title := fmt.Sprintf("%s Options", componentTypeName)

	return ui.DisplayMinUiList(strings.Join(menu, "\n"), "text", title)
}

// HandleComponentOptionsMenu processes the component options selection
func HandleComponentOptionsMenu(selection string, exitCode int) app.Screen {
	logging.LogDebug("HandleComponentOptionsMenu called with selection: '%s', exitCode: %d", selection, exitCode)

	componentType := app.GetComponentContext()

	switch exitCode {
	case 0:
		switch selection {
		case "Browse":
			// Set up for browse - reuse import functionality
			app.SetImportComponentType(componentType)
			return app.Screens.ComponentBrowseMenu

		case "Download":
			// Set up for download
			return app.Screens.ComponentDownloadMenu

		case "Export":
			// Set up for export
			app.SetExportComponentType(componentType)
			return app.Screens.ComponentExportMenu

		default:
			// Stay in the current menu if invalid selection
			return getCurrentComponentTypeMenu(componentType)
		}

	case 1, 2:
		// User pressed cancel or back
		return app.Screens.ComponentsMenu
	}

	return getCurrentComponentTypeMenu(componentType)
}

// Helper function to get the current component type menu
func getCurrentComponentTypeMenu(componentType app.ComponentType) app.Screen {
	switch componentType {
	case app.ComponentTypeWallpaper:
		return app.Screens.WallpapersMenu
	case app.ComponentTypeIcon:
		return app.Screens.IconsMenuNew
	case app.ComponentTypeAccent:
		return app.Screens.AccentsMenu
	case app.ComponentTypeLED:
		return app.Screens.LEDsMenu
	case app.ComponentTypeFont:
		return app.Screens.FontsMenu
	default:
		return app.Screens.ComponentsMenu
	}
}

// ComponentBrowseMenuScreen displays available components to browse
func ComponentBrowseMenuScreen() (string, int) {
	// Reuse the theme import selection logic but with component-specific paths
	componentType := app.GetComponentContext()

	// Set the import component type to match the current context
	app.SetImportComponentType(componentType)

	// Use the theme import selection screen directly
	return ThemeImportSelectionScreen()
}

// HandleComponentBrowseMenu processes component browse selection
func HandleComponentBrowseMenu(selection string, exitCode int) app.Screen {
	logging.LogDebug("HandleComponentBrowseMenu called with selection: '%s', exitCode: %d", selection, exitCode)

	componentType := app.GetComponentContext()

	switch exitCode {
	case 0:
		// User selected a component
		app.SetSelectedImportItem(selection)
		// Set import parameters
		app.SetImportComponentType(componentType)
		app.SetImportAllComponents(true)

		// Go directly to confirmation
		return app.Screens.ThemeImportConfirm

	case 1, 2:
		// User pressed cancel or back
		return getCurrentComponentTypeMenu(componentType)
	}

	return app.Screens.ComponentBrowseMenu
}

// ComponentDownloadMenuScreen displays available components to download (not implemented)
func ComponentDownloadMenuScreen() (string, int) {
	componentTypeName := getComponentTypeName()

	// Example components to download
	menu := []string{
		fmt.Sprintf("Classic %s Pack", componentTypeName),
		fmt.Sprintf("Modern %s Pack", componentTypeName),
		fmt.Sprintf("Retro %s Pack", componentTypeName),
		fmt.Sprintf("Dark %s Pack", componentTypeName),
		fmt.Sprintf("Light %s Pack", componentTypeName),
	}

	return ui.DisplayMinUiList(strings.Join(menu, "\n"), "text", fmt.Sprintf("Download %s Pack", componentTypeName))
}

// HandleComponentDownloadMenu processes component download selection
func HandleComponentDownloadMenu(selection string, exitCode int) app.Screen {
	logging.LogDebug("HandleComponentDownloadMenu called with selection: '%s', exitCode: %d", selection, exitCode)

	componentType := app.GetComponentContext()

	switch exitCode {
	case 0:
		// User selected a component to download
		app.SetSelectedItem(selection)
		return app.Screens.ComponentDownloadConfirm

	case 1, 2:
		// User pressed cancel or back
		return getCurrentComponentTypeMenu(componentType)
	}

	return app.Screens.ComponentDownloadMenu
}

// ComponentDownloadConfirmScreen simulates downloading a component (not implemented)
func ComponentDownloadConfirmScreen() (string, int) {
	componentTypeName := getComponentTypeName()
	selectedItem := app.GetSelectedItem()
	message := fmt.Sprintf("Downloading %s: %s", componentTypeName, selectedItem)

	// Simulate downloading
	ui.ShowMessage(message, "3")
	ui.ShowMessage("Download complete!", "3")

	return "", 1 // Return with exit code 1 to go back
}

// ComponentExportMenuScreen confirms exporting the current component
func ComponentExportMenuScreen() (string, int) {
	componentTypeName := getComponentTypeName()

	// Simple confirmation message
	message := fmt.Sprintf("Export current %s settings?\nThis will create a component package in %ss/Exports.",
		componentTypeName, componentTypeName)

	options := []string{
		"Yes",
		"No",
	}

	return ui.DisplayMinUiList(strings.Join(options, "\n"), "text", message)
}

// HandleComponentExportMenu processes the component export confirmation
func HandleComponentExportMenu(selection string, exitCode int) app.Screen {
	logging.LogDebug("HandleComponentExportMenu called with selection: '%s', exitCode: %d", selection, exitCode)

	componentType := app.GetComponentContext()

	switch exitCode {
	case 0:
		if selection == "Yes" {
			// Set component type for export
			app.SetExportComponentType(componentType)

			// Prompt for export name
			return app.Screens.ThemeExportName
		}
		// Return to component options menu
		return getCurrentComponentTypeMenu(componentType)

	case 1, 2:
		// User pressed cancel or back
		return getCurrentComponentTypeMenu(componentType)
	}

	return app.Screens.ComponentExportMenu
}

// SetupComponentDownloadUI sets up the component context and prepares for download
func SetupComponentDownloadUI() app.Screen {
	// Set the component context if not already set
	if app.GetComponentContext() == 0 {
		app.SetComponentContext(app.ComponentTypeWallpaper) // Default to wallpaper
	}

	return app.Screens.ComponentDownloadMenu
}
