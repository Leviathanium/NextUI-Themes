// src/internal/ui/screens/theme_import_screens.go
// Implementation of theme import screens for different component types

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

// ThemeImportTypeMenuScreen displays import type selection menu
func ThemeImportTypeMenuScreen() (string, int) {
	menu := []string{
		"Full Theme (.theme)",
		"Accent Pack (.acc)",
		"LED Pack (.led)",
		"Wallpaper Pack (.bg)",
		"Icon Pack (.icon)",
		"Font Pack (.font)",
	}

	return ui.DisplayMinUiList(strings.Join(menu, "\n"), "text", "Select Import Type")
}

// HandleThemeImportTypeMenu processes the user's import type selection
func HandleThemeImportTypeMenu(selection string, exitCode int) app.Screen {
	logging.LogDebug("HandleThemeImportTypeMenu called with selection: '%s', exitCode: %d", selection, exitCode)

	switch exitCode {
	case 0:
		// Set the component type in application state
		switch selection {
		case "Full Theme (.theme)":
			app.SetImportComponentType(app.ComponentTypeFullTheme)
			return app.Screens.ThemeImportSelection

		case "Accent Pack (.acc)":
			app.SetImportComponentType(app.ComponentTypeAccent)
			return app.Screens.ThemeImportSelection

		case "LED Pack (.led)":
			app.SetImportComponentType(app.ComponentTypeLED)
			return app.Screens.ThemeImportSelection

		case "Wallpaper Pack (.bg)":
			app.SetImportComponentType(app.ComponentTypeWallpaper)
			return app.Screens.ThemeImportSelection

		case "Icon Pack (.icon)":
			app.SetImportComponentType(app.ComponentTypeIcon)
			return app.Screens.ThemeImportSelection

		case "Font Pack (.font)":
			app.SetImportComponentType(app.ComponentTypeFont)
			return app.Screens.ThemeImportSelection

		default:
			logging.LogDebug("Unknown selection: %s", selection)
			return app.Screens.ThemeImportTypeMenu
		}

	case 1, 2:
		// User pressed cancel or back
		return app.Screens.ThemeManagementMenu
	}

	return app.Screens.ThemeImportTypeMenu
}

// getImportDirectory returns the appropriate directory for the selected component type
func getImportDirectory() (string, error) {
	// Get current directory
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("error getting current directory: %w", err)
	}

	// Determine directory based on component type
	componentType := app.GetImportComponentType()
	var importDir string

	switch componentType {
	case app.ComponentTypeFullTheme:
		importDir = filepath.Join(cwd, "Themes", "Imports")
	case app.ComponentTypeAccent:
		importDir = filepath.Join(cwd, "Accents", "Imports")
	case app.ComponentTypeLED:
		importDir = filepath.Join(cwd, "LEDs", "Imports")
	case app.ComponentTypeWallpaper:
		importDir = filepath.Join(cwd, "Wallpapers", "Imports")
	case app.ComponentTypeIcon:
		importDir = filepath.Join(cwd, "Icons", "Imports")
	case app.ComponentTypeFont:
		importDir = filepath.Join(cwd, "Fonts")
	default:
		return "", fmt.Errorf("invalid component type: %d", componentType)
	}

	// Ensure directory exists
	if err := os.MkdirAll(importDir, 0755); err != nil {
		return "", fmt.Errorf("error creating import directory: %w", err)
	}

	return importDir, nil
}

// getImportFileExtension returns the file extension for the selected component type
func getImportFileExtension() string {
	componentType := app.GetImportComponentType()

	switch componentType {
	case app.ComponentTypeFullTheme:
		return ".theme"
	case app.ComponentTypeAccent:
		return ".acc"
	case app.ComponentTypeLED:
		return ".led"
	case app.ComponentTypeWallpaper:
		return ".bg"
	case app.ComponentTypeIcon:
		return ".icon"
	case app.ComponentTypeFont:
		return ".font"
	default:
		return ""
	}
}

// getImportTitle returns an appropriate title for the import screen
func getImportTitle() string {
	componentType := app.GetImportComponentType()

	switch componentType {
	case app.ComponentTypeFullTheme:
		return "Select Theme to Import"
	case app.ComponentTypeAccent:
		return "Select Accent Pack to Import"
	case app.ComponentTypeLED:
		return "Select LED Pack to Import"
	case app.ComponentTypeWallpaper:
		return "Select Wallpaper Pack to Import"
	case app.ComponentTypeIcon:
		return "Select Icon Pack to Import"
	case app.ComponentTypeFont:
		return "Select Font Pack to Import"
	default:
		return "Select Item to Import"
	}
}

// ThemeImportSelectionScreen displays available items to import
func ThemeImportSelectionScreen() (string, int) {
	// Get the appropriate import directory
	importDir, err := getImportDirectory()
	if err != nil {
		logging.LogDebug("Error getting import directory: %v", err)
		ui.ShowMessage(fmt.Sprintf("Error: %s", err), "3")
		return "", 1
	}

	// Get the file extension for the current component type
	fileExt := getImportFileExtension()

	// Read the directory
	entries, err := os.ReadDir(importDir)
	if err != nil {
		logging.LogDebug("Error reading import directory: %v", err)
		ui.ShowMessage(fmt.Sprintf("Error reading directory: %s", err), "3")
		return "", 1
	}

	// Filter for directories with the appropriate extension
	var items []string
	for _, entry := range entries {
		if entry.IsDir() {
			name := entry.Name()
			if strings.HasSuffix(name, fileExt) {
				items = append(items, name)
			}
		}
	}

	// If no items found
	if len(items) == 0 {
		logging.LogDebug("No items found in import directory: %s", importDir)
		ui.ShowMessage(fmt.Sprintf("No items found in %s", importDir), "3")
		return "", 1
	}

	// Display the list
	title := getImportTitle()
	return ui.DisplayMinUiList(strings.Join(items, "\n"), "text", title)
}

// HandleThemeImportSelection processes the user's import item selection
func HandleThemeImportSelection(selection string, exitCode int) app.Screen {
	logging.LogDebug("HandleThemeImportSelection called with selection: '%s', exitCode: %d", selection, exitCode)

	switch exitCode {
	case 0:
		// Store the selected item
		app.SetSelectedImportItem(selection)
		return app.Screens.ThemeImportOptions

	case 1, 2:
		// User pressed cancel or back
		return app.Screens.ThemeImportTypeMenu
	}

	return app.Screens.ThemeImportSelection
}

// ThemeImportOptionsScreen displays import options for the selected item
func ThemeImportOptionsScreen() (string, int) {
	componentType := app.GetImportComponentType()
	itemName := app.GetSelectedImportItem()

	// For full themes, we offer component selection
	if componentType == app.ComponentTypeFullTheme {
		menu := []string{
			"Import All Components",
			"Select Components to Import",
		}

		title := fmt.Sprintf("Import Options for %s", itemName)
		return ui.DisplayMinUiList(strings.Join(menu, "\n"), "text", title)
	}

	// For other component types, we just show a confirmation
	menu := []string{
		"Import",
		"Cancel",
	}

	var typeName string
	switch componentType {
	case app.ComponentTypeAccent:
		typeName = "Accent Pack"
	case app.ComponentTypeLED:
		typeName = "LED Pack"
	case app.ComponentTypeWallpaper:
		typeName = "Wallpaper Pack"
	case app.ComponentTypeIcon:
		typeName = "Icon Pack"
	case app.ComponentTypeFont:
		typeName = "Font Pack"
	default:
		typeName = "Item"
	}

	title := fmt.Sprintf("Import %s: %s", typeName, itemName)
	return ui.DisplayMinUiList(strings.Join(menu, "\n"), "text", title)
}

// HandleThemeImportOptions processes the user's import options selection
func HandleThemeImportOptions(selection string, exitCode int) app.Screen {
	logging.LogDebug("HandleThemeImportOptions called with selection: '%s', exitCode: %d", selection, exitCode)

	componentType := app.GetImportComponentType()

	switch exitCode {
	case 0:
		if componentType == app.ComponentTypeFullTheme {
			// Full theme options
			if selection == "Import All Components" {
				app.SetImportAllComponents(true)
				return app.Screens.ThemeImportConfirm
			} else if selection == "Select Components to Import" {
				app.SetImportAllComponents(false)
				return app.Screens.ThemeImportComponentSelection
			}
		} else {
			// Other component types
			if selection == "Import" {
				return app.Screens.ThemeImportConfirm
			}
		}
		return app.Screens.ThemeImportTypeMenu

	case 1, 2:
		// User pressed cancel or back
		return app.Screens.ThemeImportSelection
	}

	return app.Screens.ThemeImportOptions
}

// ThemeImportComponentSelectionScreen displays component selection for full theme import
func ThemeImportComponentSelectionScreen() (string, int) {
	menu := []string{
		"Wallpapers",
		"Icons",
		"Accents",
		"LEDs",
		"Fonts",
		"Continue with Selected Components",
	}

	title := "Select Components to Import"
	return ui.DisplayMinUiList(strings.Join(menu, "\n"), "text", title, "--header-bar-toggles")
}

// HandleThemeImportComponentSelection processes component selection for import
func HandleThemeImportComponentSelection(selection string, exitCode int) app.Screen {
	logging.LogDebug("HandleThemeImportComponentSelection called with selection: '%s', exitCode: %d", selection, exitCode)

	switch exitCode {
	case 0:
		if selection == "Continue with Selected Components" {
			return app.Screens.ThemeImportConfirm
		}

		// Toggle the selected component
		switch selection {
		case "Wallpapers":
			app.ToggleImportComponent(app.ComponentTypeWallpaper)
		case "Icons":
			app.ToggleImportComponent(app.ComponentTypeIcon)
		case "Accents":
			app.ToggleImportComponent(app.ComponentTypeAccent)
		case "LEDs":
			app.ToggleImportComponent(app.ComponentTypeLED)
		case "Fonts":
			app.ToggleImportComponent(app.ComponentTypeFont)
		}

		// Stay on this screen to allow multiple selections
		return app.Screens.ThemeImportComponentSelection

	case 1, 2:
		// User pressed cancel or back
		return app.Screens.ThemeImportOptions
	}

	return app.Screens.ThemeImportComponentSelection
}

// ThemeImportConfirmScreen displays final confirmation before import
func ThemeImportConfirmScreen() (string, int) {
	componentType := app.GetImportComponentType()
	itemName := app.GetSelectedImportItem()

	var message string
	if componentType == app.ComponentTypeFullTheme {
		if app.GetImportAllComponents() {
			message = fmt.Sprintf("Import all components from %s?", itemName)
		} else {
			message = fmt.Sprintf("Import selected components from %s?", itemName)
		}
	} else {
		var typeName string
		switch componentType {
		case app.ComponentTypeAccent:
			typeName = "accent pack"
		case app.ComponentTypeLED:
			typeName = "LED pack"
		case app.ComponentTypeWallpaper:
			typeName = "wallpaper pack"
		case app.ComponentTypeIcon:
			typeName = "icon pack"
		case app.ComponentTypeFont:
			typeName = "font pack"
		default:
			typeName = "item"
		}

		message = fmt.Sprintf("Import %s: %s?", typeName, itemName)
	}

	menu := []string{
		"Yes",
		"No",
	}

	return ui.DisplayMinUiList(strings.Join(menu, "\n"), "text", message)
}

// HandleThemeImportConfirm processes the final import confirmation
func HandleThemeImportConfirm(selection string, exitCode int) app.Screen {
	logging.LogDebug("HandleThemeImportConfirm called with selection: '%s', exitCode: %d", selection, exitCode)

	switch exitCode {
	case 0:
		if selection == "Yes" {
			// Perform the actual import based on component type
			err := performImport()
			if err != nil {
				logging.LogDebug("Import error: %v", err)
				ui.ShowMessage(fmt.Sprintf("Import error: %s", err), "3")
			} else {
				ui.ShowMessage("Import completed successfully", "3")
			}
		}
		return app.Screens.ThemeManagementMenu

	case 1, 2:
		// User pressed cancel or back
		if app.GetImportComponentType() == app.ComponentTypeFullTheme && !app.GetImportAllComponents() {
			return app.Screens.ThemeImportComponentSelection
		}
		return app.Screens.ThemeImportOptions
	}

	return app.Screens.ThemeImportConfirm
}

// performImport executes the actual import operation based on current settings
func performImport() error {
	// Get import parameters from app state
	componentType := themes.ComponentType(app.GetImportComponentType())
	itemName := app.GetSelectedImportItem()

	// Check if we're importing selected components or all components
	var selectedComponents map[themes.ComponentType]bool

	if app.GetImportAllComponents() || componentType != themes.ComponentTypeFullTheme {
		// For full themes with all components, or for non-theme components, use an empty map
		selectedComponents = make(map[themes.ComponentType]bool)
	} else {
		// Convert app component types to themes component types
		selectedComponents = make(map[themes.ComponentType]bool)
		for compType := range app.GetSelectedImportComponents() {
			selectedComponents[themes.ComponentType(compType)] = true
		}
	}

	logging.LogDebug("Import parameters: Type=%d, Item=%s, AllComponents=%v",
		componentType, itemName, app.GetImportAllComponents())

	// Execute the import operation in the themes package
	return themes.PerformImport(componentType, itemName, selectedComponents)
}
