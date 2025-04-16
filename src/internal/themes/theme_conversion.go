// src/internal/themes/theme_conversion.go
// Implementation of theme conversion and deconstruction functionality

package themes

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"nextui-themes/internal/accents"
	"nextui-themes/internal/ui"
)

// ConvertTheme deconstructs a theme package into individual component packages
func ConvertTheme(themeName string, convertAllComponents bool, selectedComponents []ComponentType) error {
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
		filepath.Join(logsDir, "theme_conversion.log"),
		os.O_CREATE|os.O_APPEND|os.O_WRONLY,
		0644,
	)
	if err != nil {
		return fmt.Errorf("error creating log file: %w", err)
	}
	defer logFile.Close()

	// Create logger
	logger := &Logger{logFile}
	logger.Printf("Starting theme conversion for: %s", themeName)

	// Full path to theme - look in Imports directory
	themePath := filepath.Join(cwd, "Themes", "Imports", themeName)

	// Validate theme
	manifest, err := ValidateTheme(themePath, logger)
	if err != nil {
		logger.Printf("Theme validation failed: %v", err)
		return fmt.Errorf("theme validation failed: %w", err)
	}

	// Get theme name without extension
	themeName = strings.TrimSuffix(themeName, ".theme")

	// Check which components to convert
	var componentsToConvert []ComponentType
	if convertAllComponents {
		componentsToConvert = []ComponentType{
			ComponentTypeWallpaper,
			ComponentTypeIcon,
			ComponentTypeFont,
			ComponentTypeAccent,
			ComponentTypeLED,
		}
	} else {
		componentsToConvert = selectedComponents
	}

	// Process each component type
	for _, componentType := range componentsToConvert {
		switch componentType {
		case ComponentTypeWallpaper:
			if manifest.Content.Wallpapers.Present {
				if err := ConvertWallpapers(themePath, manifest, themeName, logger); err != nil {
					logger.Printf("Error converting wallpapers: %v", err)
					// Continue with other components
				}
			}

		case ComponentTypeIcon:
			if manifest.Content.Icons.Present {
				if err := ConvertIcons(themePath, manifest, themeName, logger); err != nil {
					logger.Printf("Error converting icons: %v", err)
					// Continue with other components
				}
			}

		case ComponentTypeFont:
			if manifest.Content.Fonts.Present {
				if err := ConvertFonts(themePath, manifest, themeName, logger); err != nil {
					logger.Printf("Error converting fonts: %v", err)
					// Continue with other components
				}
			}

		case ComponentTypeAccent:
			if manifest.Content.Settings.AccentsIncluded {
				if err := ConvertAccents(themePath, manifest, themeName, logger); err != nil {
					logger.Printf("Error converting accents: %v", err)
					// Continue with other components
				}
			}

		case ComponentTypeLED:
			if manifest.Content.Settings.LEDsIncluded {
				if err := ConvertLEDs(themePath, manifest, themeName, logger); err != nil {
					logger.Printf("Error converting LEDs: %v", err)
					// Continue with other components
				}
			}
		}
	}

	logger.Printf("Theme conversion completed successfully")
	ui.ShowMessage("Theme has been deconstructed into components successfully", "3")

	return nil
}

// ConvertWallpapers extracts wallpapers from a theme and saves them as a wallpaper pack
func ConvertWallpapers(themePath string, manifest *ThemeManifest, baseName string, logger *Logger) error {
	logger.Printf("Converting wallpapers from theme: %s", themePath)

	// Create a unique name for the wallpaper pack
	packName := generateUniqueComponentName(baseName, ComponentTypeWallpaper)

	// Create export directory
	exportPath, err := CreateComponentExportDirectory(packName, ComponentTypeWallpaper)
	if err != nil {
		return fmt.Errorf("error creating export directory: %w", err)
	}

	// Create a new manifest for the wallpaper pack
	wallpaperManifest := &ThemeManifest{}
	wallpaperManifest.ComponentType = "Wallpapers"
	wallpaperManifest.Content.Wallpapers.Present = true
	wallpaperManifest.PathMappings.Wallpapers = []PathMapping{}

	// Copy wallpapers from theme to wallpaper pack
	for _, mapping := range manifest.PathMappings.Wallpapers {
		srcPath := filepath.Join(themePath, mapping.ThemePath)

		// Determine destination path based on metadata
		var dstDir string
		if mapping.Metadata != nil {
			if wallpaperType, ok := mapping.Metadata["WallpaperType"]; ok {
				switch wallpaperType {
				case "Collection":
					dstDir = "CollectionWallpapers"
				default:
					dstDir = "SystemWallpapers"
				}
			} else {
				dstDir = "SystemWallpapers"
			}
		} else {
			dstDir = "SystemWallpapers"
		}

		// Extract filename from the source path
		fileName := filepath.Base(srcPath)
		dstPath := filepath.Join(exportPath, dstDir, fileName)
		relativePath := filepath.Join(dstDir, fileName)

		// Check if source file exists
		if _, err := os.Stat(srcPath); os.IsNotExist(err) {
			logger.Printf("Warning: Source wallpaper file not found: %s", srcPath)
			continue
		}

		// Create destination directory if it doesn't exist
		dstDirPath := filepath.Dir(dstPath)
		if err := os.MkdirAll(dstDirPath, 0755); err != nil {
			logger.Printf("Warning: Could not create directory: %v", err)
			continue
		}

		// Copy the file
		if err := CopyFile(srcPath, dstPath); err != nil {
			logger.Printf("Warning: Could not copy wallpaper: %v", err)
		} else {
			logger.Printf("Copied wallpaper: %s -> %s", srcPath, dstPath)

			// Add to the manifest
			newMapping := PathMapping{
				ThemePath:  relativePath,
				SystemPath: mapping.SystemPath,
				Metadata:   mapping.Metadata,
			}
			wallpaperManifest.PathMappings.Wallpapers = append(wallpaperManifest.PathMappings.Wallpapers, newMapping)
		}
	}

	// Create preview by copying a suitable wallpaper
	if err := createWallpaperPreviewForConversion(exportPath, wallpaperManifest); err != nil {
		logger.Printf("Warning: Could not create preview image: %v", err)
	}

	// Write the manifest
	if err := WriteComponentManifest(exportPath, wallpaperManifest); err != nil {
		logger.Printf("Error writing manifest: %v", err)
		return fmt.Errorf("error writing manifest: %w", err)
	}

	logger.Printf("Wallpapers converted successfully: %s", exportPath)
	return nil
}

// ConvertIcons extracts icons from a theme and saves them as an icon pack
func ConvertIcons(themePath string, manifest *ThemeManifest, baseName string, logger *Logger) error {
	logger.Printf("Converting icons from theme: %s", themePath)

	// Create a unique name for the icon pack
	packName := generateUniqueComponentName(baseName, ComponentTypeIcon)

	// Create export directory
	exportPath, err := CreateComponentExportDirectory(packName, ComponentTypeIcon)
	if err != nil {
		return fmt.Errorf("error creating export directory: %w", err)
	}

	// Create a new manifest for the icon pack
	iconManifest := &ThemeManifest{}
	iconManifest.ComponentType = "Icons"
	iconManifest.Content.Icons.Present = true
	iconManifest.PathMappings.Icons = []PathMapping{}

	// Copy icons from theme to icon pack
	for _, mapping := range manifest.PathMappings.Icons {
		srcPath := filepath.Join(themePath, mapping.ThemePath)

		// Determine destination path based on metadata
		var dstDir string
		if mapping.Metadata != nil {
			if iconType, ok := mapping.Metadata["IconType"]; ok {
				switch iconType {
				case "Tool":
					dstDir = "ToolIcons"
				case "Collection":
					dstDir = "CollectionIcons"
				default:
					dstDir = "SystemIcons"
				}
			} else {
				dstDir = "SystemIcons"
			}
		} else {
			dstDir = "SystemIcons"
		}

		// Extract filename from the source path
		fileName := filepath.Base(srcPath)
		dstPath := filepath.Join(exportPath, dstDir, fileName)
		relativePath := filepath.Join(dstDir, fileName)

		// Check if source file exists
		if _, err := os.Stat(srcPath); os.IsNotExist(err) {
			logger.Printf("Warning: Source icon file not found: %s", srcPath)
			continue
		}

		// Create destination directory if it doesn't exist
		dstDirPath := filepath.Dir(dstPath)
		if err := os.MkdirAll(dstDirPath, 0755); err != nil {
			logger.Printf("Warning: Could not create directory: %v", err)
			continue
		}

		// Copy the file
		if err := CopyFile(srcPath, dstPath); err != nil {
			logger.Printf("Warning: Could not copy icon: %v", err)
		} else {
			logger.Printf("Copied icon: %s -> %s", srcPath, dstPath)

			// Add to the manifest
			newMapping := PathMapping{
				ThemePath:  relativePath,
				SystemPath: mapping.SystemPath,
				Metadata:   mapping.Metadata,
			}
			iconManifest.PathMappings.Icons = append(iconManifest.PathMappings.Icons, newMapping)
		}
	}

	// Write the manifest
	if err := WriteComponentManifest(exportPath, iconManifest); err != nil {
		logger.Printf("Error writing manifest: %v", err)
		return fmt.Errorf("error writing manifest: %w", err)
	}

	logger.Printf("Icons converted successfully: %s", exportPath)
	return nil
}

// ConvertFonts extracts fonts from a theme and saves them as a font pack
func ConvertFonts(themePath string, manifest *ThemeManifest, baseName string, logger *Logger) error {
	logger.Printf("Converting fonts from theme: %s", themePath)

	// Create a unique name for the font pack
	packName := generateUniqueComponentName(baseName, ComponentTypeFont)

	// Create export directory
	exportPath, err := CreateComponentExportDirectory(packName, ComponentTypeFont)
	if err != nil {
		return fmt.Errorf("error creating export directory: %w", err)
	}

	// Create a new manifest for the font pack
	fontManifest := &ThemeManifest{}
	fontManifest.ComponentType = "Fonts"
	fontManifest.Content.Fonts.Present = true
	fontManifest.PathMappings.Fonts = make(map[string]PathMapping)

	// Copy fonts from theme to font pack
	for fontType, mapping := range manifest.PathMappings.Fonts {
		srcPath := filepath.Join(themePath, mapping.ThemePath)

		// Determine destination filename based on font type
		var dstFileName string
		switch fontType {
		case "og_font":
			dstFileName = "OG.ttf"
			fontManifest.Content.Fonts.OGReplaced = true
		case "og_backup":
			dstFileName = "OG.backup.ttf"
		case "next_font":
			dstFileName = "Next.ttf"
			fontManifest.Content.Fonts.NextReplaced = true
		case "next_backup":
			dstFileName = "Next.backup.ttf"
		default:
			// Use the original filename as fallback
			dstFileName = filepath.Base(srcPath)
		}

		dstPath := filepath.Join(exportPath, dstFileName)

		// Check if source file exists
		if _, err := os.Stat(srcPath); os.IsNotExist(err) {
			logger.Printf("Warning: Source font file not found: %s", srcPath)
			continue
		}

		// Copy the file
		if err := CopyFile(srcPath, dstPath); err != nil {
			logger.Printf("Warning: Could not copy font: %v", err)
		} else {
			logger.Printf("Copied font: %s -> %s", srcPath, dstPath)

			// Add to the manifest
			fontManifest.PathMappings.Fonts[fontType] = PathMapping{
				ThemePath:  dstFileName,
				SystemPath: mapping.SystemPath,
			}
		}
	}

	// Write the manifest
	if err := WriteComponentManifest(exportPath, fontManifest); err != nil {
		logger.Printf("Error writing manifest: %v", err)
		return fmt.Errorf("error writing manifest: %w", err)
	}

	logger.Printf("Fonts converted successfully: %s", exportPath)
	return nil
}

// ConvertAccents extracts accent settings from a theme and saves them as an accent pack
func ConvertAccents(themePath string, manifest *ThemeManifest, baseName string, logger *Logger) error {
	logger.Printf("Converting accent settings from theme: %s", themePath)

	// Create a unique name for the accent pack
	packName := generateUniqueComponentName(baseName, ComponentTypeAccent)

	// Create export directory
	exportPath, err := CreateComponentExportDirectory(packName, ComponentTypeAccent)
	if err != nil {
		return fmt.Errorf("error creating export directory: %w", err)
	}

	// Create a new manifest for the accent pack
	accentManifest := &ThemeManifest{}
	accentManifest.ComponentType = "Accents"
	accentManifest.Content.Settings.AccentsIncluded = true

	// Check if we already have accent colors in the manifest
	if manifest.AccentColors != nil && len(manifest.AccentColors) > 0 {
		// Copy accent colors from the theme manifest
		accentManifest.AccentColors = make(map[string]string)
		for key, value := range manifest.AccentColors {
			accentManifest.AccentColors[key] = value
		}
		logger.Printf("Copied accent colors from theme manifest")
	} else {
		// Legacy approach - extract from file
		srcPath := filepath.Join(themePath, "Settings", "minuisettings.txt")

		// Check if source file exists
		if _, err := os.Stat(srcPath); os.IsNotExist(err) {
			logger.Printf("Warning: Source accent settings file not found: %s", srcPath)
			return fmt.Errorf("accent settings file not found: %s", srcPath)
		}

		// Parse accent colors from file
		accentColors, err := accents.ParseAccentColors(srcPath)
		if err != nil {
			logger.Printf("Warning: Could not parse accent colors: %v", err)
			return fmt.Errorf("failed to parse accent colors: %w", err)
		}

		accentManifest.AccentColors = accentColors
		logger.Printf("Extracted accent colors from settings file")
	}

	// Write the manifest
	if err := WriteComponentManifest(exportPath, accentManifest); err != nil {
		logger.Printf("Error writing manifest: %v", err)
		return fmt.Errorf("error writing manifest: %w", err)
	}

	logger.Printf("Accent settings converted successfully: %s", exportPath)
	return nil
}

// ConvertLEDs extracts LED settings from a theme and saves them as a LED pack
func ConvertLEDs(themePath string, manifest *ThemeManifest, baseName string, logger *Logger) error {
	logger.Printf("Converting LED settings from theme: %s", themePath)

	// Create a unique name for the LED pack
	packName := generateUniqueComponentName(baseName, ComponentTypeLED)

	// Create export directory
	exportPath, err := CreateComponentExportDirectory(packName, ComponentTypeLED)
	if err != nil {
		return fmt.Errorf("error creating export directory: %w", err)
	}

	// Create a new manifest for the LED pack
	ledManifest := &ThemeManifest{}
	ledManifest.ComponentType = "LEDs"
	ledManifest.Content.Settings.LEDsIncluded = true

	// Check if we already have LED settings in the manifest
	if manifest.LEDSettings != nil && len(manifest.LEDSettings) > 0 {
		// Copy LED settings from the theme manifest
		ledManifest.LEDSettings = make(map[string]map[string]interface{})
		for section, settings := range manifest.LEDSettings {
			ledManifest.LEDSettings[section] = make(map[string]interface{})
			for key, value := range settings {
				ledManifest.LEDSettings[section][key] = value
			}
		}
		logger.Printf("Copied LED settings from theme manifest")
	} else {
		// Legacy approach - extract from file
		srcPath := filepath.Join(themePath, "Settings", "ledsettings_brick.txt")

		// Check if source file exists
		if _, err := os.Stat(srcPath); os.IsNotExist(err) {
			logger.Printf("Warning: Source LED settings file not found: %s", srcPath)
			return fmt.Errorf("LED settings file not found: %s", srcPath)
		}

		// Read the settings file
		data, err := os.ReadFile(srcPath)
		if err != nil {
			logger.Printf("Warning: Could not read LED settings file: %v", err)
			return fmt.Errorf("failed to read LED settings file: %w", err)
		}

		// Parse LED settings
		ledSettings := parseLEDSettingsForConversion(string(data))
		ledManifest.LEDSettings = ledSettings
		logger.Printf("Extracted LED settings from settings file")
	}

	// Write the manifest
	if err := WriteComponentManifest(exportPath, ledManifest); err != nil {
		logger.Printf("Error writing manifest: %v", err)
		return fmt.Errorf("error writing manifest: %w", err)
	}

	logger.Printf("LED settings converted successfully: %s", exportPath)
	return nil
}

// generateUniqueComponentName creates a unique name for a component based on the theme name
func generateUniqueComponentName(baseName string, componentType ComponentType) string {
	// Get the current directory
	cwd, err := os.Getwd()
	if err != nil {
		return "default_1"
	}

	// Determine prefix and file extension based on component type
	var prefix string
	var fileExt string

	switch componentType {
	case ComponentTypeWallpaper:
		prefix = "wallpaper"
		fileExt = ".bg"
	case ComponentTypeIcon:
		prefix = "icon"
		fileExt = ".icon"
	case ComponentTypeFont:
		prefix = "font"
		fileExt = ".font"
	case ComponentTypeAccent:
		prefix = "accent"
		fileExt = ".acc"
	case ComponentTypeLED:
		prefix = "led"
		fileExt = ".led"
	default:
		prefix = "component"
		fileExt = ""
	}

	// Use a single exports directory
	exportDir := filepath.Join(cwd, "Theme-Manager.pak", "Exports")

	// Generate a base name using the theme name
	baseName = strings.ReplaceAll(baseName, " ", "_")
	proposedName := fmt.Sprintf("%s_from_%s", prefix, baseName)

	// Check if the name already exists
	_, err = os.Stat(filepath.Join(exportDir, proposedName+fileExt))
	if os.IsNotExist(err) {
		return proposedName
	}

	// Find a unique name
	counter := 1
	for {
		uniqueName := fmt.Sprintf("%s_%d", proposedName, counter)
		_, err = os.Stat(filepath.Join(exportDir, uniqueName+fileExt))
		if os.IsNotExist(err) {
			return uniqueName
		}
		counter++
	}
}

// Connect the UI's performThemeConversion function to our backend implementation
func PerformThemeConversion(themeName string, convertAllComponents bool, selectedComponents []ComponentType) error {
	return ConvertTheme(themeName, convertAllComponents, selectedComponents)
}

// parseLEDSettingsForConversion parses LED settings from a string
// This is a duplicate of the function in manifest_updater.go with a different name to avoid conflicts
func parseLEDSettingsForConversion(data string) map[string]map[string]interface{} {
	// Simple implementation to parse LED settings
	result := make(map[string]map[string]interface{})

	// Example parsing - would be replaced by actual parsing logic
	sections := strings.Split(data, "[")
	for _, section := range sections {
		if len(section) == 0 {
			continue
		}

		lines := strings.Split(section, "\n")
		if len(lines) < 2 {
			continue
		}

		// Extract section name
		sectionName := strings.TrimRight(lines[0], "]")
		result[sectionName] = make(map[string]interface{})

		// Extract key-value pairs
		for i := 1; i < len(lines); i++ {
			line := strings.TrimSpace(lines[i])
			if line == "" {
				continue
			}

			parts := strings.SplitN(line, "=", 2)
			if len(parts) != 2 {
				continue
			}

			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])

			// Try to convert numeric values
			if strings.HasPrefix(value, "#") {
				// Color value
				result[sectionName][key] = value
			} else if intVal, err := strconv.Atoi(value); err == nil {
				result[sectionName][key] = intVal
			} else {
				result[sectionName][key] = value
			}
		}
	}

	return result
}

// createWallpaperPreviewForConversion creates a preview image for a wallpaper pack during conversion
func createWallpaperPreviewForConversion(themePath string, manifest *ThemeManifest) error {
	previewPath := filepath.Join(themePath, "preview.png")

	// If preview already exists, update manifest and skip
	if _, err := os.Stat(previewPath); err == nil {
		manifest.PreviewImage = "preview.png"
		return nil
	}

	// Look for Recently Played wallpaper
	var sourcePath string
	for _, mapping := range manifest.PathMappings.Wallpapers {
		if mapping.Metadata != nil {
			if name, ok := mapping.Metadata["SystemName"]; ok && name == "Recently Played" {
				sourcePath = filepath.Join(themePath, mapping.ThemePath)
				break
			}
		}
	}

	// If not found, try to find Root wallpaper
	if sourcePath == "" {
		for _, mapping := range manifest.PathMappings.Wallpapers {
			if mapping.Metadata != nil {
				if name, ok := mapping.Metadata["SystemName"]; ok && name == "Root" {
					sourcePath = filepath.Join(themePath, mapping.ThemePath)
					break
				}
			}
		}
	}

	// If still not found, use the first wallpaper
	if sourcePath == "" && len(manifest.PathMappings.Wallpapers) > 0 {
		sourcePath = filepath.Join(themePath, manifest.PathMappings.Wallpapers[0].ThemePath)
	}

	// If we have a source path, copy it as preview
	if sourcePath != "" {
		if err := CopyFile(sourcePath, previewPath); err != nil {
			return err
		}
		manifest.PreviewImage = "preview.png"
	}

	return nil
}
