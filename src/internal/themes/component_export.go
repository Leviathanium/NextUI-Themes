// src/internal/themes/component_export.go
// Implementation of component-specific export functionality

package themes

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"nextui-themes/internal/accents"
	"nextui-themes/internal/ui"
)

// ExportComponent exports a specific component type
func ExportComponent(componentType ComponentType, exportName string, selectedComponents map[ComponentType]bool) error {
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
		filepath.Join(logsDir, "component_exports.log"),
		os.O_CREATE|os.O_APPEND|os.O_WRONLY,
		0644,
	)
	if err != nil {
		return fmt.Errorf("error creating log file: %w", err)
	}
	defer logFile.Close()

	// Create logger
	logger := &Logger{logFile}
	logger.Printf("Starting component export: Type=%d, Name=%s", componentType, exportName)

	// Handle different component types
	switch componentType {
	case ComponentTypeFullTheme:
		// For full themes, check if we're exporting all components or specific ones
		if len(selectedComponents) == 0 {
			// Export all components (legacy mode)
			return ExportTheme()
		}

		// Export specific components
		return ExportThemeComponents(exportName, selectedComponents, logger)

	case ComponentTypeAccent:
		return ExportAccentPack(exportName, logger)

	case ComponentTypeLED:
		return ExportLEDPack(exportName, logger)

	case ComponentTypeWallpaper:
		return ExportWallpaperPack(exportName, logger)

	case ComponentTypeIcon:
		return ExportIconPack(exportName, logger)

	case ComponentTypeFont:
		return ExportFontPack(exportName, logger)

	default:
		return fmt.Errorf("unsupported component type: %d", componentType)
	}
}

// ExportThemeComponents exports specific components to a theme package
func ExportThemeComponents(themeName string, selectedComponents map[ComponentType]bool, logger *Logger) error {
	logger.Printf("Exporting specific components to theme: %s", themeName)

	// Create theme directory
	themePath, err := CreateComponentExportDirectory(themeName, ComponentTypeFullTheme)
	if err != nil {
		logger.Printf("Error creating theme directory: %v", err)
		return fmt.Errorf("error creating theme directory: %w", err)
	}

	// Initialize manifest
	manifest := &ThemeManifest{}
	manifest.ThemeInfo.Name = themeName
	manifest.ThemeInfo.Version = "1.0.0"
	manifest.ThemeInfo.Author = "AuthorName" // Default author name as requested
	manifest.ThemeInfo.CreationDate = getCurrentTime()
	manifest.ThemeInfo.ExportedBy = GetVersionString()

	// Export selected components
	if selectedComponents[ComponentTypeWallpaper] {
		logger.Printf("Exporting wallpapers")
		if err := ExportWallpapers(themePath, manifest, logger); err != nil {
			logger.Printf("Error exporting wallpapers: %v", err)
			// Continue with other components
		}
	}

	if selectedComponents[ComponentTypeIcon] {
		logger.Printf("Exporting icons")
		if err := ExportIcons(themePath, manifest, logger); err != nil {
			logger.Printf("Error exporting icons: %v", err)
			// Continue with other components
		}
	}

	if selectedComponents[ComponentTypeFont] {
		logger.Printf("Exporting fonts")
		if err := ExportFonts(themePath, manifest, logger); err != nil {
			logger.Printf("Error exporting fonts: %v", err)
			// Continue with other components
		}
	}

	if selectedComponents[ComponentTypeAccent] || selectedComponents[ComponentTypeLED] {
		logger.Printf("Exporting settings")
		if err := ExportSettings(themePath, manifest, logger); err != nil {
			logger.Printf("Error exporting settings: %v", err)
			// Continue with other components
		}
	}

	// Generate preview image
	if err := GeneratePreview(themePath, logger); err != nil {
		logger.Printf("Error generating preview: %v", err)
		// Continue anyway
	}

	// Write manifest
	if err := WriteManifest(themePath, manifest, logger); err != nil {
		logger.Printf("Error writing manifest: %v", err)
		return fmt.Errorf("error writing manifest: %w", err)
	}

	logger.Printf("Theme component export completed successfully: %s", themePath)
	ui.ShowMessage(fmt.Sprintf("Theme '%s' exported successfully", themeName), "3")

	return nil
}

// CreateComponentExportDirectory creates a directory for exporting a component
func CreateComponentExportDirectory(exportName string, componentType ComponentType) (string, error) {
	// Get the current directory
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("error getting current directory: %w", err)
	}

	// Determine export directory and file extension based on component type
	var exportDir string
	var fileExtension string

	switch componentType {
	case ComponentTypeFullTheme:
		exportDir = filepath.Join(cwd, "Themes", "Exports")
		fileExtension = ".theme"
	case ComponentTypeAccent:
		exportDir = filepath.Join(cwd, "Accents", "Exports")
		fileExtension = ".acc"
	case ComponentTypeLED:
		exportDir = filepath.Join(cwd, "LEDs", "Exports")
		fileExtension = ".led"
	case ComponentTypeWallpaper:
		exportDir = filepath.Join(cwd, "Wallpapers", "Exports")
		fileExtension = ".bg"
	case ComponentTypeIcon:
		exportDir = filepath.Join(cwd, "Icons", "Exports")
		fileExtension = ".icon"
	case ComponentTypeFont:
		exportDir = filepath.Join(cwd, "Fonts", "Exports")
		fileExtension = ".font"
	default:
		return "", fmt.Errorf("invalid component type: %d", componentType)
	}

	// Ensure directory exists
	if err := os.MkdirAll(exportDir, 0755); err != nil {
		return "", fmt.Errorf("error creating export directory: %w", err)
	}

	// If exportName already has the correct extension, use it as is
	if !strings.HasSuffix(exportName, fileExtension) {
		exportName = exportName + fileExtension
	}

	// Full path to export directory
	exportPath := filepath.Join(exportDir, exportName)

	// Create the export directory
	if err := os.MkdirAll(exportPath, 0755); err != nil {
		return "", fmt.Errorf("error creating export directory: %w", err)
	}

	// For different component types, create appropriate subdirectories
	switch componentType {
	case ComponentTypeFullTheme:
		// Create theme subdirectories
		subDirs := []string{
			"Wallpapers/SystemWallpapers",
			"Wallpapers/CollectionWallpapers",
			"Icons/SystemIcons",
			"Icons/ToolIcons",
			"Icons/CollectionIcons",
			"Fonts",
			"Settings",
		}

		for _, dir := range subDirs {
			path := filepath.Join(exportPath, dir)
			if err := os.MkdirAll(path, 0755); err != nil {
				return "", fmt.Errorf("error creating subdirectory %s: %w", dir, err)
			}
		}

	case ComponentTypeWallpaper:
		// Create wallpaper subdirectories
		subDirs := []string{
			"SystemWallpapers",
			"CollectionWallpapers",
		}

		for _, dir := range subDirs {
			path := filepath.Join(exportPath, dir)
			if err := os.MkdirAll(path, 0755); err != nil {
				return "", fmt.Errorf("error creating subdirectory %s: %w", dir, err)
			}
		}

	case ComponentTypeIcon:
		// Create icon subdirectories
		subDirs := []string{
			"SystemIcons",
			"ToolIcons",
			"CollectionIcons",
		}

		for _, dir := range subDirs {
			path := filepath.Join(exportPath, dir)
			if err := os.MkdirAll(path, 0755); err != nil {
				return "", fmt.Errorf("error creating subdirectory %s: %w", dir, err)
			}
		}
	}

	return exportPath, nil
}

// ExportAccentPack exports an accent pack (.acc)
func ExportAccentPack(exportName string, logger *Logger) error {
	logger.Printf("Exporting accent pack: %s", exportName)

	// Create export directory
	exportPath, err := CreateComponentExportDirectory(exportName, ComponentTypeAccent)
	if err != nil {
		return fmt.Errorf("error creating export directory: %w", err)
	}

	// Copy the current accent settings file
	srcPath := accents.SettingsPath
	dstPath := filepath.Join(exportPath, "minuisettings.txt")

	// Check if source file exists
	if _, err := os.Stat(srcPath); os.IsNotExist(err) {
		return fmt.Errorf("accent settings file not found: %s", srcPath)
	}

	// Copy the file
	if err := CopyFile(srcPath, dstPath); err != nil {
		return fmt.Errorf("error copying accent settings: %w", err)
	}

	logger.Printf("Accent pack exported successfully: %s", exportPath)
	ui.ShowMessage(fmt.Sprintf("Accent pack '%s' exported successfully", exportName), "3")

	return nil
}

// ExportLEDPack exports a LED pack (.led)
func ExportLEDPack(exportName string, logger *Logger) error {
	logger.Printf("Exporting LED pack: %s", exportName)

	// Create export directory
	exportPath, err := CreateComponentExportDirectory(exportName, ComponentTypeLED)
	if err != nil {
		return fmt.Errorf("error creating export directory: %w", err)
	}

	// Copy the current LED settings file
	srcPath := "/mnt/SDCARD/.userdata/shared/ledsettings_brick.txt"
	dstPath := filepath.Join(exportPath, "ledsettings_brick.txt")

	// Check if source file exists
	if _, err := os.Stat(srcPath); os.IsNotExist(err) {
		return fmt.Errorf("LED settings file not found: %s", srcPath)
	}

	// Copy the file
	if err := CopyFile(srcPath, dstPath); err != nil {
		return fmt.Errorf("error copying LED settings: %w", err)
	}

	logger.Printf("LED pack exported successfully: %s", exportPath)
	ui.ShowMessage(fmt.Sprintf("LED pack '%s' exported successfully", exportName), "3")

	return nil
}

// ExportWallpaperPack exports a wallpaper pack (.bg)
func ExportWallpaperPack(exportName string, logger *Logger) error {
	logger.Printf("Exporting wallpaper pack: %s", exportName)

	// Create export directory
	exportPath, err := CreateComponentExportDirectory(exportName, ComponentTypeWallpaper)
	if err != nil {
		return fmt.Errorf("error creating export directory: %w", err)
	}

	// Get system paths
	// systemPaths, err := system.GetSystemPaths()
	// if err != nil {
	//     return fmt.Errorf("error getting system paths: %w", err)
	// }

	// Create a manifest just for tracking what we've exported
	manifest := &ThemeManifest{}

	// Export wallpapers only
	if err := ExportWallpapers(exportPath, manifest, logger); err != nil {
		return fmt.Errorf("error exporting wallpapers: %w", err)
	}

	logger.Printf("Wallpaper pack exported successfully: %s", exportPath)
	ui.ShowMessage(fmt.Sprintf("Wallpaper pack '%s' exported successfully", exportName), "3")

	return nil
}

// ExportIconPack exports an icon pack (.icon)
func ExportIconPack(exportName string, logger *Logger) error {
	logger.Printf("Exporting icon pack: %s", exportName)

	// Create export directory
	exportPath, err := CreateComponentExportDirectory(exportName, ComponentTypeIcon)
	if err != nil {
		return fmt.Errorf("error creating export directory: %w", err)
	}

	// Create a manifest just for tracking what we've exported
	manifest := &ThemeManifest{}

	// Export icons only
	if err := ExportIcons(exportPath, manifest, logger); err != nil {
		return fmt.Errorf("error exporting icons: %w", err)
	}

	logger.Printf("Icon pack exported successfully: %s", exportPath)
	ui.ShowMessage(fmt.Sprintf("Icon pack '%s' exported successfully", exportName), "3")

	return nil
}

// ExportFontPack exports a font pack (.font)
func ExportFontPack(exportName string, logger *Logger) error {
	logger.Printf("Exporting font pack: %s", exportName)

	// Create export directory
	exportPath, err := CreateComponentExportDirectory(exportName, ComponentTypeFont)
	if err != nil {
		return fmt.Errorf("error creating export directory: %w", err)
	}

	// Create a manifest just for tracking what we've exported
	manifest := &ThemeManifest{}

	// Export fonts only
	if err := ExportFonts(exportPath, manifest, logger); err != nil {
		return fmt.Errorf("error exporting fonts: %w", err)
	}

	logger.Printf("Font pack exported successfully: %s", exportPath)
	ui.ShowMessage(fmt.Sprintf("Font pack '%s' exported successfully", exportName), "3")

	return nil
}

// generateSequentialFileName generates a sequential file name based on component type
func generateSequentialFileName(componentType ComponentType) string {
	// Get the current directory
	cwd, err := os.Getwd()
	if err != nil {
		return "export_1"
	}

	// Determine prefix and directory based on component type
	var prefix string
	var exportDir string
	var fileExt string

	switch componentType {
	case ComponentTypeFullTheme:
		prefix = "theme"
		exportDir = filepath.Join(cwd, "Themes", "Exports")
		fileExt = ".theme"
	case ComponentTypeAccent:
		prefix = "accent"
		exportDir = filepath.Join(cwd, "Accents", "Exports")
		fileExt = ".acc"
	case ComponentTypeLED:
		prefix = "led"
		exportDir = filepath.Join(cwd, "LEDs", "Exports")
		fileExt = ".led"
	case ComponentTypeWallpaper:
		prefix = "wallpaper"
		exportDir = filepath.Join(cwd, "Wallpapers", "Exports")
		fileExt = ".bg"
	case ComponentTypeIcon:
		prefix = "icon"
		exportDir = filepath.Join(cwd, "Icons", "Exports")
		fileExt = ".icon"
	case ComponentTypeFont:
		prefix = "font"
		exportDir = filepath.Join(cwd, "Fonts", "Exports")
		fileExt = ".font"
	default:
		prefix = "export"
		exportDir = filepath.Join(cwd, "Exports")
		fileExt = ""
	}

	// Create the directory if it doesn't exist
	if err := os.MkdirAll(exportDir, 0755); err != nil {
		return prefix + "_1"
	}

	// Find highest existing number
	highestNum := 0
	regex := regexp.MustCompile(fmt.Sprintf(`^%s_(\d+)%s$`, prefix, fileExt))

	entries, err := os.ReadDir(exportDir)
	if err == nil {
		for _, entry := range entries {
			if entry.IsDir() {
				name := entry.Name()
				matches := regex.FindStringSubmatch(name)
				if len(matches) == 2 {
					num, err := strconv.Atoi(matches[1])
					if err == nil && num > highestNum {
						highestNum = num
					}
				}
			}
		}
	}

	// Generate new file name with the next number
	return fmt.Sprintf("%s_%d", prefix, highestNum+1)
}

// Connect the UI's performExport function to our backend implementation
func PerformExport(componentType ComponentType, exportName string, selectedComponents map[ComponentType]bool) error {
	return ExportComponent(componentType, exportName, selectedComponents)
}
