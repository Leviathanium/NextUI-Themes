// src/internal/ui/screens/theme_convert_screens.go
// Implementation of screens for converting/deconstructing themes

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

// ThemeConvertSelectionScreen displays available themes for conversion
func ThemeConvertSelectionScreen() (string, int) {
	// Get current directory
	cwd, err := os.Getwd()
	if err != nil {
		logging.LogDebug("Error getting current directory: %v", err)
		ui.ShowMessage(fmt.Sprintf("Error: %s", err), "3")
		return "", 1
	}

	// Path to Themes/Imports directory where full themes are stored
	themesDir := filepath.Join(cwd, "Themes", "Imports")

	// Read the directory
	entries, err := os.ReadDir(themesDir)
	if err != nil {
		logging.LogDebug("Error reading themes directory: %v", err)
		ui.ShowMessage(fmt.Sprintf("Error reading themes directory: %s", err), "3")
		return "", 1
	}

	// Filter for directories with .theme extension
	var themes []string
	for _, entry := range entries {
		if entry.IsDir() {
			name := entry.Name()
			if strings.HasSuffix(name, ".theme") {
				themes = append(themes, name)
			}
		}
	}

	// If no themes found
	if len(themes) == 0 {
		logging.LogDebug("No themes found in: %s", themesDir)
		ui.ShowMessage("No themes found. Import themes first.", "3")
		return "", 1
	}

	return ui.DisplayMinUiList(strings.Join(themes, "\n"), "text", "Select Theme to Convert")
}

// HandleThemeConvertSelection processes theme selection for conversion
func HandleThemeConvertSelection(selection string, exitCode int) app.Screen {
	logging.LogDebug("HandleThemeConvertSelection called with selection: '%s', exitCode: %d", selection, exitCode)

	switch exitCode {
	case 0:
		// Store the selected theme for conversion
		app.SetSelectedConvertTheme(selection)
		return app.Screens.ThemeConvertOptions

	case 1, 2:
		// User pressed cancel or back
		return app.Screens.ThemeManagementMenu
	}

	return app.Screens.ThemeConvertSelection
}

// ThemeConvertOptionsScreen displays options for theme conversion
func ThemeConvertOptionsScreen() (string, int) {
	themeName := app.GetSelectedConvertTheme()

	menu := []string{
		"Deconstruct into Components",
		"Cancel",
	}

	title := fmt.Sprintf("Convert: %s", themeName)
	return ui.DisplayMinUiList(strings.Join(menu, "\n"), "text", title)
}

// HandleThemeConvertOptions processes conversion options selection
func HandleThemeConvertOptions(selection string, exitCode int) app.Screen {
	logging.LogDebug("HandleThemeConvertOptions called with selection: '%s', exitCode: %d", selection, exitCode)

	switch exitCode {
	case 0:
		if selection == "Deconstruct into Components" {
			return app.Screens.ThemeConvertComponentSelection
		}
		return app.Screens.ThemeConvertSelection

	case 1, 2:
		// User pressed cancel or back
		return app.Screens.ThemeConvertSelection
	}

	return app.Screens.ThemeConvertOptions
}

// ThemeConvertComponentSelectionScreen displays component selection for conversion
func ThemeConvertComponentSelectionScreen() (string, int) {
	menu := []string{
		"All Components",
		"Accents Only",
		"LEDs Only",
		"Wallpapers Only",
		"Icons Only",
		"Fonts Only",
	}

	themeName := app.GetSelectedConvertTheme()
	title := fmt.Sprintf("Extract from: %s", themeName)

	return ui.DisplayMinUiList(strings.Join(menu, "\n"), "text", title)
}

// HandleThemeConvertComponentSelection processes component selection for conversion
func HandleThemeConvertComponentSelection(selection string, exitCode int) app.Screen {
	logging.LogDebug("HandleThemeConvertComponentSelection called with selection: '%s', exitCode: %d", selection, exitCode)

	switch exitCode {
	case 0:
		// Set the components to extract based on selection
		switch selection {
		case "All Components":
			app.SetConvertAllComponents(true)
			app.ClearSelectedConvertComponents()
			app.AddConvertComponent(app.ComponentTypeAccent)
			app.AddConvertComponent(app.ComponentTypeLED)
			app.AddConvertComponent(app.ComponentTypeWallpaper)
			app.AddConvertComponent(app.ComponentTypeIcon)
			app.AddConvertComponent(app.ComponentTypeFont)
		case "Accents Only":
			app.SetConvertAllComponents(false)
			app.ClearSelectedConvertComponents()
			app.AddConvertComponent(app.ComponentTypeAccent)
		case "LEDs Only":
			app.SetConvertAllComponents(false)
			app.ClearSelectedConvertComponents()
			app.AddConvertComponent(app.ComponentTypeLED)
		case "Wallpapers Only":
			app.SetConvertAllComponents(false)
			app.ClearSelectedConvertComponents()
			app.AddConvertComponent(app.ComponentTypeWallpaper)
		case "Icons Only":
			app.SetConvertAllComponents(false)
			app.ClearSelectedConvertComponents()
			app.AddConvertComponent(app.ComponentTypeIcon)
		case "Fonts Only":
			app.SetConvertAllComponents(false)
			app.ClearSelectedConvertComponents()
			app.AddConvertComponent(app.ComponentTypeFont)
		}

		return app.Screens.ThemeConvertConfirm

	case 1, 2:
		// User pressed cancel or back
		return app.Screens.ThemeConvertOptions
	}

	return app.Screens.ThemeConvertComponentSelection
}

// ThemeConvertConfirmScreen displays final confirmation before theme conversion
func ThemeConvertConfirmScreen() (string, int) {
	themeName := app.GetSelectedConvertTheme()

	var message string
	if app.GetConvertAllComponents() {
		message = fmt.Sprintf("Extract all components from %s?", themeName)
	} else {
		// In a real implementation, we'd list the specific components
		message = fmt.Sprintf("Extract selected components from %s?", themeName)
	}

	menu := []string{
		"Yes",
		"No",
	}

	return ui.DisplayMinUiList(strings.Join(menu, "\n"), "text", message)
}

// HandleThemeConvertConfirm processes final confirmation for theme conversion
func HandleThemeConvertConfirm(selection string, exitCode int) app.Screen {
	logging.LogDebug("HandleThemeConvertConfirm called with selection: '%s', exitCode: %d", selection, exitCode)

	switch exitCode {
	case 0:
		if selection == "Yes" {
			// Perform the actual conversion
			err := performThemeConversion()
			if err != nil {
				logging.LogDebug("Conversion error: %v", err)
				ui.ShowMessage(fmt.Sprintf("Conversion error: %s", err), "3")
			} else {
				ui.ShowMessage("Theme components extracted successfully", "3")
			}
		}
		return app.Screens.ThemeManagementMenu

	case 1, 2:
		// User pressed cancel or back
		return app.Screens.ThemeConvertComponentSelection
	}

	return app.Screens.ThemeConvertConfirm
}

// performThemeConversion executes the actual theme conversion based on settings
func performThemeConversion() error {
	// Get conversion parameters from app state
	themeName := app.GetSelectedConvertTheme()
	convertAllComponents := app.GetConvertAllComponents()

	// Convert selected components
	var selectedComponents []themes.ComponentType
	if !convertAllComponents {
		// Get selected components from app state
		appComponents := app.GetSelectedConvertComponents()
		for _, compType := range appComponents {
			selectedComponents = append(selectedComponents, themes.ComponentType(compType))
		}
	}

	logging.LogDebug("Conversion parameters: Theme=%s, AllComponents=%v, SelectedCount=%d",
		themeName, convertAllComponents, len(selectedComponents))

	// Execute the theme conversion operation in the themes package
	return themes.PerformThemeConversion(themeName, convertAllComponents, selectedComponents)
}
