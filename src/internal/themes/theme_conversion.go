// src/internal/themes/theme_conversion.go
// Implementation of theme conversion and deconstruction functionality

package themes

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

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
		}
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
		}
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

	// Copy fonts from theme to font pack
	for fontType, mapping := range manifest.PathMappings.Fonts {
		srcPath := filepath.Join(themePath, mapping.ThemePath)

		// Determine destination filename based on font type
		var dstFileName string
		switch fontType {
		case "og_font":
			dstFileName = "OG.ttf"
		case "og_backup":
			dstFileName = "OG.backup.ttf"
		case "next_font":
			dstFileName = "Next.ttf"
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
		}
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

	// Check if the theme has accent settings
	srcPath := filepath.Join(themePath, "Settings", "minuisettings.txt")
	dstPath := filepath.Join(exportPath, "minuisettings.txt")

	// Check if source file exists
	if _, err := os.Stat(srcPath); os.IsNotExist(err) {
		logger.Printf("Warning: Source accent settings file not found: %s", srcPath)
		return fmt.Errorf("accent settings file not found: %s", srcPath)
	}

	// Copy the file
	if err := CopyFile(srcPath, dstPath); err != nil {
		logger.Printf("Warning: Could not copy accent settings: %v", err)
		return fmt.Errorf("failed to copy accent settings: %w", err)
	}

	logger.Printf("Copied accent settings: %s -> %s", srcPath, dstPath)
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

	// Check if the theme has LED settings
	srcPath := filepath.Join(themePath, "Settings", "ledsettings_brick.txt")
	dstPath := filepath.Join(exportPath, "ledsettings_brick.txt")

	// Check if source file exists
	if _, err := os.Stat(srcPath); os.IsNotExist(err) {
		logger.Printf("Warning: Source LED settings file not found: %s", srcPath)
		return fmt.Errorf("LED settings file not found: %s", srcPath)
	}

	// Copy the file
	if err := CopyFile(srcPath, dstPath); err != nil {
		logger.Printf("Warning: Could not copy LED settings: %v", err)
		return fmt.Errorf("failed to copy LED settings: %w", err)
	}

	logger.Printf("Copied LED settings: %s -> %s", srcPath, dstPath)
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

	// Determine directory and prefix based on component type
	var exportDir string
	var prefix string
	var fileExt string

	switch componentType {
	case ComponentTypeWallpaper:
		exportDir = filepath.Join(cwd, "Wallpapers", "Exports")
		prefix = "wallpaper"
		fileExt = ".bg"
	case ComponentTypeIcon:
		exportDir = filepath.Join(cwd, "Icons", "Exports")
		prefix = "icon"
		fileExt = ".icon"
	case ComponentTypeFont:
		exportDir = filepath.Join(cwd, "Fonts", "Exports")
		prefix = "font"
		fileExt = ".font"
	case ComponentTypeAccent:
		exportDir = filepath.Join(cwd, "Accents", "Exports")
		prefix = "accent"
		fileExt = ".acc"
	case ComponentTypeLED:
		exportDir = filepath.Join(cwd, "LEDs", "Exports")
		prefix = "led"
		fileExt = ".led"
	default:
		exportDir = filepath.Join(cwd, "Exports")
		prefix = "component"
		fileExt = ""
	}

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
