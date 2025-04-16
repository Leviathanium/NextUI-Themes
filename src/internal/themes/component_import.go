// src/internal/themes/component_import.go
// Implementation of component-specific import functionality

package themes

import (
	"fmt"
	"os"
	"path/filepath"

	"nextui-themes/internal/accents"
	"nextui-themes/internal/ui"
)

// ComponentType represents different component types that can be imported
type ComponentType int

const (
	ComponentTypeFullTheme ComponentType = iota + 1
	ComponentTypeAccent
	ComponentTypeLED
	ComponentTypeWallpaper
	ComponentTypeIcon
	ComponentTypeFont
)

// ImportComponent imports a specific component type
func ImportComponent(componentType ComponentType, itemName string, selectedComponents map[ComponentType]bool) error {
	// Create logging directory if it doesn't exist
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting current directory: %w", err)
	}

	logsDir := filepath.Join(cwd, "Logs")
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		return fmt.Errorf("error creating logs directory: %w", err)
	}

	// Create log file
	logFile, err := os.OpenFile(
		filepath.Join(logsDir, "component_imports.log"),
		os.O_CREATE|os.O_APPEND|os.O_WRONLY,
		0644,
	)
	if err != nil {
		return fmt.Errorf("error creating log file: %w", err)
	}
	defer logFile.Close()

	// Create logger
	logger := &Logger{logFile}
	logger.Printf("Starting component import: Type=%d, Item=%s", componentType, itemName)

	// Handle different component types
	switch componentType {
	case ComponentTypeFullTheme:
		// For full themes, check if we're importing all components or specific ones
		if len(selectedComponents) == 0 {
			// Import all components (legacy mode)
			return ImportTheme(itemName)
		}

		// Import specific components from the theme
		return ImportThemeComponents(itemName, selectedComponents, logger)

	case ComponentTypeAccent:
		return ImportAccentPack(itemName, logger)

	case ComponentTypeLED:
		return ImportLEDPack(itemName, logger)

	case ComponentTypeWallpaper:
		return ImportWallpaperPack(itemName, logger)

	case ComponentTypeIcon:
		return ImportIconPack(itemName, logger)

	case ComponentTypeFont:
		return ImportFontPack(itemName, logger)

	default:
		return fmt.Errorf("unsupported component type: %d", componentType)
	}
}

// ImportThemeComponents imports specific components from a theme
func ImportThemeComponents(themeName string, selectedComponents map[ComponentType]bool, logger *Logger) error {
	logger.Printf("Importing specific components from theme: %s", themeName)

	// Full path to theme - look in Themes/Imports directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting current directory: %w", err)
	}

	themePath := filepath.Join(cwd, "Themes", "Imports", themeName)

	// Validate theme
	manifest, err := ValidateTheme(themePath, logger)
	if err != nil {
		logger.Printf("Theme validation failed: %v", err)
		return fmt.Errorf("theme validation failed: %w", err)
	}

	// Import selected components
	if selectedComponents[ComponentTypeWallpaper] {
		logger.Printf("Importing wallpapers from theme")
		if err := ImportWallpapers(themePath, manifest, logger); err != nil {
			logger.Printf("Error importing wallpapers: %v", err)
			// Continue with other components
		}
	}

	if selectedComponents[ComponentTypeIcon] {
		logger.Printf("Importing icons from theme")
		if err := ImportIcons(themePath, manifest, logger); err != nil {
			logger.Printf("Error importing icons: %v", err)
			// Continue with other components
		}
	}

	if selectedComponents[ComponentTypeFont] {
		logger.Printf("Importing fonts from theme")
		if err := ImportFonts(themePath, manifest, logger); err != nil {
			logger.Printf("Error importing fonts: %v", err)
			// Continue with other components
		}
	}

	if selectedComponents[ComponentTypeAccent] {
		logger.Printf("Importing accent settings from theme")
		// Extract just accent settings from the theme
		if manifest.Content.Settings.AccentsIncluded {
			if err := ImportSettings(themePath, manifest, logger); err != nil {
				logger.Printf("Error importing settings: %v", err)
				// Continue with other components
			}
		} else {
			logger.Printf("No accent settings found in theme")
		}
	}

	if selectedComponents[ComponentTypeLED] {
		logger.Printf("Importing LED settings from theme")
		// Extract just LED settings from the theme
		if manifest.Content.Settings.LEDsIncluded {
			if err := ImportSettings(themePath, manifest, logger); err != nil {
				logger.Printf("Error importing settings: %v", err)
				// Continue with other components
			}
		} else {
			logger.Printf("No LED settings found in theme")
		}
	}

	logger.Printf("Component import completed")
	return nil
}

// ImportAccentPack imports an accent pack (.acc)
func ImportAccentPack(packName string, logger *Logger) error {
	logger.Printf("Importing accent pack: %s", packName)

	// Get current directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting current directory: %w", err)
	}

	// Path to accent pack
	packPath := filepath.Join(cwd, "Accents", "Imports", packName)

	// Ensure directory exists
	_, err = os.Stat(packPath)
	if os.IsNotExist(err) {
		return fmt.Errorf("accent pack not found: %s", packPath)
	}

	// Check for manifest.json
	manifestPath := filepath.Join(packPath, "manifest.json")
	var manifest *ThemeManifest

	if _, err := os.Stat(manifestPath); os.IsNotExist(err) {
		// No manifest, use legacy approach to find settings file
		settingsPath := filepath.Join(packPath, "minuisettings.txt")
		if _, err := os.Stat(settingsPath); os.IsNotExist(err) {
			return fmt.Errorf("settings file not found in accent pack: %s", settingsPath)
		}

		// Copy settings file to system location
		systemSettingsPath := accents.SettingsPath
		if err := copyComponentFile(settingsPath, systemSettingsPath); err != nil {
			return fmt.Errorf("failed to copy accent settings: %w", err)
		}
	} else {
		// Load manifest
		manifest, err = ValidateTheme(packPath, logger)
		if err != nil {
			logger.Printf("Error validating manifest: %v", err)

			// Try to update the manifest if it's invalid
			logger.Printf("Attempting to update manifest...")
			if err := EnhanceManifestUpdater(packPath); err != nil {
				logger.Printf("Error updating manifest: %v", err)
				return fmt.Errorf("error updating manifest: %w", err)
			}

			// Load updated manifest
			manifest, err = ValidateTheme(packPath, logger)
			if err != nil {
				logger.Printf("Error validating updated manifest: %v", err)
				return fmt.Errorf("error validating updated manifest: %w", err)
			}
		}

		// Apply accent colors from manifest
		if manifest.Content.Settings.AccentsIncluded && manifest.AccentColors != nil {
			logger.Printf("Applying accent colors from manifest")

			// Update the system settings file with the accent colors
			if err := accents.ApplyAccentColors(manifest.AccentColors); err != nil {
				logger.Printf("Error applying accent colors: %v", err)
				return fmt.Errorf("error applying accent colors: %w", err)
			}
		} else {
			logger.Printf("No accent colors found in manifest")
			return fmt.Errorf("no accent colors found in manifest")
		}
	}

	logger.Printf("Accent pack imported successfully")
	ui.ShowMessage(fmt.Sprintf("Accent pack '%s' imported successfully", packName), "3")
	return nil
}

// ImportLEDPack imports a LED pack (.led)
func ImportLEDPack(packName string, logger *Logger) error {
	logger.Printf("Importing LED pack: %s", packName)

	// Get current directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting current directory: %w", err)
	}

	// Path to LED pack
	packPath := filepath.Join(cwd, "LEDs", "Imports", packName)

	// Ensure directory exists
	_, err = os.Stat(packPath)
	if os.IsNotExist(err) {
		return fmt.Errorf("LED pack not found: %s", packPath)
	}

	// Check for manifest.json
	manifestPath := filepath.Join(packPath, "manifest.json")
	var manifest *ThemeManifest

	if _, err := os.Stat(manifestPath); os.IsNotExist(err) {
		// No manifest, use legacy approach to find settings file
		settingsPath := filepath.Join(packPath, "ledsettings_brick.txt")
		if _, err := os.Stat(settingsPath); os.IsNotExist(err) {
			return fmt.Errorf("settings file not found in LED pack: %s", settingsPath)
		}

		// Copy settings file to system location
		systemSettingsPath := "/mnt/SDCARD/.userdata/shared/ledsettings_brick.txt"
		if err := copyComponentFile(settingsPath, systemSettingsPath); err != nil {
			return fmt.Errorf("failed to copy LED settings: %w", err)
		}
	} else {
		// Load manifest
		manifest, err = ValidateTheme(packPath, logger)
		if err != nil {
			logger.Printf("Error validating manifest: %v", err)

			// Try to update the manifest if it's invalid
			logger.Printf("Attempting to update manifest...")
			if err := EnhanceManifestUpdater(packPath); err != nil {
				logger.Printf("Error updating manifest: %v", err)
				return fmt.Errorf("error updating manifest: %w", err)
			}

			// Load updated manifest
			manifest, err = ValidateTheme(packPath, logger)
			if err != nil {
				logger.Printf("Error validating updated manifest: %v", err)
				return fmt.Errorf("error validating updated manifest: %w", err)
			}
		}

		// Apply LED settings from manifest
		if manifest.Content.Settings.LEDsIncluded && manifest.LEDSettings != nil {
			logger.Printf("Applying LED settings from manifest")

			// Generate settings file from manifest data
			if err := ApplyLEDSettings(manifest.LEDSettings); err != nil {
				logger.Printf("Error applying LED settings: %v", err)
				return fmt.Errorf("error applying LED settings: %w", err)
			}
		} else {
			logger.Printf("No LED settings found in manifest")
			return fmt.Errorf("no LED settings found in manifest")
		}
	}

	logger.Printf("LED pack imported successfully")
	ui.ShowMessage(fmt.Sprintf("LED pack '%s' imported successfully", packName), "3")
	return nil
}

// ImportWallpaperPack imports a wallpaper pack (.bg)
func ImportWallpaperPack(packName string, logger *Logger) error {
	logger.Printf("Importing wallpaper pack: %s", packName)

	// Get current directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting current directory: %w", err)
	}

	// Path to wallpaper pack
	packPath := filepath.Join(cwd, "Wallpapers", "Imports", packName)

	// Ensure directory exists
	_, err = os.Stat(packPath)
	if os.IsNotExist(err) {
		return fmt.Errorf("wallpaper pack not found: %s", packPath)
	}

	// Check for manifest.json
	manifestPath := filepath.Join(packPath, "manifest.json")
	var manifest *ThemeManifest

	if _, err := os.Stat(manifestPath); os.IsNotExist(err) {
		// No manifest, create one
		logger.Printf("No manifest found in wallpaper pack, creating one...")

		// Initialize empty manifest
		manifest = &ThemeManifest{}
		manifest.ComponentType = "Wallpapers"

		// Update the manifest based on directory contents
		if err := EnhanceManifestUpdater(packPath); err != nil {
			logger.Printf("Error creating manifest: %v", err)
			return fmt.Errorf("error creating manifest: %w", err)
		}

		// Load the newly created manifest
		manifest, err = ValidateTheme(packPath, logger)
		if err != nil {
			logger.Printf("Error validating new manifest: %v", err)
			return fmt.Errorf("error validating new manifest: %w", err)
		}
	} else {
		// Load manifest
		manifest, err = ValidateTheme(packPath, logger)
		if err != nil {
			logger.Printf("Error validating manifest: %v", err)

			// Try to update the manifest if it's invalid
			logger.Printf("Attempting to update manifest...")
			if err := EnhanceManifestUpdater(packPath); err != nil {
				logger.Printf("Error updating manifest: %v", err)
				return fmt.Errorf("error updating manifest: %w", err)
			}

			// Load updated manifest
			manifest, err = ValidateTheme(packPath, logger)
			if err != nil {
				logger.Printf("Error validating updated manifest: %v", err)
				return fmt.Errorf("error validating updated manifest: %w", err)
			}
		}
	}

	// Import wallpapers from manifest
	if err := ImportWallpapers(packPath, manifest, logger); err != nil {
		logger.Printf("Error importing wallpapers: %v", err)
		return fmt.Errorf("error importing wallpapers: %w", err)
	}

	logger.Printf("Wallpaper pack imported successfully")
	ui.ShowMessage(fmt.Sprintf("Wallpaper pack '%s' imported successfully", packName), "3")
	return nil
}

// ImportIconPack imports an icon pack (.icon)
func ImportIconPack(packName string, logger *Logger) error {
	logger.Printf("Importing icon pack: %s", packName)

	// Get current directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting current directory: %w", err)
	}

	// Path to icon pack
	packPath := filepath.Join(cwd, "Icons", "Imports", packName)

	// Ensure directory exists
	_, err = os.Stat(packPath)
	if os.IsNotExist(err) {
		return fmt.Errorf("icon pack not found: %s", packPath)
	}

	// Check for manifest.json
	manifestPath := filepath.Join(packPath, "manifest.json")
	var manifest *ThemeManifest

	if _, err := os.Stat(manifestPath); os.IsNotExist(err) {
		// No manifest, create one
		logger.Printf("No manifest found in icon pack, creating one...")

		// Initialize empty manifest
		manifest = &ThemeManifest{}
		manifest.ComponentType = "Icons"

		// Update the manifest based on directory contents
		if err := EnhanceManifestUpdater(packPath); err != nil {
			logger.Printf("Error creating manifest: %v", err)
			return fmt.Errorf("error creating manifest: %w", err)
		}

		// Load the newly created manifest
		manifest, err = ValidateTheme(packPath, logger)
		if err != nil {
			logger.Printf("Error validating new manifest: %v", err)
			return fmt.Errorf("error validating new manifest: %w", err)
		}
	} else {
		// Load manifest
		manifest, err = ValidateTheme(packPath, logger)
		if err != nil {
			logger.Printf("Error validating manifest: %v", err)

			// Try to update the manifest if it's invalid
			logger.Printf("Attempting to update manifest...")
			if err := EnhanceManifestUpdater(packPath); err != nil {
				logger.Printf("Error updating manifest: %v", err)
				return fmt.Errorf("error updating manifest: %w", err)
			}

			// Load updated manifest
			manifest, err = ValidateTheme(packPath, logger)
			if err != nil {
				logger.Printf("Error validating updated manifest: %v", err)
				return fmt.Errorf("error validating updated manifest: %w", err)
			}
		}
	}

	// Import icons from manifest
	if err := ImportIcons(packPath, manifest, logger); err != nil {
		logger.Printf("Error importing icons: %v", err)
		return fmt.Errorf("error importing icons: %w", err)
	}

	logger.Printf("Icon pack imported successfully")
	ui.ShowMessage(fmt.Sprintf("Icon pack '%s' imported successfully", packName), "3")
	return nil
}

// ImportFontPack imports a font pack (.font)
func ImportFontPack(packName string, logger *Logger) error {
	logger.Printf("Importing font pack: %s", packName)

	// Get current directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting current directory: %w", err)
	}

	// Path to font pack
	packPath := filepath.Join(cwd, "Fonts", "Imports", packName)

	// Ensure directory exists
	_, err = os.Stat(packPath)
	if os.IsNotExist(err) {
		return fmt.Errorf("font pack not found: %s", packPath)
	}

	// Check for manifest.json
	manifestPath := filepath.Join(packPath, "manifest.json")
	var manifest *ThemeManifest

	if _, err := os.Stat(manifestPath); os.IsNotExist(err) {
		// No manifest, create one
		logger.Printf("No manifest found in font pack, creating one...")

		// Initialize empty manifest
		manifest = &ThemeManifest{}
		manifest.ComponentType = "Fonts"

		// Update the manifest based on directory contents
		if err := EnhanceManifestUpdater(packPath); err != nil {
			logger.Printf("Error creating manifest: %v", err)
			return fmt.Errorf("error creating manifest: %w", err)
		}

		// Load the newly created manifest
		manifest, err = ValidateTheme(packPath, logger)
		if err != nil {
			logger.Printf("Error validating new manifest: %v", err)
			return fmt.Errorf("error validating new manifest: %w", err)
		}
	} else {
		// Load manifest
		manifest, err = ValidateTheme(packPath, logger)
		if err != nil {
			logger.Printf("Error validating manifest: %v", err)

			// Try to update the manifest if it's invalid
			logger.Printf("Attempting to update manifest...")
			if err := EnhanceManifestUpdater(packPath); err != nil {
				logger.Printf("Error updating manifest: %v", err)
				return fmt.Errorf("error updating manifest: %w", err)
			}

			// Load updated manifest
			manifest, err = ValidateTheme(packPath, logger)
			if err != nil {
				logger.Printf("Error validating updated manifest: %v", err)
				return fmt.Errorf("error validating updated manifest: %w", err)
			}
		}
	}

	// Import fonts from manifest
	if err := ImportFonts(packPath, manifest, logger); err != nil {
		logger.Printf("Error importing fonts: %v", err)
		return fmt.Errorf("error importing fonts: %w", err)
	}

	logger.Printf("Font pack imported successfully")
	ui.ShowMessage(fmt.Sprintf("Font pack '%s' imported successfully", packName), "3")
	return nil
}

// ApplyLEDSettings applies LED settings from a map to the system
func ApplyLEDSettings(ledSettings map[string]map[string]interface{}) error {
	// Path to LED settings
	ledSettingsPath := "/mnt/SDCARD/.userdata/shared/ledsettings_brick.txt"

	// Convert LED settings to a string
	var settingsContent string

	for sectionName, settings := range ledSettings {
		settingsContent += "[" + sectionName + "]\n"

		for key, value := range settings {
			settingsContent += fmt.Sprintf("%s=%v\n", key, value)
		}

		settingsContent += "\n"
	}

	// Write settings to file
	return os.WriteFile(ledSettingsPath, []byte(settingsContent), 0644)
}

// PerformImport is the entry point for imports from the UI
func PerformImport(componentType ComponentType, itemName string, selectedComponents map[ComponentType]bool) error {
	return ImportComponent(componentType, itemName, selectedComponents)
}
