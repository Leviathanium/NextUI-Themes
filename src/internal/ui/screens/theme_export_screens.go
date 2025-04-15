// src/internal/ui/screens/theme_export_screens.go
// Implementation of theme export screens for different component types

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

// ThemeExportTypeMenuScreen displays export type selection menu
func ThemeExportTypeMenuScreen() (string, int) {
	menu := []string{
		"Full Theme (.theme)",
		"Current Accents (.acc)",
		"Current LEDs (.led)",
		"Current Wallpapers (.bg)",
		"Current Icons (.icon)",
		"Current Fonts (.font)",
	}

	return ui.DisplayMinUiList(strings.Join(menu, "\n"), "text", "Select Export Type")
}

// HandleThemeExportTypeMenu processes the user's export type selection
func HandleThemeExportTypeMenu(selection string, exitCode int) app.Screen {
	logging.LogDebug("HandleThemeExportTypeMenu called with selection: '%s', exitCode: %d", selection, exitCode)

	switch exitCode {
	case 0:
		// Set the component type in application state
		switch selection {
		case "Full Theme (.theme)":
			app.SetExportComponentType(app.ComponentTypeFullTheme)
			return app.Screens.ThemeExportName

		case "Current Accents (.acc)":
			app.SetExportComponentType(app.ComponentTypeAccent)
			return app.Screens.ThemeExportName

		case "Current LEDs (.led)":
			app.SetExportComponentType(app.ComponentTypeLED)
			return app.Screens.ThemeExportName

		case "Current Wallpapers (.bg)":
			app.SetExportComponentType(app.ComponentTypeWallpaper)
			return app.Screens.ThemeExportName

		case "Current Icons (.icon)":
			app.SetExportComponentType(app.ComponentTypeIcon)
			return app.Screens.ThemeExportName

		case "Current Fonts (.font)":
			app.SetExportComponentType(app.ComponentTypeFont)
			return app.Screens.ThemeExportName

		default:
			logging.LogDebug("Unknown selection: %s", selection)
			return app.Screens.ThemeExportTypeMenu
		}

	case 1, 2:
		// User pressed cancel or back
		return app.Screens.ThemesMenu
	}

	return app.Screens.ThemeExportTypeMenu
}

// getExportDirectory returns the appropriate directory for the selected component type
func getExportDirectory() (string, error) {
	// Get current directory
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("error getting current directory: %w", err)
	}

	// Determine directory based on component type
	componentType := app.GetExportComponentType()
	var exportDir string

	switch componentType {
	case app.ComponentTypeFullTheme:
		exportDir = filepath.Join(cwd, "Themes", "Exports")
	case app.ComponentTypeAccent:
		exportDir = filepath.Join(cwd, "Accents", "Exports")
	case app.ComponentTypeLED:
		exportDir = filepath.Join(cwd, "LEDs", "Exports")
	case app.ComponentTypeWallpaper:
		exportDir = filepath.Join(cwd, "Wallpapers", "Exports")
	case app.ComponentTypeIcon:
		exportDir = filepath.Join(cwd, "Icons", "Exports")
	case app.ComponentTypeFont:
		exportDir = filepath.Join(cwd, "Fonts", "Exports")
	default:
		return "", fmt.Errorf("invalid component type: %d", componentType)
	}

	// Ensure directory exists
	if err := os.MkdirAll(exportDir, 0755); err != nil {
		return "", fmt.Errorf("error creating export directory: %w", err)
	}

	return exportDir, nil
}

// getExportFileExtension returns the file extension for the selected component type
func getExportFileExtension() string {
	componentType := app.GetExportComponentType()

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

// getExportTitle returns an appropriate title for the export screen
func getExportTitle() string {
	componentType := app.GetExportComponentType()

	switch componentType {
	case app.ComponentTypeFullTheme:
		return "Export Theme"
	case app.ComponentTypeAccent:
		return "Export Accent Pack"
	case app.ComponentTypeLED:
		return "Export LED Pack"
	case app.ComponentTypeWallpaper:
		return "Export Wallpaper Pack"
	case app.ComponentTypeIcon:
		return "Export Icon Pack"
	case app.ComponentTypeFont:
		return "Export Font Pack"
	default:
		return "Export Item"
	}
}

// ThemeExportNameScreen allows the user to name the export
func ThemeExportNameScreen() (string, int) {
	// For simplicity, we'll use the same minui-list to offer sequential names
	// In a real implementation, this could be a text input screen

	// Generate sequential name options
	names := generateExportNameOptions()
	title := getExportTitle()

	return ui.DisplayMinUiList(strings.Join(names, "\n"), "text", title)
}

// generateExportNameOptions creates a list of potential export names
func generateExportNameOptions() []string {
	componentType := app.GetExportComponentType()
	exportDir, err := getExportDirectory()
	if err != nil {
		logging.LogDebug("Error getting export directory: %v", err)
		return []string{"export_1"}
	}

	fileExt := getExportFileExtension()

	// Determine prefix based on component type
	var prefix string
	switch componentType {
	case app.ComponentTypeFullTheme:
		prefix = "theme"
	case app.ComponentTypeAccent:
		prefix = "accent"
	case app.ComponentTypeLED:
		prefix = "led"
	case app.ComponentTypeWallpaper:
		prefix = "wallpaper"
	case app.ComponentTypeIcon:
		prefix = "icon"
	case app.ComponentTypeFont:
		prefix = "font"
	default:
		prefix = "export"
	}

	// Find highest existing number
	highestNum := 0
	entries, err := os.ReadDir(exportDir)
	if err == nil {
		for _, entry := range entries {
			if entry.IsDir() {
				name := entry.Name()
				if strings.HasPrefix(name, prefix+"_") && strings.HasSuffix(name, fileExt) {
					// Extract number from name
					numPart := strings.TrimPrefix(name, prefix+"_")
					numPart = strings.TrimSuffix(numPart, fileExt)

					var num int
					if _, err := fmt.Sscanf(numPart, "%d", &num); err == nil {
						if num > highestNum {
							highestNum = num
						}
					}
				}
			}
		}
	}

	// Generate a few sequential name options
	var names []string
	for i := 1; i <= 3; i++ {
		names = append(names, fmt.Sprintf("%s_%d", prefix, highestNum+i))
	}

	// Add custom option
	names = append(names, "Custom Name...")

	return names
}

// HandleThemeExportName processes the user's export name selection
func HandleThemeExportName(selection string, exitCode int) app.Screen {
	logging.LogDebug("HandleThemeExportName called with selection: '%s', exitCode: %d", selection, exitCode)

	switch exitCode {
	case 0:
		// Store the selected name
		if selection == "Custom Name..." {
			// In a real implementation, this would go to a text input screen
			// For now, we'll just use a default name
			app.SetExportName("custom_export")
		} else {
			app.SetExportName(selection)
		}

		// For full themes, proceed to component selection
		if app.GetExportComponentType() == app.ComponentTypeFullTheme {
			return app.Screens.ThemeExportOptions
		}

		// For other component types, go straight to confirmation
		return app.Screens.ThemeExportConfirm

	case 1, 2:
		// User pressed cancel or back
		return app.Screens.ThemeExportTypeMenu
	}

	return app.Screens.ThemeExportName
}

// ThemeExportOptionsScreen displays options for theme export
func ThemeExportOptionsScreen() (string, int) {
	// This is only used for full themes
	menu := []string{
		"Export All Components",
		"Select Components to Export",
	}

	title := "Export Options"
	return ui.DisplayMinUiList(strings.Join(menu, "\n"), "text", title)
}

// HandleThemeExportOptions processes the user's export options selection
func HandleThemeExportOptions(selection string, exitCode int) app.Screen {
	logging.LogDebug("HandleThemeExportOptions called with selection: '%s', exitCode: %d", selection, exitCode)

	switch exitCode {
	case 0:
		if selection == "Export All Components" {
			app.SetExportAllComponents(true)
			return app.Screens.ThemeExportConfirm
		} else if selection == "Select Components to Export" {
			app.SetExportAllComponents(false)
			return app.Screens.ThemeExportComponentSelection
		}
		return app.Screens.ThemeExportOptions

	case 1, 2:
		// User pressed cancel or back
		return app.Screens.ThemeExportName
	}

	return app.Screens.ThemeExportOptions
}

// ThemeExportComponentSelectionScreen displays component selection for theme export
func ThemeExportComponentSelectionScreen() (string, int) {
	menu := []string{
		"Wallpapers",
		"Icons",
		"Accents",
		"LEDs",
		"Fonts",
		"Continue with Selected Components",
	}

	title := "Select Components to Export"
	return ui.DisplayMinUiList(strings.Join(menu, "\n"), "text", title, "--header-bar-toggles")
}

// HandleThemeExportComponentSelection processes component selection for export
func HandleThemeExportComponentSelection(selection string, exitCode int) app.Screen {
	logging.LogDebug("HandleThemeExportComponentSelection called with selection: '%s', exitCode: %d", selection, exitCode)

	switch exitCode {
	case 0:
		if selection == "Continue with Selected Components" {
			return app.Screens.ThemeExportConfirm
		}

		// Toggle the selected component
		switch selection {
		case "Wallpapers":
			app.ToggleExportComponent(app.ComponentTypeWallpaper)
		case "Icons":
			app.ToggleExportComponent(app.ComponentTypeIcon)
		case "Accents":
			app.ToggleExportComponent(app.ComponentTypeAccent)
		case "LEDs":
			app.ToggleExportComponent(app.ComponentTypeLED)
		case "Fonts":
			app.ToggleExportComponent(app.ComponentTypeFont)
		}

		// Stay on this screen to allow multiple selections
		return app.Screens.ThemeExportComponentSelection

	case 1, 2:
		// User pressed cancel or back
		return app.Screens.ThemeExportOptions
	}

	return app.Screens.ThemeExportComponentSelection
}

// ThemeExportConfirmScreen displays final confirmation before export
func ThemeExportConfirmScreen() (string, int) {
	componentType := app.GetExportComponentType()
	exportName := app.GetExportName()
	fileExt := getExportFileExtension()

	var message string
	if componentType == app.ComponentTypeFullTheme {
		if app.GetExportAllComponents() {
			message = fmt.Sprintf("Export all components to %s%s?", exportName, fileExt)
		} else {
			message = fmt.Sprintf("Export selected components to %s%s?", exportName, fileExt)
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

		message = fmt.Sprintf("Export current %s to %s%s?", typeName, exportName, fileExt)
	}

	menu := []string{
		"Yes",
		"No",
	}

	return ui.DisplayMinUiList(strings.Join(menu, "\n"), "text", message)
}

// HandleThemeExportConfirm processes the final export confirmation
func HandleThemeExportConfirm(selection string, exitCode int) app.Screen {
	logging.LogDebug("HandleThemeExportConfirm called with selection: '%s', exitCode: %d", selection, exitCode)

	switch exitCode {
	case 0:
		if selection == "Yes" {
			// Perform the actual export based on component type
			err := performExport()
			if err != nil {
				logging.LogDebug("Export error: %v", err)
				ui.ShowMessage(fmt.Sprintf("Export error: %s", err), "3")
			} else {
				ui.ShowMessage("Export completed successfully", "3")
			}
		}
		return app.Screens.ThemesMenu

	case 1, 2:
		// User pressed cancel or back
		if app.GetExportComponentType() == app.ComponentTypeFullTheme {
			if app.GetExportAllComponents() {
				return app.Screens.ThemeExportOptions
			} else {
				return app.Screens.ThemeExportComponentSelection
			}
		}
		return app.Screens.ThemeExportName
	}

	return app.Screens.ThemeExportConfirm
}

// performExport executes the actual export operation based on current settings
func performExport() error {
	// Get export parameters from app state
	componentType := themes.ComponentType(app.GetExportComponentType())
	exportName := app.GetExportName()

	// Check if we're exporting selected components or all components
	var selectedComponents map[themes.ComponentType]bool

	if app.GetExportAllComponents() || componentType != themes.ComponentTypeFullTheme {
		// For full themes with all components, or for non-theme components, use an empty map
		selectedComponents = make(map[themes.ComponentType]bool)
	} else {
		// Convert app component types to themes component types
		selectedComponents = make(map[themes.ComponentType]bool)
		for compType := range app.GetSelectedExportComponents() {
			selectedComponents[themes.ComponentType(compType)] = true
		}
	}

	logging.LogDebug("Export parameters: Type=%d, Name=%s, AllComponents=%v",
		componentType, exportName, app.GetExportAllComponents())

	// Execute the export operation in the themes package
	return themes.PerformExport(componentType, exportName, selectedComponents)
}
